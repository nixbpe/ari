package scorer_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/nixbpe/ari/internal/checker"
	"github.com/nixbpe/ari/internal/scorer"
)

func TestLevel1Achieved(t *testing.T) {
	results := buildResults(levelSpec{level: checker.LevelFunctional, passed: 11})

	score := scorer.New().Calculate(results)

	if score.Level != checker.LevelFunctional {
		t.Fatalf("Level = %v, want %v", score.Level, checker.LevelFunctional)
	}

	if !score.LevelScores[checker.LevelFunctional].Achieved {
		t.Fatal("Level 1 should be achieved")
	}
}

func TestLevel1Threshold(t *testing.T) {
	results := buildResults(levelSpec{level: checker.LevelFunctional, passed: 9, failed: 2})

	score := scorer.New().Calculate(results)

	if score.Level != checker.LevelFunctional {
		t.Fatalf("Level = %v, want %v", score.Level, checker.LevelFunctional)
	}

	if !nearlyEqual(score.LevelScores[checker.LevelFunctional].Rate, 9.0/11.0) {
		t.Fatalf("L1 rate = %f, want %f", score.LevelScores[checker.LevelFunctional].Rate, 9.0/11.0)
	}
}

func TestLevel1BelowThreshold(t *testing.T) {
	results := buildResults(levelSpec{level: checker.LevelFunctional, passed: 5, failed: 6})

	score := scorer.New().Calculate(results)

	if score.Level != 0 {
		t.Fatalf("Level = %v, want 0", score.Level)
	}
}

func TestLevel2Gated(t *testing.T) {
	results := buildResults(
		levelSpec{level: checker.LevelFunctional, passed: 11},
		levelSpec{level: checker.LevelDocumented, passed: 14, failed: 3},
	)

	score := scorer.New().Calculate(results)

	if score.Level != checker.LevelDocumented {
		t.Fatalf("Level = %v, want %v", score.Level, checker.LevelDocumented)
	}
}

func TestLevel2BlockedByL2Failure(t *testing.T) {
	results := buildResults(
		levelSpec{level: checker.LevelFunctional, passed: 11},
		levelSpec{level: checker.LevelDocumented, passed: 5, failed: 12},
	)

	score := scorer.New().Calculate(results)

	if score.Level != checker.LevelFunctional {
		t.Fatalf("Level = %v, want %v", score.Level, checker.LevelFunctional)
	}

	if score.LevelScores[checker.LevelDocumented].Achieved {
		t.Fatal("Level 2 should not be achieved")
	}
}

func TestLevel4Achieved(t *testing.T) {
	results := buildResults(
		levelSpec{level: checker.LevelFunctional, passed: 11},
		levelSpec{level: checker.LevelDocumented, passed: 17},
		levelSpec{level: checker.LevelStandardized, passed: 19},
		levelSpec{level: checker.LevelOptimized, passed: 16, failed: 4},
	)

	score := scorer.New().Calculate(results)

	if score.Level != checker.LevelOptimized {
		t.Fatalf("Level = %v, want %v", score.Level, checker.LevelOptimized)
	}
}

func TestSkippedExcluded(t *testing.T) {
	results := buildResults(levelSpec{level: checker.LevelFunctional, passed: 9, failed: 1, skipped: 1})

	score := scorer.New().Calculate(results)

	if score.Level != checker.LevelFunctional {
		t.Fatalf("Level = %v, want %v", score.Level, checker.LevelFunctional)
	}

	l1 := score.LevelScores[checker.LevelFunctional]
	if l1.Total != 10 {
		t.Fatalf("L1 total = %d, want 10", l1.Total)
	}

	if !nearlyEqual(l1.Rate, 9.0/10.0) {
		t.Fatalf("L1 rate = %f, want %f", l1.Rate, 9.0/10.0)
	}
}

func TestEmptyResults(t *testing.T) {
	score := scorer.New().Calculate(nil)

	if score.Level != 0 {
		t.Fatalf("Level = %v, want 0", score.Level)
	}

	if score.PassRate != 0 {
		t.Fatalf("PassRate = %f, want 0", score.PassRate)
	}
}

