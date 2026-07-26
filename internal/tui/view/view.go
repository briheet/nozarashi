package view

import (
	"fmt"
	"strings"
	"text/tabwriter"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/briheet/nozarashi/internal/specs"
	"github.com/briheet/nozarashi/internal/tui/model"
)

const (
	brandText   = " Nozarashi"
	runtimeText = "Apple Container "
)

type teaModel struct {
	m *model.Model
}

type pollMsg struct {
	err error
	ok  bool
}

// viewStyles contains the immutable palette shared by every rendered frame.
type viewStyles struct {
	frame            lipgloss.Style
	border           lipgloss.Style
	accent           lipgloss.Style
	muted            lipgloss.Style
	header           lipgloss.Style
	selected         lipgloss.Style
	inactiveSelected lipgloss.Style
	success          lipgloss.Style
	failure          lipgloss.Style
	truncate         lipgloss.Style
}

var styles = viewStyles{
	frame: lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#4B5563")),
	border: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4B5563")),
	accent: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7DD3FC")).
		Bold(true),
	muted: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#94A3B8")),
	header: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CBD5E1")).
		Bold(true),
	selected: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#0F172A")).
		Background(lipgloss.Color("#7DD3FC")).
		Bold(true),
	inactiveSelected: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E2E8F0")).
		Background(lipgloss.Color("#334155")).
		Bold(true),
	success: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#86EFAC")).
		Bold(true),
	failure: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FCA5A5")).
		Bold(true),
	truncate: lipgloss.NewStyle(),
}

// frameRenderer contains the derived values shared by one rendered frame.
type frameRenderer struct {
	m          *model.Model
	snapshot   model.Snapshot
	innerWidth int
	listRows   int
	detailRows int
	cursor     int
}

func (t teaModel) Init() tea.Cmd {
	return t.waitForUpdate()
}

func (t teaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Switch on messages
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.m.Width = msg.Width
		t.m.Height = msg.Height
	case pollMsg:
		if !msg.ok {
			return t, tea.Quit
		}

		t.m.Err = msg.err
		if snapshot, ok := t.m.RingBuffer.Latest(); ok {
			switch t.m.ActivePanel {
			case 1, 2, 3:
				t.m.Cursor = min(t.m.Cursor, max(len(snapshot.Containers.Items)-1, 0))
			case 4:
				t.m.ResourceCursor = min(t.m.ResourceCursor, max(len(snapshot.Images)-1, 0))
			case 5:
				t.m.ResourceCursor = min(t.m.ResourceCursor, max(len(snapshot.Volumes)-1, 0))
			case 6:
				t.m.ResourceCursor = min(t.m.ResourceCursor, max(len(snapshot.Networks)-1, 0))
			}
		}
		t.scrollLogs(0)
		return t, t.waitForUpdate()
	case tea.MouseWheelMsg:
		mouse := msg.Mouse()
		height := max(t.m.Height, 24)
		listRows := max((height-10)*2/5, 5)
		if t.m.ActivePanel != 2 || mouse.Y <= listRows+6 || mouse.Y >= height-2 {
			break
		}

		switch mouse.Button {
		case tea.MouseWheelUp:
			t.scrollLogs(3)
		case tea.MouseWheelDown:
			t.scrollLogs(-3)
		}
	case tea.KeyPressMsg:
		if msg.Code == 'c' && msg.Mod.Contains(tea.ModCtrl) {
			return t, tea.Quit
		}

		switch msg.Code {
		case tea.KeyTab:
			if t.m.ActivePanel == 1 {
				t.m.ActivePanel = 2
			} else {
				t.m.ActivePanel = 1
			}
		case tea.KeyUp:
			switch t.m.ActivePanel {
			case 1:
				t.m.Cursor = max(t.m.Cursor-1, 0)
				t.m.LogOffset = 0
			case 2:
				t.scrollLogs(1)
			case 4, 5, 6:
				t.m.ResourceCursor = max(t.m.ResourceCursor-1, 0)
			}
		case tea.KeyDown:
			snapshot, _ := t.m.RingBuffer.Latest()
			switch t.m.ActivePanel {
			case 1:
				t.m.Cursor = min(t.m.Cursor+1, max(len(snapshot.Containers.Items)-1, 0))
				t.m.LogOffset = 0
			case 2:
				t.scrollLogs(-1)
			case 4:
				t.m.ResourceCursor = min(t.m.ResourceCursor+1, max(len(snapshot.Images)-1, 0))
			case 5:
				t.m.ResourceCursor = min(t.m.ResourceCursor+1, max(len(snapshot.Volumes)-1, 0))
			case 6:
				t.m.ResourceCursor = min(t.m.ResourceCursor+1, max(len(snapshot.Networks)-1, 0))
			}
		case tea.KeyPgUp:
			if t.m.ActivePanel == 2 {
				t.scrollLogs(max(t.m.Height/4, 1))
			}
		case tea.KeyPgDown:
			if t.m.ActivePanel == 2 {
				t.scrollLogs(-max(t.m.Height/4, 1))
			}
		case tea.KeyHome:
			if t.m.ActivePanel == 2 {
				t.scrollLogs(1 << 30)
			}
		case tea.KeyEnd:
			if t.m.ActivePanel == 2 {
				t.m.LogOffset = 0
			}
		default:
			if msg.Mod != 0 {
				break
			}

			switch msg.Code {
			case 'q':
				return t, tea.Quit
			case 'k':
				switch t.m.ActivePanel {
				case 1:
					t.m.Cursor = max(t.m.Cursor-1, 0)
					t.m.LogOffset = 0
				case 2:
					t.scrollLogs(1)
				case 4, 5, 6:
					t.m.ResourceCursor = max(t.m.ResourceCursor-1, 0)
				}
			case 'j':
				snapshot, _ := t.m.RingBuffer.Latest()
				switch t.m.ActivePanel {
				case 1:
					t.m.Cursor = min(t.m.Cursor+1, max(len(snapshot.Containers.Items)-1, 0))
					t.m.LogOffset = 0
				case 2:
					t.scrollLogs(-1)
				case 4:
					t.m.ResourceCursor = min(t.m.ResourceCursor+1, max(len(snapshot.Images)-1, 0))
				case 5:
					t.m.ResourceCursor = min(t.m.ResourceCursor+1, max(len(snapshot.Volumes)-1, 0))
				case 6:
					t.m.ResourceCursor = min(t.m.ResourceCursor+1, max(len(snapshot.Networks)-1, 0))
				}
			case 'u':
				if t.m.ActivePanel == 2 {
					t.scrollLogs(1)
				}
			case 'd':
				if t.m.ActivePanel == 2 {
					t.scrollLogs(-1)
				}
			case '1':
				t.m.ActivePanel = 1
			case '2':
				t.m.ActivePanel = 2
			case '3':
				t.m.ActivePanel = 3
			case '4':
				t.m.ActivePanel = 4
				t.m.ResourceCursor = 0
			case '5':
				t.m.ActivePanel = 5
				t.m.ResourceCursor = 0
			case '6':
				t.m.ActivePanel = 6
				t.m.ResourceCursor = 0
			}
		}
	}

	return t, nil
}

