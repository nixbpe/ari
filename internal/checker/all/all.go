package all

import (
	"github.com/nixbpe/ari/internal/checker"
	"github.com/nixbpe/ari/internal/checker/build"
	"github.com/nixbpe/ari/internal/checker/docs"
	"github.com/nixbpe/ari/internal/checker/style"
	checktesting "github.com/nixbpe/ari/internal/checker/testing"
	"github.com/nixbpe/ari/internal/llm"
)

// RegisterAll registers all 40 checkers into the registry.
// If eval is non-nil, LLM-capable checkers (naming_consistency, code_modularization) will use it.
func RegisterAll(r *checker.Registry, eval llm.Evaluator) {
	// Style & Validation — 12 checkers
	mustRegister(r, &style.LintConfigChecker{})
	mustRegister(r, &style.FormatterChecker{})
	mustRegister(r, &style.TypeCheckChecker{})
	mustRegister(r, &style.StrictTypingChecker{})
	mustRegister(r, &style.PreCommitHooksChecker{})
	mustRegister(r, style.NewNamingConsistencyChecker(eval))
	mustRegister(r, &style.CyclomaticComplexityChecker{})
	mustRegister(r, &style.DeadCodeDetectionChecker{})
	mustRegister(r, &style.DuplicateCodeDetectionChecker{})
	mustRegister(r, &style.CodeModularizationChecker{Evaluator: eval})
	mustRegister(r, &style.LargeFileDetectionChecker{})
	mustRegister(r, &style.TechDebtTrackingChecker{})

	// Build System — 13 checkers
	mustRegister(r, &build.BuildCmdDocChecker{})
	mustRegister(r, &build.SingleCommandSetupChecker{})
	mustRegister(r, &build.DepsPinnedChecker{})
	mustRegister(r, &build.FastCIFeedbackChecker{})
	mustRegister(r, &build.ReleaseAutomationChecker{})
	mustRegister(r, &build.DeploymentFrequencyChecker{})
	mustRegister(r, &build.VCSCliToolsChecker{})
	mustRegister(r, &build.AgenticDevelopmentChecker{})
	mustRegister(r, &build.AutomatedPRReviewChecker{})
	mustRegister(r, &build.BuildPerformanceTrackingChecker{})
	mustRegister(r, &build.FeatureFlagInfrastructureChecker{})
	mustRegister(r, &build.ReleaseNotesAutomationChecker{})
	mustRegister(r, &build.UnusedDependenciesDetectionChecker{})

	// Testing — 8 checkers
	mustRegister(r, &checktesting.UnitTestsExistChecker{})
	mustRegister(r, &checktesting.UnitTestsRunnableChecker{})
	mustRegister(r, &checktesting.TestNamingConventionsChecker{})
	mustRegister(r, &checktesting.TestIsolationChecker{})
	mustRegister(r, &checktesting.IntegrationTestsExistChecker{})
	mustRegister(r, &checktesting.TestCoverageThresholdsChecker{})
	mustRegister(r, &checktesting.FlakyTestDetectionChecker{})
	mustRegister(r, &checktesting.TestPerformanceTrackingChecker{})

	// Documentation — 7 checkers
	mustRegister(r, &docs.ReadmeChecker{})
	mustRegister(r, &docs.AgentsMdChecker{})
	mustRegister(r, &docs.DocumentationFreshnessChecker{})
	mustRegister(r, &docs.SkillsChecker{})
	mustRegister(r, &docs.AutomatedDocGenerationChecker{})
	mustRegister(r, &docs.ServiceFlowDocumentedChecker{})
	mustRegister(r, &docs.ApiSchemaDocsChecker{})
}

func mustRegister(r *checker.Registry, ch checker.Checker) {
	if err := r.Register(ch); err != nil {
		panic("checker registration failed: " + err.Error())
	}
}
