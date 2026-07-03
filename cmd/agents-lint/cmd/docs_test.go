package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunDocs(t *testing.T) {
	t.Run("known schema rule ID prints all fields", func(t *testing.T) {
		var buf bytes.Buffer
		err := runDocs(&buf, "S002")

		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "S002")
		assert.Contains(t, out, "error")
		assert.Contains(t, out, "Valid:")
		assert.Contains(t, out, "Invalid:")
		assert.Contains(t, out, "Fix:")
	})

	t.Run("known codebase-awareness rule ID prints all fields", func(t *testing.T) {
		var buf bytes.Buffer
		err := runDocs(&buf, "C001")

		require.NoError(t, err)
		out := buf.String()
		assert.Contains(t, out, "C001")
		assert.Contains(t, out, "Valid:")
		assert.Contains(t, out, "Invalid:")
		assert.Contains(t, out, "Fix:")
	})

	t.Run("S001 (no struct type) is documented too", func(t *testing.T) {
		var buf bytes.Buffer
		err := runDocs(&buf, "S001")

		require.NoError(t, err)
		assert.Contains(t, buf.String(), "S001")
	})

	t.Run("unknown rule ID errors without writing output", func(t *testing.T) {
		var buf bytes.Buffer
		err := runDocs(&buf, "S999")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "S999")
		assert.Empty(t, buf.String())
	})

	t.Run("lowercase rule ID is treated as unknown", func(t *testing.T) {
		var buf bytes.Buffer
		err := runDocs(&buf, "s002")

		require.Error(t, err)
		assert.Empty(t, buf.String())
	})
}

func TestDocsCmd(t *testing.T) {
	t.Run("missing rule ID argument is a usage error", func(t *testing.T) {
		var out bytes.Buffer
		rootCmd.SetOut(&out)
		rootCmd.SetErr(&out)
		rootCmd.SetArgs([]string{"docs"})
		t.Cleanup(func() { rootCmd.SetArgs(nil) })

		err := rootCmd.Execute()

		assert.Error(t, err)
	})

	t.Run("known rule ID via the command tree exits cleanly", func(t *testing.T) {
		var out bytes.Buffer
		rootCmd.SetOut(&out)
		rootCmd.SetErr(&out)
		rootCmd.SetArgs([]string{"docs", "S002"})
		t.Cleanup(func() { rootCmd.SetArgs(nil) })

		err := rootCmd.Execute()

		require.NoError(t, err)
		assert.Contains(t, out.String(), "S002")
	})

	t.Run("unknown rule ID via the command tree errors", func(t *testing.T) {
		var out bytes.Buffer
		rootCmd.SetOut(&out)
		rootCmd.SetErr(&out)
		rootCmd.SetArgs([]string{"docs", "S999"})
		t.Cleanup(func() { rootCmd.SetArgs(nil) })

		err := rootCmd.Execute()

		require.Error(t, err)
		assert.Contains(t, out.String(), "S999")
	})
}
