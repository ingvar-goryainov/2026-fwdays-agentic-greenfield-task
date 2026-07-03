package rules_test

import (
	"testing"

	"github.com/ingvar-goryainov/agents-lint/internal/lint/rules"
)

// BenchmarkRun_LargeFile measures rules.Run against a ~500-line AGENTS.md
// fixture, the boundary NFR-PERF-01 specifies (< 100ms for <= 500 lines).
// Reported as ns/op rather than asserted as a hard pass/fail, since a fixed
// wall-clock threshold in a regular test is flaky across machines/CI
// runners (see scan-command's design.md).
func BenchmarkRun_LargeFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := rules.Run("../../../testdata/scan/perf/large.md"); err != nil {
			b.Fatal(err)
		}
	}
}
