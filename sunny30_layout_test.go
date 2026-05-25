package main

import (
	"os"
	"strings"
	"testing"
)

func TestSunny30TableFillsDesktopWidthWithoutLooseColumns(t *testing.T) {
	css := readLayoutTestFile(t, "styles.css")
	sunny30Table := betweenLayoutTestMarkers(t, css, `.sunny30-table {`, `}`)

	requireLayoutTestContains(t, sunny30Table, `width: 100%;`)
	requireLayoutTestContains(t, sunny30Table, `min-width: 1180px;`)
	requireLayoutTestNotContains(t, sunny30Table, `width: fit-content;`)
	requireLayoutTestContains(t, css, ".sunny30-table th:nth-child(3),\n.sunny30-table td:nth-child(3) {\n  width: 38%;\n}")
}

func readLayoutTestFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}

func requireLayoutTestContains(t *testing.T, content, needle string) {
	t.Helper()
	if !strings.Contains(content, needle) {
		t.Fatalf("expected content to contain %q", needle)
	}
}

func requireLayoutTestNotContains(t *testing.T, content, needle string) {
	t.Helper()
	if strings.Contains(content, needle) {
		t.Fatalf("expected content not to contain %q", needle)
	}
}

func betweenLayoutTestMarkers(t *testing.T, content, start, end string) string {
	t.Helper()
	startIndex := strings.Index(content, start)
	if startIndex < 0 {
		t.Fatalf("missing start marker %q", start)
	}
	rest := content[startIndex+len(start):]
	endIndex := strings.Index(rest, end)
	if endIndex < 0 {
		t.Fatalf("missing end marker %q", end)
	}
	return rest[:endIndex]
}
