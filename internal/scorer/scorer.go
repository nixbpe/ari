package scorer

import "github.com/nixbpe/ari/internal/checker"

const passThreshold = 0.8

var levelCriterionCounts = map[checker.Level]int{
	checker.LevelFunctional:   11,
	checker.LevelDocumented:   17,
	checker.LevelStandardized: 19,
	checker.LevelOptimized:    20,
	checker.LevelAutonomous:   5,
}

type Scorer struct{}

func New() *Scorer {
	return &Scorer{}
}

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
	PillarScores map[checker.Pillar]PillarScore
	LevelScores  map[checker.Level]LevelScore
}

// LevelScore holds pass/fail stats for a maturity level
type LevelScore struct {
	Level    checker.Level
	Passed   int
	Total    int
	Rate     float64
	Achieved bool
}

func (s *Scorer) Calculate(results []*checker.Result) *Score {
	_ = s

	levelScores := make(map[checker.Level]LevelScore, len(levelCriterionCounts))
	skippedByLevel := make(map[checker.Level]int, len(levelCriterionCounts))

	for level, total := range levelCriterionCounts {
		levelScores[level] = LevelScore{Level: level, Total: total}
	}

	pillarScores := map[checker.Pillar]PillarScore{}

	overallPassed := 0
	overallTotal := 0

	for _, result := range results {
		if result == nil {
			continue
		}

		if result.Skipped {
			skippedByLevel[result.Level]++
			continue
		}

		overallTotal++
		if result.Passed {
			overallPassed++
		}

		pillar := pillarScores[result.Pillar]
		pillar.Pillar = result.Pillar
		pillar.Total++
		if result.Passed {
			pillar.Passed++
		}
		pillarScores[result.Pillar] = pillar

		level := levelScores[result.Level]
		level.Level = result.Level
		if level.Total == 0 {
			if expectedTotal, ok := levelCriterionCounts[result.Level]; ok {
				level.Total = expectedTotal
			} else {
				level.Total = 1
			}
		} else if _, ok := levelCriterionCounts[result.Level]; !ok {
			level.Total++
		}
		if result.Passed {
			level.Passed++
		}
		levelScores[result.Level] = level
	}

	for level, score := range levelScores {
		score.Total -= skippedByLevel[level]
		if score.Total < 0 {
			score.Total = 0
		}
		if score.Total > 0 {
			score.Rate = float64(score.Passed) / float64(score.Total)
		}
		levelScores[level] = score
	}

	for pillar, score := range pillarScores {
		if score.Total > 0 {
			score.Rate = float64(score.Passed) / float64(score.Total)
		}
		pillarScores[pillar] = score
	}

	achieved := checker.Level(0)
	for level := checker.LevelFunctional; level <= checker.LevelAutonomous; level++ {
		levelScore := levelScores[level]
		if levelScore.Total == 0 || levelScore.Rate < passThreshold {
			levelScores[level] = levelScore
			break
		}

		levelScore.Achieved = true
		levelScores[level] = levelScore
		achieved = level
	}

	passRate := 0.0
	if overallTotal > 0 {
		passRate = float64(overallPassed) / float64(overallTotal)
	}

	return &Score{
		Level:        achieved,
		PassRate:     passRate,
		PillarScores: pillarScores,
		LevelScores:  levelScores,
	}
}
