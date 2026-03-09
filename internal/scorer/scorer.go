package scorer

import "github.com/bbik/ari/internal/checker"

// PillarScore holds pass/fail stats for a pillar
type PillarScore struct {
	Pillar checker.Pillar
	Passed int
	Total  int
	Rate   float64
}

// Score holds the overall evaluation result
type Score struct {
	Level        checker.Level
	PassRate     float64
	PillarScores map[checker.Pillar]*PillarScore
	LevelScores  map[checker.Level]*LevelScore
}

// LevelScore holds pass/fail stats for a maturity level
type LevelScore struct {
	Level  checker.Level
	Passed int
	Total  int
	Rate   float64
}
