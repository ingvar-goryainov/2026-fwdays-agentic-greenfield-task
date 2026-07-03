package rules_test

import "os"

func writeTestFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o644)
}