// scrollLogs moves the selected container's log window and clamps it to available rows.
func (t teaModel) scrollLogs(lines int) {
	snapshot, ok := t.m.RingBuffer.Latest()
	if !ok || len(snapshot.Containers.Items) == 0 {
		return
	}

	cursor := min(t.m.Cursor, len(snapshot.Containers.Items)-1)
	logs := snapshot.Logs[snapshot.Containers.Items[cursor].ID]
	height := max(t.m.Height, 24)
	listRows := max((height-10)*2/5, 5)
	detailRows := height - listRows - 10
	t.m.LogOffset = min(
		max(t.m.LogOffset+lines, 0),
		max(len(logs)-detailRows, 0),
	)
}

// View composes the independent header, container, detail and footer components.
func (t teaModel) View() tea.View {
	renderer := newFrameRenderer(t.m)
	lines := renderer.renderHeader()
	lines = append(lines, renderer.renderPanel()...)
	lines = append(lines, renderer.renderFooter()...)

	content := styles.frame.
		Width(renderer.innerWidth + 2).
		Render(strings.Join(lines, "\n"))
	view := tea.NewView(content)
	view.AltScreen = true
	view.MouseMode = tea.MouseModeCellMotion
	return view
}

func newFrameRenderer(m *model.Model) frameRenderer {
	width := max(m.Width, 80)
	height := max(m.Height, 24)
	listRows := max((height-10)*2/5, 5)
	snapshot, _ := m.RingBuffer.Latest()

	return frameRenderer{
		m:          m,
		snapshot:   snapshot,
		innerWidth: width - 2,
		listRows:   listRows,
		detailRows: height - listRows - 10,
		cursor:     min(m.Cursor, max(len(snapshot.Containers.Items)-1, 0)),
	}
}

