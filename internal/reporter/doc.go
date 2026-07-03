// Package reporter renders []lint.Finding as human-readable text
// (FR-OUT-01, FR-OUT-02, FR-OUT-04). It writes to a caller-supplied
// io.Writer and never calls fmt.Print* directly, keeping it usable by any
// future consumer (CLI, HTTP API) without depending on one.
package reporter
