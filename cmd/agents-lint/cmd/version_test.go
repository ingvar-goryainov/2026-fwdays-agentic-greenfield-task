package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionFlag(t *testing.T) {
	t.Run("--version prints the version and exits 0", func(t *testing.T) {
		var out bytes.Buffer
		rootCmd.SetOut(&out)
		rootCmd.SetErr(&out)
		rootCmd.SetArgs([]string{"--version"})

		err := rootCmd.Execute()

		require.NoError(t, err)
		assert.Contains(t, out.String(), version)
	})

	t.Run("fallback version is dev with no ldflags override", func(t *testing.T) {
		assert.Equal(t, "dev", version)
	})

	t.Run("--version is not accepted on subcommands", func(t *testing.T) {
		var out bytes.Buffer
		rootCmd.SetOut(&out)
		rootCmd.SetErr(&out)
		rootCmd.SetArgs([]string{"scan", "--version"})

		err := rootCmd.Execute()

		assert.Error(t, err)
	})
}
