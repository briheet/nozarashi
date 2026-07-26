package poller

import (
	"context"
	"time"

	"github.com/briheet/nozarashi/internal/containers"
	"github.com/briheet/nozarashi/internal/specs"
	"github.com/briheet/nozarashi/internal/tui/model"
	"github.com/briheet/nozarashi/internal/tui/ringbuffer"
	"golang.org/x/sync/errgroup"
)

type poller struct {
	ringBuffer  *ringbuffer.RingBuffer[model.Snapshot]
	updates     chan<- error
	pollingTime time.Duration
}

type containerResult struct {
	snapshot model.Snapshot
	err      error
}

type statsResult struct {
	stats map[string]specs.ContainerStats
	err   error
}

func NewPoller(
	ringBuffer *ringbuffer.RingBuffer[model.Snapshot],
	updates chan<- error,
	pollingTime time.Duration,
) *poller {
	return &poller{
		ringBuffer:  ringBuffer,
		updates:     updates,
		pollingTime: pollingTime,
	}
}

func (p *poller) Poll(ctx context.Context) {
	containerResults := make(chan containerResult, 1)
	statsResults := make(chan statsResult, 1)
	group, ctx := errgroup.WithContext(ctx)

	// Container data and stats have different sampling times, so poll them independently.
	group.Go(func() error {
		return p.pollContainers(ctx, containerResults)
	})
	group.Go(func() error {
		return p.pollStats(ctx, statsResults)
	})

	defer close(p.updates)
	defer group.Wait()

	latestStats := make(map[string]specs.ContainerStats)
	var statsErr error

	// Merge both streams here so only one goroutine writes snapshots to the ring buffer.
	for {
		select {
		case result := <-statsResults:
			statsErr = result.err
			if result.err == nil {
				latestStats = result.stats
			}
		case result := <-containerResults:
			if result.err == nil {
				result.snapshot.Stats = latestStats
				p.ringBuffer.Push(result.snapshot)
			}

			if result.err == nil {
				result.err = statsErr
			}

			select {
			case p.updates <- result.err:
			case <-ctx.Done():
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (p *poller) pollContainers(ctx context.Context, results chan<- containerResult) error {
	ticker := time.NewTicker(p.pollingTime)
	defer ticker.Stop()

	for {
		result := containerResult{}
		runtimeContainers, err := containers.GetContainers(ctx, false)
		result.err = err

		if err == nil {
			result.snapshot = model.Snapshot{
				Containers: *runtimeContainers,
				Logs:       make(map[string][]string, len(runtimeContainers.Items)),
			}

			for _, container := range runtimeContainers.Items {
				lines, logErr := containers.GetContainerLogLines(ctx, container.ID, 200)
				if logErr != nil {
					lines = []string{"error: " + logErr.Error()}
				}
				result.snapshot.Logs[container.ID] = lines
			}
		}

		select {
		case results <- result:
		case <-ctx.Done():
			return ctx.Err()
		}

		select {
		case <-ticker.C:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (p *poller) pollStats(ctx context.Context, results chan<- statsResult) error {
	previousStats := make(map[string]specs.ContainerStats)
	var previousSample time.Time

	for {
		stats, err := containers.GetContainerStats(ctx)
		sampledAt := time.Now()
		result := statsResult{err: err}

		if err == nil {
			result.stats = make(map[string]specs.ContainerStats, len(stats))
			elapsed := sampledAt.Sub(previousSample).Microseconds()

			for _, current := range stats {
				previous, ok := previousStats[current.ID]
				if ok && elapsed > 0 && current.CPUUsageUsec >= previous.CPUUsageUsec {
					current.CPUPercent = float64(
						current.CPUUsageUsec-previous.CPUUsageUsec,
					) / float64(elapsed) * 100
					current.CPUPercentValid = true
				}
				result.stats[current.ID] = current
			}

			previousStats = result.stats
			previousSample = sampledAt
		}

		select {
		case results <- result:
		case <-ctx.Done():
			return ctx.Err()
		}

		if err != nil {
			select {
			case <-time.After(p.pollingTime):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}
