package view

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/briheet/nozarashi/internal/specs"
	"github.com/briheet/nozarashi/internal/tui/model"
	"github.com/briheet/nozarashi/internal/tui/ringbuffer"
)

func TestViewLayout(t *testing.T) {
	logs := make([]string, 20)
	for index := range logs {
		logs[index] = fmt.Sprintf("log %02d", index+1)
	}

	buffer := ringbuffer.NewRingBuffer[model.Snapshot](1)
	buffer.Push(model.Snapshot{
		Containers: specs.Containers{
			Items: []specs.Container{
				{
					ID: "api",
					Configuration: specs.ContainerConfiguration{
						Image: specs.ContainerImage{Reference: "example/api:latest"},
					},
					Status: specs.ContainerStatus{State: "running"},
				},
			},
		},
		Logs: map[string][]string{
			"api": logs,
		},
		Stats: map[string]specs.ContainerStats{
			"api": {
				ID:               "api",
				CPUPercent:       12.5,
				CPUPercentValid:  true,
				MemoryUsageBytes: 64 * 1024 * 1024,
				MemoryLimitBytes: 1024 * 1024 * 1024,
				NetworkRxBytes:   2 * 1024,
				NetworkTxBytes:   1024,
				BlockReadBytes:   4 * 1024,
				BlockWriteBytes:  2 * 1024,
				NumProcesses:     7,
			},
		},
	})

	baseModel := model.InitialModel(buffer, make(chan error))
	baseModel.Width = 220
	baseModel.Height = 30
	view := teaModel{m: baseModel}.View()

	if !view.AltScreen {
		t.Fatal("expected alternate screen view")
	}
	if view.MouseMode != tea.MouseModeCellMotion {
		t.Fatal("expected cell motion mouse mode")
	}
	if view.Cursor != nil {
		t.Fatal("expected terminal cursor to remain hidden")
	}
	lines := strings.Split(view.Content, "\n")
	if len(lines) != baseModel.Height {
		t.Fatalf("unexpected view height: %d", len(lines))
	}
	for index, line := range lines {
		if width := lipgloss.Width(line); width != baseModel.Width {
			t.Fatalf("unexpected width on line %d: %d", index+1, width)
		}
	}
	if !strings.Contains(view.Content, "example/api:latest") {
		t.Fatal("expected container image in view")
	}
	if !strings.Contains(view.Content, "12.50%") ||
		!strings.Contains(view.Content, "64.00 MiB (6.2%)") ||
		!strings.Contains(view.Content, "2.00 KiB/1.00 KiB") {
		t.Fatal("expected live container stats in view")
	}
	if !strings.Contains(view.Content, "log 20") {
		t.Fatal("expected latest container logs in view")
	}
	if baseModel.ActivePanel != 1 {
		t.Fatal("expected containers panel to be active initially")
	}

	updated, _ := (teaModel{m: baseModel}).Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if updated.(teaModel).m.ActivePanel != 2 {
		t.Fatal("expected tab to focus logs")
	}
	updated, _ = updated.(teaModel).Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if updated.(teaModel).m.ActivePanel != 1 {
		t.Fatal("expected tab to return to containers")
	}
	updated, _ = updated.(teaModel).Update(tea.KeyPressMsg{Code: '3'})
	if !strings.Contains(updated.(teaModel).View().Content, "3 CONFIG") {
		t.Fatal("expected config panel")
	}
	updated, _ = updated.(teaModel).Update(tea.KeyPressMsg{Code: '2'})
	updated, _ = updated.(teaModel).Update(tea.KeyPressMsg{Code: tea.KeyPgUp})
	scrolled := updated.(teaModel).View()
	if strings.Contains(scrolled.Content, "log 20") {
		t.Fatal("expected latest logs outside scrolled view")
	}
	if !strings.Contains(scrolled.Content, "lines from latest") {
		t.Fatal("expected log scroll position in view")
	}

	baseModel.LogOffset = 0
	updated, _ = (teaModel{m: baseModel}).Update(tea.MouseWheelMsg{
		X:      10,
		Y:      20,
		Button: tea.MouseWheelUp,
	})
	if updated.(teaModel).m.LogOffset == 0 {
		t.Fatal("expected mouse wheel to scroll logs")
	}
}
