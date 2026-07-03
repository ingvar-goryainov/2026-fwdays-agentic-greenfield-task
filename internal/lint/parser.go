package lint

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

const frontmatterDelimiter = "---"

// Frontmatter is the raw YAML frontmatter block at the top of an AGENTS.md
// file, delimited by "---" lines. Parse does not validate the YAML itself;
// that is rule S005's job.
type Frontmatter struct {
	Raw  string
	Line int
}

// Heading is a top-level ("##") section heading.
type Heading struct {
	Text string
	Line int
}

// AgentBlock is a "###" agent definition found inside the required
// Agent(s) section, with the presence of each recognized field recorded.
type AgentBlock struct {
	Name            string
	Line            int
	HasRole         bool
	HasInstructions bool
	HasTools        bool
	HasContext      bool
}

// Document is the structural model rules operate on. It is produced once
// by Parse and never exposes the underlying goldmark AST, so rule
// implementations never need to import goldmark themselves.
type Document struct {
	Path        string
	Frontmatter *Frontmatter
	H2Sections  []Heading
	AgentBlocks []AgentBlock
}

// Parse reads the AGENTS.md file at path and builds a Document from it.
func Parse(path string) (*Document, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("lint: reading %s: %w", path, err)
	}

	frontmatter, body := extractFrontmatter(src)

	md := goldmark.New()
	root := md.Parser().Parse(text.NewReader(body))

	doc := &Document{Path: path, Frontmatter: frontmatter}
	populateSections(doc, root, body)
	return doc, nil
}

// extractFrontmatter detects a leading "---"-delimited block and returns
// it alongside a copy of src with that block's bytes blanked out (newlines
// preserved) so goldmark's line numbers for the remaining content still
// line up with the original file.
func extractFrontmatter(src []byte) (*Frontmatter, []byte) {
	starts := lineStarts(src)
	if lineText(src, starts, 0) != frontmatterDelimiter {
		return nil, src
	}

	closeIdx := -1
	for i := 1; i < len(starts); i++ {
		if lineText(src, starts, i) == frontmatterDelimiter {
			closeIdx = i
			break
		}
	}

	rawStart := len(src)
	if len(starts) > 1 {
		rawStart = starts[1]
	}

	var raw string
	var blankEnd int
	if closeIdx == -1 {
		raw = string(src[rawStart:])
		blankEnd = len(src)
	} else {
		raw = string(src[rawStart:starts[closeIdx]])
		if closeIdx+1 < len(starts) {
			blankEnd = starts[closeIdx+1]
		} else {
			blankEnd = len(src)
		}
	}

	blanked := make([]byte, len(src))
	copy(blanked, src)
	for i := 0; i < blankEnd; i++ {
		if blanked[i] != '\n' {
			blanked[i] = ' '
		}
	}

	return &Frontmatter{Raw: raw, Line: 1}, blanked
}

// lineStarts returns the byte offset of the start of every line in src.
func lineStarts(src []byte) []int {
	starts := []int{0}
	for i, b := range src {
		if b == '\n' {
			starts = append(starts, i+1)
		}
	}
	return starts
}

// lineText returns the trimmed (no trailing \r or \n) text of the line at
// index idx, using the offsets produced by lineStarts. The caller must
// ensure idx is a valid index into starts.
func lineText(src []byte, starts []int, idx int) string {
	start := starts[idx]
	end := len(src)
	if idx+1 < len(starts) {
		end = starts[idx+1] - 1
	}
	return string(bytes.TrimSuffix(src[start:end], []byte("\r")))
}

// lineNumber converts a byte offset into a 1-based line number.
func lineNumber(src []byte, pos int) int {
	line := 1
	for i := 0; i < pos && i < len(src); i++ {
		if src[i] == '\n' {
			line++
		}
	}
	return line
}

func populateSections(doc *Document, root ast.Node, src []byte) {
	inAgentSection := false
	var currentBlock *AgentBlock

	flushBlock := func() {
		if currentBlock != nil {
			doc.AgentBlocks = append(doc.AgentBlocks, *currentBlock)
			currentBlock = nil
		}
	}

	for n := root.FirstChild(); n != nil; n = n.NextSibling() {
		switch node := n.(type) {
		case *ast.Heading:
			headingText := nodeText(node, src)
			line := lineOf(node, src)
			switch node.Level {
			case 1:
				flushBlock()
				inAgentSection = false
			case 2:
				flushBlock()
				doc.H2Sections = append(doc.H2Sections, Heading{Text: headingText, Line: line})
				inAgentSection = isAgentSectionTitle(headingText)
			case 3:
				if inAgentSection {
					flushBlock()
					currentBlock = &AgentBlock{Name: headingText, Line: line}
				}
			}
		case *ast.Paragraph:
			if inAgentSection && currentBlock != nil {
				classifyParagraph(currentBlock, node, src)
			}
		}
	}
	flushBlock()
}

func isAgentSectionTitle(headingText string) bool {
	t := strings.ToLower(strings.TrimSpace(headingText))
	return t == "agent" || t == "agents"
}

// classifyParagraph inspects a paragraph inside an agent block and records
// which field (role/instructions/tools/context) it satisfies.
func classifyParagraph(block *AgentBlock, p *ast.Paragraph, src []byte) {
	label, ok := boldLabel(p, src)
	if !ok {
		if strings.TrimSpace(nodeText(p, src)) != "" {
			block.HasRole = true
		}
		return
	}

	switch strings.ToLower(label) {
	case "role":
		block.HasRole = true
	case "instructions", "responsibilities":
		block.HasInstructions = true
	case "tools":
		block.HasTools = true
	case "context":
		block.HasContext = true
	}
}

// boldLabel returns the label text (e.g. "Role") if the paragraph starts
// with a "**Label:**"-style strong-emphasis run, and false otherwise.
func boldLabel(p *ast.Paragraph, src []byte) (string, bool) {
	em, ok := p.FirstChild().(*ast.Emphasis)
	if !ok || em.Level != 2 {
		return "", false
	}
	label := strings.TrimSpace(nodeText(em, src))
	label = strings.TrimSuffix(label, ":")
	return label, true
}

// nodeText concatenates the raw text of every descendant *ast.Text node,
// skipping markup characters (bold/italic/code-span delimiters).
func nodeText(n ast.Node, src []byte) string {
	var buf bytes.Buffer
	var walk func(ast.Node)
	walk = func(n ast.Node) {
		if t, ok := n.(*ast.Text); ok {
			buf.Write(t.Segment.Value(src))
			return
		}
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			walk(c)
		}
	}
	walk(n)
	return strings.TrimSpace(buf.String())
}

func lineOf(n ast.Node, src []byte) int {
	lines := n.Lines()
	if lines.Len() == 0 {
		return 0
	}
	return lineNumber(src, lines.At(0).Start)
}
