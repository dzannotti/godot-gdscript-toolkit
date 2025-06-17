package linter

import (
	"regexp"
	"strings"

	"github.com/dzannotti/gdtoolkit/internal/core/ast"
)

// RuleDirective represents a linting directive found in comments
type RuleDirective struct {
	Type  DirectiveType
	Rules []string
	Line  int
}

// DirectiveType represents the type of linting directive
type DirectiveType int

const (
	DirectiveIgnore DirectiveType = iota
	DirectiveDisable
	DirectiveEnable
)

// ParseDirectives extracts linting directives from source code comments
func ParseDirectives(source string) []RuleDirective {
	var directives []RuleDirective

	// Regular expressions for different directive types
	ignorePattern := regexp.MustCompile(`#\s*gdlint\s*:\s*ignore\s*=\s*([^#\n]+)`)
	disablePattern := regexp.MustCompile(`#\s*gdlint\s*:\s*disable\s*=\s*([^#\n]+)`)
	enablePattern := regexp.MustCompile(`#\s*gdlint\s*:\s*enable\s*=\s*([^#\n]+)`)

	lines := strings.Split(source, "\n")

	for i, line := range lines {
		lineNum := i + 1

		// Check for ignore directive
		if matches := ignorePattern.FindStringSubmatch(line); len(matches) > 1 {
			rules := parseRuleList(matches[1])
			directives = append(directives, RuleDirective{
				Type:  DirectiveIgnore,
				Rules: rules,
				Line:  lineNum,
			})
		}

		// Check for disable directive
		if matches := disablePattern.FindStringSubmatch(line); len(matches) > 1 {
			rules := parseRuleList(matches[1])
			directives = append(directives, RuleDirective{
				Type:  DirectiveDisable,
				Rules: rules,
				Line:  lineNum,
			})
		}

		// Check for enable directive
		if matches := enablePattern.FindStringSubmatch(line); len(matches) > 1 {
			rules := parseRuleList(matches[1])
			directives = append(directives, RuleDirective{
				Type:  DirectiveEnable,
				Rules: rules,
				Line:  lineNum,
			})
		}
	}

	return directives
}

// parseRuleList parses a comma-separated list of rule names
func parseRuleList(ruleStr string) []string {
	var rules []string
	parts := strings.Split(ruleStr, ",")

	for _, part := range parts {
		rule := strings.TrimSpace(part)
		if rule != "" {
			rules = append(rules, rule)
		}
	}

	return rules
}

// RuleContext tracks which rules are enabled/disabled at different positions
type RuleContext struct {
	globallyDisabled map[string]bool
	ignoredAtLine    map[int]map[string]bool
}

// NewRuleContext creates a new rule context
func NewRuleContext() *RuleContext {
	return &RuleContext{
		globallyDisabled: make(map[string]bool),
		ignoredAtLine:    make(map[int]map[string]bool),
	}
}

// ProcessDirectives processes all directives and updates the rule context
func (rc *RuleContext) ProcessDirectives(directives []RuleDirective) {
	for _, directive := range directives {
		switch directive.Type {
		case DirectiveIgnore:
			// Mark rules as ignored for the next line
			nextLine := directive.Line + 1
			if rc.ignoredAtLine[nextLine] == nil {
				rc.ignoredAtLine[nextLine] = make(map[string]bool)
			}
			for _, rule := range directive.Rules {
				rc.ignoredAtLine[nextLine][rule] = true
			}

		case DirectiveDisable:
			// Mark rules as globally disabled
			for _, rule := range directive.Rules {
				rc.globallyDisabled[rule] = true
			}

		case DirectiveEnable:
			// Mark rules as globally enabled (remove from disabled list)
			for _, rule := range directive.Rules {
				delete(rc.globallyDisabled, rule)
			}
		}
	}
}

// IsRuleEnabled checks if a rule is enabled at a specific position
func (rc *RuleContext) IsRuleEnabled(ruleName string, pos ast.Position) bool {
	// Check if rule is globally disabled
	if rc.globallyDisabled[ruleName] {
		return false
	}

	// Check if rule is ignored at this line
	if ignoredRules, exists := rc.ignoredAtLine[pos.Line]; exists {
		if ignoredRules[ruleName] {
			return false
		}
	}

	return true
}