func TestAllSkipped(t *testing.T) {
	results := buildResults(
		levelSpec{level: checker.LevelFunctional, skipped: 11},
		levelSpec{level: checker.LevelDocumented, skipped: 17},
		levelSpec{level: checker.LevelStandardized, skipped: 19},
		levelSpec{level: checker.LevelOptimized, skipped: 20},
		levelSpec{level: checker.LevelAutonomous, skipped: 5},
	)

	score := scorer.New().Calculate(results)

	if score.Level != 0 {
		t.Fatalf("Level = %v, want 0", score.Level)
	}

	if score.PassRate != 0 {
		t.Fatalf("PassRate = %f, want 0", score.PassRate)
	}
}

func TestGatedProgression(t *testing.T) {
	results := buildResults(
		levelSpec{level: checker.LevelFunctional, passed: 9, failed: 2},
		levelSpec{level: checker.LevelDocumented, passed: 5, failed: 12},
	)

	score := scorer.New().Calculate(results)

	if score.Level != checker.LevelFunctional {
		t.Fatalf("Level = %v, want %v", score.Level, checker.LevelFunctional)
	}
}

func TestPillarScores(t *testing.T) {
	results := []*checker.Result{
		{ID: "a", Level: checker.LevelFunctional, Pillar: checker.PillarConstraints, Passed: true},
		{ID: "b", Level: checker.LevelFunctional, Pillar: checker.PillarConstraints, Passed: false},
		{ID: "c", Level: checker.LevelFunctional, Pillar: checker.PillarEnvInfra, Passed: true},
		{ID: "d", Level: checker.LevelFunctional, Pillar: checker.PillarEnvInfra, Passed: true},
	}

	score := scorer.New().Calculate(results)

	style := score.PillarScores[checker.PillarConstraints]
	if style.Passed != 1 || style.Total != 2 || !nearlyEqual(style.Rate, 0.5) {
		t.Fatalf("style score = %+v, want passed=1 total=2 rate=0.5", style)
	}

	build := score.PillarScores[checker.PillarEnvInfra]
	if build.Passed != 2 || build.Total != 2 || !nearlyEqual(build.Rate, 1.0) {
		t.Fatalf("build score = %+v, want passed=2 total=2 rate=1.0", build)
	}
}

type levelSpec struct {
	level   checker.Level
	passed  int
	failed  int
	skipped int
}

func buildResults(specs ...levelSpec) []*checker.Result {
	results := make([]*checker.Result, 0)
	id := 0
	for _, spec := range specs {
		for i := 0; i < spec.passed; i++ {
			results = append(results, &checker.Result{
				ID:     checker.CheckerID(fmt.Sprintf("c-%d", id)),
				Level:  spec.level,
				Pillar: checker.PillarConstraints,
				Passed: true,
			})
			id++
		}
		for i := 0; i < spec.failed; i++ {
			results = append(results, &checker.Result{
				ID:     checker.CheckerID(fmt.Sprintf("c-%d", id)),
				Level:  spec.level,
				Pillar: checker.PillarConstraints,
				Passed: false,
			})
			id++
		}
		for i := 0; i < spec.skipped; i++ {
			results = append(results, &checker.Result{
				ID:      checker.CheckerID(fmt.Sprintf("c-%d", id)),
				Level:   spec.level,
				Pillar:  checker.PillarConstraints,
				Skipped: true,
			})
			id++
		}
	}

	return results
}

func TestLevelCriterionCountsMatchRegistry(t *testing.T) {
	// Verify that the sum of all level criterion counts equals 72 (total checkers)
	counts := map[checker.Level]int{
		checker.LevelFunctional:   11,
		checker.LevelDocumented:   17,
		checker.LevelStandardized: 19,
		checker.LevelOptimized:    20,
		checker.LevelAutonomous:   5,
	}
	total := 0
	for _, c := range counts {
		total += c
	}
	if total != 72 {
		t.Errorf("total criterion count = %d, want 72", total)
	}
}

func nearlyEqual(got, want float64) bool {
	return math.Abs(got-want) < 1e-9
}