// renderHeader displays host and polling information above the fixed container panel.
func (r frameRenderer) renderHeader() []string {
	brandGap := strings.Repeat(" ", max(r.innerWidth-len(brandText)-len(runtimeText), 1))

	return []string{
		styles.accent.Render(brandText) +
			brandGap +
			styles.muted.Render(runtimeText),
		styles.muted.Render(r.fit(" Hostname: " + r.m.Hostname)),
		styles.muted.Render(r.fit(fmt.Sprintf(
			" Containers: %d    Images: %d    Volumes: %d    Networks: %d    Poll: %s",
			len(r.snapshot.Containers.Items),
			len(r.snapshot.Images),
			len(r.snapshot.Volumes),
			len(r.snapshot.Networks),
			pollInterval,
		))),
	}
}

// renderContainerPanel renders the visible container window and its live metrics.
func (r frameRenderer) renderContainerPanel() []string {
	containers := r.snapshot.Containers
	start := max(r.cursor-r.listRows+1, 0)
	end := min(start+r.listRows, len(containers.Items))
	nameWidth := max((r.innerWidth-101)/2, 12)
	imageWidth := max(r.innerWidth-nameWidth-101, 12)
	table := &strings.Builder{}
	writer := tabwriter.NewWriter(table, 0, 4, 2, ' ', 0)

	fmt.Fprintln(writer, "  NAME\tIMAGE\tSTATE\tCPU %\tMEMORY\tNET RX/TX\tBLOCK R/W\tPIDS")
	for index := start; index < end; index++ {
		container := containers.Items[index]
		cpu, memory, network, block, processes := r.renderMetrics(container.ID)
		fmt.Fprintf(
			writer,
			"  %.*s\t%.*s\t%.11s\t%s\t%s\t%s\t%s\t%s\n",
			nameWidth,
			container.ID,
			imageWidth,
			container.Configuration.Image.Reference,
			container.Status.State,
			cpu,
			memory,
			network,
			block,
			processes,
		)
	}
	writer.Flush()
	tableRows := strings.Split(strings.TrimSuffix(table.String(), "\n"), "\n")

	lines := []string{
		r.title("1 CONTAINERS", r.m.ActivePanel == 1),
		styles.header.Render(r.fit(tableRows[0])),
	}

	for index := start; index < end; index++ {
		row := r.fit(tableRows[index-start+1])
		if index == r.cursor {
			if r.m.ActivePanel == 1 {
				row = styles.selected.Render(row)
			} else {
				row = styles.inactiveSelected.Render(row)
			}
		}
		lines = append(lines, row)
	}
	for range r.listRows - (end - start) {
		lines = append(lines, r.fit(""))
	}

	return lines
}

// renderPanel routes all six panel shortcuts.
func (r frameRenderer) renderPanel() []string {
	switch r.m.ActivePanel {
	case 1, 2, 3:
		return append(r.renderContainerPanel(), r.renderDetailPanel()...)
	case 4:
		rows := make([][]string, 0, len(r.snapshot.Images))
		for _, image := range r.snapshot.Images {
			id := strings.TrimPrefix(image.ID, "sha256:")
			if len(id) > 12 {
				id = id[:12]
			}
			var size uint64
			platforms := make([]string, 0, len(image.Variants))
			for _, variant := range image.Variants {
				if variant.Size > 0 {
					size += uint64(variant.Size)
				}
				platforms = append(platforms, variant.Platform.OS+"/"+variant.Platform.Architecture)
			}
			rows = append(rows, []string{
				image.Configuration.Name,
				id,
				formatBytes(size),
				strings.Join(platforms, ","),
				image.Configuration.CreationDate,
			})
		}
		return r.renderResourceTable(
			"4 IMAGES",
			r.m.ResourceCursor,
			[]string{"REFERENCE", "IMAGE ID", "SIZE", "PLATFORM", "CREATED"},
			rows,
		)
	case 5:
		rows := make([][]string, 0, len(r.snapshot.Volumes))
		for _, volume := range r.snapshot.Volumes {
			rows = append(rows, []string{
				volume.ID,
				volume.Configuration.Driver,
				volume.Configuration.Format,
				formatBytes(volume.Configuration.SizeInBytes),
				volume.Configuration.CreationDate,
			})
		}
		return r.renderResourceTable(
			"5 VOLUMES",
			r.m.ResourceCursor,
			[]string{"VOLUME", "DRIVER", "FORMAT", "CAPACITY", "CREATED"},
			rows,
		)
	case 6:
		rows := make([][]string, 0, len(r.snapshot.Networks))
		for _, network := range r.snapshot.Networks {
			rows = append(rows, []string{
				network.ID,
				network.Configuration.Mode,
				network.Status.IPv4Subnet,
				network.Status.IPv4Gateway,
				network.Status.IPv6Subnet,
			})
		}
		return r.renderResourceTable(
			"6 NETWORKS",
			r.m.ResourceCursor,
			[]string{"NETWORK", "MODE", "IPv4 SUBNET", "GATEWAY", "IPv6 SUBNET"},
			rows,
		)
	default:
		return nil
	}
}

