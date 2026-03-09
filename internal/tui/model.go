package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/nixbpe/ari/internal/checker"
	"github.com/nixbpe/ari/internal/tui/views"
)

type ProgressModel = views.ProgressModel
type ReportModel = views.ReportModel
type DetailModel = views.DetailModel

type ViewID int

const (
	ProgressView ViewID = iota
	ReportView
	DetailView
)

type Model struct {
	currentView ViewID
	progress    ProgressModel
	report      *ReportModel
	detail      *DetailModel
	results     []*checker.Result
	err         error
	quitting    bool
	width       int
	height      int
}

func NewModel() Model {
	return Model{
		currentView: ProgressView,
		progress:    ProgressModel{},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case ScanStartMsg:
		m.progress.SetTotal(msg.Total)

	case CheckerStartMsg:
		m.progress.SetCurrent(msg.Name)

	case CheckerCompleteMsg:
		m.progress.UpdateProgress(msg.Result, msg.Done, msg.Total)

	case ScanCompleteMsg:
		report := ReportModel{Score: msg.Score, Report: msg.Report}
		m.report = &report
		m.results = msg.Results
		m.currentView = ReportView

	case DrillDownMsg:
		var pillarResults []*checker.Result
		for _, r := range m.results {
			if r != nil && r.Pillar == msg.Pillar {
				pillarResults = append(pillarResults, r)
			}
		}
		detail := DetailModel{Pillar: msg.Pillar, Results: pillarResults}
		m.detail = &detail
		m.currentView = DetailView

	case BackMsg:
		m.currentView = ReportView

	case ErrorMsg:
		m.err = msg.Err

	case tea.InterruptMsg:
		m.quitting = true
		return m, tea.Quit

	case tea.KeyPressMsg:
		keyStr := msg.String()

		// Ctrl+C / q always quit
		if keyStr == "ctrl+c" || keyStr == "q" {
			m.quitting = true
			return m, tea.Quit
		}

		switch m.currentView {
		case ReportView:
			if m.report != nil {
				switch keyStr {
				case "up":
					m.report.MoveUp()
				case "down":
					m.report.MoveDown()
				case "enter":
					pillar := m.report.CurrentPillar()
					return m, func() tea.Msg { return DrillDownMsg{Pillar: pillar} }
				case "h":
					return m, func() tea.Msg { return OpenBrowserMsg{} }
				case "J":
					return m, func() tea.Msg { return ExportJSONMsg{} }
				}
			}

		case DetailView:
			if m.detail != nil {
				switch keyStr {
				case "up":
					m.detail.ScrollUp()
				case "down":
					m.detail.ScrollDown()
				case "esc", "backspace":
					return m, func() tea.Msg { return BackMsg{} }
				}
			}
		}
	}

	return m, nil
}

func (m Model) View() tea.View {
	if m.err != nil {
		return tea.NewView(fmt.Sprintf("Error: %v", m.err))
	}

	var content string
	switch m.currentView {
	case ProgressView:
		content = m.progress.View()
	case ReportView:
		if m.report != nil {
			content = m.report.View()
		}
	case DetailView:
		if m.detail != nil {
			content = m.detail.View()
		}
	}

	return tea.NewView(content)
}
