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

type imageResult struct {
	images []specs.Image
	err    error
}

type volumeResult struct {
	volumes []specs.Volume
	err     error
}

type networkResult struct {
	networks []specs.Network
	err      error
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
	imageResults := make(chan imageResult, 1)
	volumeResults := make(chan volumeResult, 1)
	networkResults := make(chan networkResult, 1)
	group, ctx := errgroup.WithContext(ctx)

	// Each data source polls independently; snapshots are merged below.
	group.Go(func() error {
		return p.pollContainers(ctx, containerResults)
	})
	group.Go(func() error {
		return p.pollStats(ctx, statsResults)
	})
	group.Go(func() error {
		return p.pollImages(ctx, imageResults)
	})
	group.Go(func() error {
		return p.pollVolumes(ctx, volumeResults)
	})
	group.Go(func() error {
		return p.pollNetworks(ctx, networkResults)
	})

	defer close(p.updates)
	defer group.Wait()

	snapshot := model.Snapshot{
		Stats: make(map[string]specs.ContainerStats),
	}
	var containerErr error
	var statsErr error
	var imageErr error
	var volumeErr error
	var networkErr error

	for {
		select {
		case result := <-statsResults:
			statsErr = result.err
			if result.err == nil {
				snapshot.Stats = result.stats
				p.ringBuffer.Push(snapshot)
			}
		case result := <-imageResults:
			imageErr = result.err
			if result.err == nil {
				snapshot.Images = result.images
				p.ringBuffer.Push(snapshot)
			}
		case result := <-volumeResults:
			volumeErr = result.err
			if result.err == nil {
				snapshot.Volumes = result.volumes
				p.ringBuffer.Push(snapshot)
			}
		case result := <-networkResults:
			networkErr = result.err
			if result.err == nil {
				snapshot.Networks = result.networks
				p.ringBuffer.Push(snapshot)
			}
		case result := <-containerResults:
			containerErr = result.err
			if containerErr == nil {
				snapshot.Containers = result.snapshot.Containers
				snapshot.Logs = result.snapshot.Logs
				p.ringBuffer.Push(snapshot)
			}
		case <-ctx.Done():
			return
		}

		select {
		case p.updates <- firstError(
			containerErr,
			statsErr,
			imageErr,
			volumeErr,
			networkErr,
		):
		case <-ctx.Done():
			return
		}
	}
}

func firstError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
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

func (p *poller) pollImages(ctx context.Context, results chan<- imageResult) error {
	ticker := time.NewTicker(p.pollingTime)
	defer ticker.Stop()

	for {
		images, err := containers.GetImages(ctx)

		select {
		case results <- imageResult{images: images, err: err}:
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

func (p *poller) pollVolumes(ctx context.Context, results chan<- volumeResult) error {
	ticker := time.NewTicker(p.pollingTime)
	defer ticker.Stop()

	for {
		volumes, err := containers.GetVolumes(ctx)

		select {
		case results <- volumeResult{volumes: volumes, err: err}:
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

func (p *poller) pollNetworks(ctx context.Context, results chan<- networkResult) error {
	ticker := time.NewTicker(p.pollingTime)
	defer ticker.Stop()

	for {
		networks, err := containers.GetNetworks(ctx)

		select {
		case results <- networkResult{networks: networks, err: err}:
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
	ticker := time.NewTicker(p.pollingTime)
	defer ticker.Stop()

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

		select {
		case <-ticker.C:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