func (r frameRenderer) renderResourceTable(
	title string,
	cursor int,
	header []string,
	rows [][]string,
) []string {
	table := &strings.Builder{}
	writer := tabwriter.NewWriter(table, 0, 4, 2, ' ', 0)
	fmt.Fprintln(writer, "  "+strings.Join(header, "\t"))
	for _, row := range rows {
		fmt.Fprintln(writer, "  "+strings.Join(row, "\t"))
	}
	writer.Flush()
	tableRows := strings.Split(strings.TrimSuffix(table.String(), "\n"), "\n")

	visibleRows := r.listRows + r.detailRows + 1
	cursor = min(cursor, max(len(rows)-1, 0))
	start := max(cursor-visibleRows+1, 0)
	end := min(start+visibleRows, len(rows))
	lines := []string{
		r.title(title, true),
		styles.header.Render(r.fit(tableRows[0])),
	}

	for index := start; index < end; index++ {
		row := r.fit(tableRows[index+1])
		if index == cursor {
			row = styles.selected.Render(row)
		}
		lines = append(lines, row)
	}

	for len(lines) < visibleRows+2 {
		lines = append(lines, r.fit(""))
	}
	return lines
}

// renderMetrics formats one raw stats sample for the container table.
func (r frameRenderer) renderMetrics(id string) (
	cpu string,
	memory string,
	network string,
	block string,
	processes string,
) {
	cpu, memory, network, block, processes = "-", "-", "-", "-", "-"
	stats, ok := r.snapshot.Stats[id]
	if !ok {
		return
	}

	if stats.CPUPercentValid {
		cpu = fmt.Sprintf("%.2f%%", stats.CPUPercent)
	}
	if stats.MemoryLimitBytes > 0 {
		memory = fmt.Sprintf(
			"%s (%.1f%%)",
			formatBytes(stats.MemoryUsageBytes),
			float64(stats.MemoryUsageBytes)/float64(stats.MemoryLimitBytes)*100,
		)
	}
	network = formatBytes(stats.NetworkRxBytes) + "/" + formatBytes(stats.NetworkTxBytes)
	block = formatBytes(stats.BlockReadBytes) + "/" + formatBytes(stats.BlockWriteBytes)
	processes = fmt.Sprintf("%d", stats.NumProcesses)
	return
}

// renderDetailPanel displays logs by default and config for the selected container.
func (r frameRenderer) renderDetailPanel() []string {
	selected := "No container selected"
	details := make([]string, 0, r.detailRows)

	if len(r.snapshot.Containers.Items) > 0 {
		container := r.snapshot.Containers.Items[r.cursor]
		selected = container.ID
		if r.m.ActivePanel == 3 {
			details = r.renderConfig(container)
		} else {
			details = r.renderLogs(container.ID)
		}
	}

	panel := "2 LOGS"
	lowerActive := r.m.ActivePanel == 2
	if r.m.ActivePanel == 3 {
		panel = "3 CONFIG"
		lowerActive = true
	} else if r.m.LogOffset > 0 {
		panel += fmt.Sprintf(" · %d lines from latest", r.m.LogOffset)
	}

	lines := []string{r.title(panel+" · "+selected, lowerActive)}
	for index := range r.detailRows {
		if index >= len(details) {
			lines = append(lines, r.fit(""))
			continue
		}

		detail := r.fit(details[index])
		if r.m.ActivePanel == 3 {
			if separator := strings.Index(detail, ":"); separator >= 0 {
				detail = styles.accent.Render(detail[:separator+1]) + detail[separator+1:]
			}
		}
		lines = append(lines, detail)
	}

	return lines
}

// renderLogs selects the visible log window relative to the latest line.
func (r frameRenderer) renderLogs(id string) []string {
	logs := r.snapshot.Logs[id]
	if len(logs) == 0 {
		return []string{" No logs"}
	}

	end := max(len(logs)-r.m.LogOffset, 0)
	start := max(end-r.detailRows, 0)
	details := make([]string, 0, end-start)
	for _, log := range logs[start:end] {
		details = append(details, " "+log)
	}
	return details
}

