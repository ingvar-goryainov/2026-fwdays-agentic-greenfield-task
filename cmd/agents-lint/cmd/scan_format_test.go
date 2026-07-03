package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateFormat(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{"text is valid", "text", false},
		{"sarif is valid", "sarif", false},
		{"empty is invalid", "", true},
		{"unknown value is invalid", "xml", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFormat(tt.format)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.format)
				assert.Contains(t, err.Error(), "text")
				assert.Contains(t, err.Error(), "sarif")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRunScan_Format(t *testing.T) {
	t.Run("default text format is unchanged", func(t *testing.T) {
		var buf bytes.Buffer
		exitCode, err := runScan(&buf, "../../../testdata/S002/invalid.md", "", formatText)
		require.NoError(t, err)
		assert.Equal(t, 1, exitCode)
		assert.Contains(t, buf.String(), "error  S002")
	})

	t.Run("sarif format produces a SARIF document", func(t *testing.T) {
		var buf bytes.Buffer
		exitCode, err := runScan(&buf, "../../../testdata/S002/invalid.md", "", formatSARIF)
		require.NoError(t, err)
		assert.Equal(t, 1, exitCode)

		var doc struct {
			Version string `json:"version"`
			Runs    []struct {
				Results []struct {
					RuleID string `json:"ruleId"`
				} `json:"results"`
			} `json:"runs"`
		}
		require.NoError(t, json.Unmarshal(buf.Bytes(), &doc))
		assert.Equal(t, "2.1.0", doc.Version)
		require.Len(t, doc.Runs, 1)
		require.NotEmpty(t, doc.Runs[0].Results)
		assert.Equal(t, "S002", doc.Runs[0].Results[0].RuleID)
	})

	t.Run("exit code is identical across formats for the same findings", func(t *testing.T) {
		var textBuf, sarifBuf bytes.Buffer
		textExit, err := runScan(&textBuf, "../../../testdata/S002/invalid.md", "", formatText)
		require.NoError(t, err)
		sarifExit, err := runScan(&sarifBuf, "../../../testdata/S002/invalid.md", "", formatSARIF)
		require.NoError(t, err)
		assert.Equal(t, textExit, sarifExit)
	})
}
