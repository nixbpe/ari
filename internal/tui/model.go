package tui

// ViewID identifies which view is currently displayed
type ViewID int

const (
	ViewProgress ViewID = iota
	ViewReport
	ViewDetail
)

// Model is the root Bubbletea model (placeholder - implemented in T22)
type Model struct {
	currentView ViewID
}