// renderConfig builds the detailed static configuration for the selected container.
func (r frameRenderer) renderConfig(container specs.Container) []string {
	command := strings.Join(
		append(
			[]string{container.Configuration.InitProcess.Executable},
			container.Configuration.InitProcess.Arguments...,
		),
		" ",
	)
	details := []string{
		" Image:     " + container.Configuration.Image.Reference,
		" State:     " + container.Status.State,
		" Platform:  " + container.Configuration.Platform.OperatingSystem + "/" +
			container.Configuration.Platform.Architecture,
		fmt.Sprintf(
			" Resources: %d CPU, %d MB memory",
			container.Configuration.Resources.CPUs,
			container.Configuration.Resources.MemoryInBytes/(1024*1024),
		),
		" Started:   " + container.Status.StartedDate,
		" Command:   " + command,
		" Workdir:   " + container.Configuration.InitProcess.WorkingDirectory,
	}

	for _, network := range container.Status.Networks {
		details = append(
			details,
			fmt.Sprintf(
				" Network:   %s  %s  %s",
				network.Network,
				network.IPv4Address,
				network.MacAddress,
			),
		)
	}
	for _, port := range container.Configuration.PublishedPorts {
		details = append(
			details,
			fmt.Sprintf(
				" Port:      %s:%d -> %d/%s",
				port.HostAddress,
				port.HostPort,
				port.ContainerPort,
				port.Protocol,
			),
		)
	}
	for _, mount := range container.Configuration.Mounts {
		source := mount.Type.Volume.Name
		if source == "" {
			source = mount.Source
		}
		details = append(details, " Mount:     "+source+" -> "+mount.Destination)
	}
	for _, environment := range container.Configuration.InitProcess.Environment {
		details = append(details, " Env:       "+environment)
	}

	return details
}

// renderFooter shows panel shortcuts and the latest poll status.
func (r frameRenderer) renderFooter() []string {
	status := "connected"
	statusStyle := styles.success
	if r.m.Err != nil {
		status = "error: " + r.m.Err.Error()
		statusStyle = styles.failure
	}

	panelKey := func(panel int) string {
		key := fmt.Sprintf(" %d ", panel)
		if r.m.ActivePanel == panel {
			return styles.selected.Render(key)
		}
		return styles.accent.Render(key)
	}

	footer := panelKey(1) + " Containers  " +
		panelKey(2) + " Logs  " +
		panelKey(3) + " Config  " +
		panelKey(4) + " Images  " +
		panelKey(5) + " Volumes  " +
		panelKey(6) + " Networks  " +
		styles.accent.Render("Tab") + " Switch  " +
		styles.accent.Render("↑/↓") + " Navigate  " +
		styles.accent.Render("q") + " Quit  " +
		statusStyle.Render(status)
	if footerWidth := lipgloss.Width(footer); footerWidth < r.innerWidth {
		footer += strings.Repeat(" ", r.innerWidth-footerWidth)
	} else if footerWidth > r.innerWidth {
		footer = styles.truncate.MaxWidth(r.innerWidth).Render(footer)
	}

	return []string{
		styles.border.Render(strings.Repeat("─", r.innerWidth)),
		footer,
	}
}

// fit clamps content to the panel width and pads the remaining cells.
func (r frameRenderer) fit(value string) string {
	value = strings.ReplaceAll(value, "\n", " ")
	valueWidth := lipgloss.Width(value)
	if valueWidth > r.innerWidth {
		return styles.truncate.MaxWidth(r.innerWidth).Render(value)
	}
	return value + strings.Repeat(" ", r.innerWidth-valueWidth)
}

func (r frameRenderer) title(value string, active bool) string {
	value = " " + value + " "
	style := styles.muted
	if active {
		style = styles.accent
	}
	return styles.border.Render("─") +
		style.Render(value) +
		styles.border.Render(strings.Repeat("─", max(r.innerWidth-lipgloss.Width(value)-1, 0)))
}

func formatBytes(value uint64) string {
	switch {
	case value >= 1024*1024*1024:
		return fmt.Sprintf("%.2f GiB", float64(value)/(1024*1024*1024))
	case value >= 1024*1024:
		return fmt.Sprintf("%.2f MiB", float64(value)/(1024*1024))
	case value >= 1024:
		return fmt.Sprintf("%.2f KiB", float64(value)/1024)
	default:
		return fmt.Sprintf("%d B", value)
	}
}

func (t teaModel) waitForUpdate() tea.Cmd {
	return func() tea.Msg {
		err, ok := <-t.m.Updates
		return pollMsg{err: err, ok: ok}
	}
}
