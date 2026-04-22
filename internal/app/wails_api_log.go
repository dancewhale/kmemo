package app

import (
	"unicode/utf8"

	"go.uber.org/zap"
)

// truncateRunes shortens large strings for Info logs (e.g. HTML, JSON).
func truncateRunes(s string, maxRunes int) string {
	if maxRunes <= 0 {
		return ""
	}
	if utf8.RuneCountInString(s) <= maxRunes {
		return s
	}
	r := []rune(s)
	return string(r[:maxRunes]) + "…"
}

func zapOptionalString(name string, p *string) zap.Field {
	if p == nil {
		return zap.String(name, "")
	}
	return zap.String(name, *p)
}

func countKnowledgeTreeNodes(roots []*KnowledgeTreeNode) int {
	n := 0
	var walk func(nodes []*KnowledgeTreeNode)
	walk = func(nodes []*KnowledgeTreeNode) {
		for _, node := range nodes {
			if node == nil {
				continue
			}
			n++
			walk(node.Children)
		}
	}
	walk(roots)
	return n
}

func firstStrings(ids []string, n int) []string {
	if n <= 0 || len(ids) == 0 {
		return nil
	}
	if len(ids) <= n {
		out := make([]string, len(ids))
		copy(out, ids)
		return out
	}
	return append([]string(nil), ids[:n]...)
}

func firstRootKnowledgeIDs(roots []*KnowledgeTreeNode, n int) []string {
	if n <= 0 {
		return nil
	}
	out := make([]string, 0, n)
	for _, r := range roots {
		if r == nil {
			continue
		}
		out = append(out, r.ID)
		if len(out) >= n {
			break
		}
	}
	return out
}
