package scorer_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/bbik/ari/internal/checker"
	"github.com/bbik/ari/internal/scorer"
)

func TestLevel1Achieved(t *testing.T) {
	results := buildResults(levelSpec{level: checker.LevelFunctional, passed: 7})

	score := scorer.New().Calculate(results)

	if score.Level != checker.LevelFunctional {
		t.Fatalf("Level = %v, want %v", score.Level, checker.LevelFunctional)
	}

	if !score.LevelScores[checker.LevelFunctional].Achieved {
		t.Fatal("Level 1 should be achieved")
	}
}

func TestLevel1Threshold(t *testing.T) {
	results := buildResults(levelSpec{level: checker.LevelFunctional, passed: 6, failed: 1})

	score := scorer.New().Calculate(results)

	if score.Level != checker.LevelFunctional {
		t.Fatalf("Level = %v, want %v", score.Level, checker.LevelFunctional)
	}

	if !nearlyEqual(score.LevelScores[checker.LevelFunctional].Rate, 6.0/7.0) {
		t.Fatalf("L1 rate = %f, want %f", score.LevelScores[checker.LevelFunctional].Rate, 6.0/7.0)
	}
}

func TestLevel1BelowThreshold(t *testing.T) {
	results := buildResults(levelSpec{level: checker.LevelFunctional, passed: 5, failed: 2})

	score := scorer.New().Calculate(results)

	if score.Level != 0 {
		t.Fatalf("Level = %v, want 0", score.Level)
	}
}

func TestLevel2Gated(t *testing.T) {
	results := buildResults(
		levelSpec{level: checker.LevelFunctional, passed: 7},
		levelSpec{level: checker.LevelDocumented, passed: 7, failed: 1},
	)

	score := scorer.New().Calculate(results)

	if score.Level != checker.LevelDocumented {
		t.Fatalf("Level = %v, want %v", score.Level, checker.LevelDocumented)
	}
}

func TestLevel2BlockedByL2Failure(t *testing.T) {
	results := buildResults(
		levelSpec{level: checker.LevelFunctional, passed: 7},
		levelSpec{level: checker.LevelDocumented, passed: 5, failed: 3},
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
		levelSpec{level: checker.LevelFunctional, passed: 7},
		levelSpec{level: checker.LevelDocumented, passed: 8},
		levelSpec{level: checker.LevelStandardized, passed: 9},
		levelSpec{level: checker.LevelOptimized, passed: 10, failed: 2},
	)

	score := scorer.New().Calculate(results)

	if score.Level != checker.LevelOptimized {
		t.Fatalf("Level = %v, want %v", score.Level, checker.LevelOptimized)
	}
}

func TestSkippedExcluded(t *testing.T) {
	results := buildResults(levelSpec{level: checker.LevelFunctional, passed: 5, failed: 1, skipped: 1})

	score := scorer.New().Calculate(results)

	if score.Level != checker.LevelFunctional {
		t.Fatalf("Level = %v, want %v", score.Level, checker.LevelFunctional)
	}

	l1 := score.LevelScores[checker.LevelFunctional]
	if l1.Total != 6 {
		t.Fatalf("L1 total = %d, want 6", l1.Total)
	}

	if !nearlyEqual(l1.Rate, 5.0/6.0) {
		t.Fatalf("L1 rate = %f, want %f", l1.Rate, 5.0/6.0)
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
		levelSpec{level: checker.LevelFunctional, skipped: 7},
		levelSpec{level: checker.LevelDocumented, skipped: 8},
		levelSpec{level: checker.LevelStandardized, skipped: 9},
		levelSpec{level: checker.LevelOptimized, skipped: 12},
		levelSpec{level: checker.LevelAutonomous, skipped: 4},
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
		levelSpec{level: checker.LevelFunctional, passed: 6, failed: 1},
		levelSpec{level: checker.LevelDocumented, passed: 5, failed: 3},
	)

	score := scorer.New().Calculate(results)

	if score.Level != checker.LevelFunctional {
		t.Fatalf("Level = %v, want %v", score.Level, checker.LevelFunctional)
	}
}

func TestPillarScores(t *testing.T) {
	results := []*checker.Result{
		{ID: "a", Level: checker.LevelFunctional, Pillar: checker.PillarStyleValidation, Passed: true},
		{ID: "b", Level: checker.LevelFunctional, Pillar: checker.PillarStyleValidation, Passed: false},
		{ID: "c", Level: checker.LevelFunctional, Pillar: checker.PillarBuildSystem, Passed: true},
		{ID: "d", Level: checker.LevelFunctional, Pillar: checker.PillarBuildSystem, Passed: true},
	}

	score := scorer.New().Calculate(results)

	style := score.PillarScores[checker.PillarStyleValidation]
	if style.Passed != 1 || style.Total != 2 || !nearlyEqual(style.Rate, 0.5) {
		t.Fatalf("style score = %+v, want passed=1 total=2 rate=0.5", style)
	}

	build := score.PillarScores[checker.PillarBuildSystem]
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
				Pillar: checker.PillarStyleValidation,
				Passed: true,
			})
			id++
		}
		for i := 0; i < spec.failed; i++ {
			results = append(results, &checker.Result{
				ID:     checker.CheckerID(fmt.Sprintf("c-%d", id)),
				Level:  spec.level,
				Pillar: checker.PillarStyleValidation,
				Passed: false,
			})
			id++
		}
		for i := 0; i < spec.skipped; i++ {
			results = append(results, &checker.Result{
				ID:      checker.CheckerID(fmt.Sprintf("c-%d", id)),
				Level:   spec.level,
				Pillar:  checker.PillarStyleValidation,
				Skipped: true,
			})
			id++
		}
	}

	return results
}

func nearlyEqual(got, want float64) bool {
	return math.Abs(got-want) < 1e-9
}
