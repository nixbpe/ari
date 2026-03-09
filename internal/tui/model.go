package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/bbik/ari/internal/tui/views"
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
		m.currentView = ReportView
	case ErrorMsg:
		m.err = msg.Err
	case tea.KeyMsg:
		key := msg.Key()
		switch {
		case msg.String() == "q":
			m.quitting = true
			return m, tea.Quit
		case key.Mod&tea.ModCtrl != 0 && key.Code == 'c':
			m.quitting = true
			return m, tea.Quit
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
