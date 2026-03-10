package all

import (
	"github.com/nixbpe/ari/internal/checker"
	"github.com/nixbpe/ari/internal/checker/analytics"
	"github.com/nixbpe/ari/internal/checker/build"
	"github.com/nixbpe/ari/internal/checker/devenv"
	"github.com/nixbpe/ari/internal/checker/docs"
	"github.com/nixbpe/ari/internal/checker/observability"
	"github.com/nixbpe/ari/internal/checker/security"
	"github.com/nixbpe/ari/internal/checker/style"
	"github.com/nixbpe/ari/internal/checker/taskdiscovery"
	checktesting "github.com/nixbpe/ari/internal/checker/testing"
	"github.com/nixbpe/ari/internal/llm"
)

// RegisterAll registers all 72 checkers into the registry.
// If eval is non-nil, LLM-capable checkers will use it for enhanced evaluation.
func RegisterAll(r *checker.Registry, eval llm.Evaluator) {
	// Constraints & Governance — style (12 checkers)
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

	// Build (split: ContextIntent=1, EnvInfra=4, Constraints=2, Verification=6)
	mustRegister(r, &build.BuildCmdDocChecker{Evaluator: eval})
	mustRegister(r, &build.SingleCommandSetupChecker{})
	mustRegister(r, &build.DepsPinnedChecker{})
	mustRegister(r, &build.FastCIFeedbackChecker{})
	mustRegister(r, &build.ReleaseAutomationChecker{})
	mustRegister(r, &build.DeploymentFrequencyChecker{})
	mustRegister(r, &build.VCSCliToolsChecker{})
	mustRegister(r, &build.AgenticDevelopmentChecker{Evaluator: eval})
	mustRegister(r, &build.AutomatedPRReviewChecker{})
	mustRegister(r, &build.BuildPerformanceTrackingChecker{})
	mustRegister(r, &build.FeatureFlagInfrastructureChecker{})
	mustRegister(r, &build.ReleaseNotesAutomationChecker{})
	mustRegister(r, &build.UnusedDependenciesDetectionChecker{})

	// Verification & Feedback — testing (8 checkers)
	mustRegister(r, &checktesting.UnitTestsExistChecker{})
	mustRegister(r, &checktesting.UnitTestsRunnableChecker{Evaluator: eval})
	mustRegister(r, &checktesting.TestNamingConventionsChecker{})
	mustRegister(r, &checktesting.TestIsolationChecker{})
	mustRegister(r, &checktesting.IntegrationTestsExistChecker{})
	mustRegister(r, &checktesting.TestCoverageThresholdsChecker{})
	mustRegister(r, &checktesting.FlakyTestDetectionChecker{})
	mustRegister(r, &checktesting.TestPerformanceTrackingChecker{})

	// Context & Intent — docs (7 checkers)
	mustRegister(r, &docs.ReadmeChecker{Evaluator: eval})
	mustRegister(r, &docs.AgentsMdChecker{Evaluator: eval})
	mustRegister(r, &docs.DocumentationFreshnessChecker{})
	mustRegister(r, &docs.SkillsChecker{Evaluator: eval})
	mustRegister(r, &docs.AutomatedDocGenerationChecker{})
	mustRegister(r, &docs.ServiceFlowDocumentedChecker{})
	mustRegister(r, &docs.ApiSchemaDocsChecker{})

	// Environment & Infra — devenv (7 checkers)
	mustRegister(r, &devenv.EnvTemplateChecker{})
	mustRegister(r, &devenv.DevcontainerChecker{})
	mustRegister(r, &devenv.VersionPinningChecker{})
	mustRegister(r, &devenv.LocalServicesSetupChecker{})
	mustRegister(r, &devenv.IDEConfigChecker{})
	mustRegister(r, &devenv.DevcontainerQualityChecker{Evaluator: eval})
	mustRegister(r, &devenv.DatabaseSchemaChecker{})

	// Verification & Feedback — observability (8 checkers)
	mustRegister(r, &observability.StructuredLoggingChecker{})
	mustRegister(r, &observability.HealthChecksChecker{})
	mustRegister(r, &observability.ErrorTrackingChecker{})
	mustRegister(r, &observability.DistributedTracingChecker{})
	mustRegister(r, &observability.MetricsCollectionChecker{})
	mustRegister(r, &observability.AlertingConfiguredChecker{})
	mustRegister(r, &observability.ProfilingInstrumentationChecker{})
	mustRegister(r, &observability.RunbooksDocumentedChecker{Evaluator: eval})

	// Constraints & Governance — security (7 checkers)
	mustRegister(r, &security.SecurityPolicyChecker{})
	mustRegister(r, &security.GitignoreComprehensiveChecker{})
	mustRegister(r, &security.CodeownersChecker{})
	mustRegister(r, &security.DepUpdateAutomationChecker{})
	mustRegister(r, &security.SecretScanningConfigChecker{})
	mustRegister(r, &security.SASTConfigChecker{})
	mustRegister(r, &security.DependencyAuditCIChecker{})

	// Context & Intent — taskdiscovery (5 checkers)
	mustRegister(r, &taskdiscovery.ContributingGuideChecker{})
	mustRegister(r, &taskdiscovery.IssueTemplatesChecker{})
	mustRegister(r, &taskdiscovery.PRTemplateChecker{})
	mustRegister(r, &taskdiscovery.IssueLabelingSystemChecker{})
	mustRegister(r, &taskdiscovery.BacklogStructureDocsChecker{Evaluator: eval})

	// Context & Intent — analytics (5 checkers)
	mustRegister(r, &analytics.AnalyticsSdkChecker{})
	mustRegister(r, &analytics.TrackingPlanDocsChecker{})
	mustRegister(r, &analytics.ExperimentInfrastructureChecker{})
	mustRegister(r, &analytics.ProductMetricsDocsChecker{Evaluator: eval})
	mustRegister(r, &analytics.ErrorToInsightPipelineChecker{})
}

func mustRegister(r *checker.Registry, ch checker.Checker) {
	if err := r.Register(ch); err != nil {
		panic("checker registration failed: " + err.Error())
	}
}
