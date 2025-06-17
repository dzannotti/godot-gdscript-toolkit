package rules

import (
	"github.com/dzannotti/gdtoolkit/internal/core/linter"
)

// GetDefaultRules returns the default set of linting rules
func GetDefaultRules() []linter.Rule {
	var rules []linter.Rule

	// Basic checks
	rules = append(rules, GetDefaultBasicRules()...)

	// Class checks
	rules = append(rules, GetDefaultClassRules()...)

	// TODO: Enable these once basic rules are working correctly
	// rules = append(rules, GetDefaultNameRules()...)
	// rules = append(rules, GetDefaultDesignRules()...)
	// rules = append(rules, GetDefaultFormatRules()...)
	// rules = append(rules, GetDefaultIfReturnRules()...)

	return rules
}

// GetDefaultBasicRules returns the default basic linting rules
func GetDefaultBasicRules() []linter.Rule {
	return []linter.Rule{
		&ExpressionNotAssigned{},
		&UnnecessaryPass{},
		&DuplicatedLoad{},
		&UnusedArgument{},
		&ComparisonWithItself{},
	}
}

// GetDefaultClassRules returns the default class linting rules
func GetDefaultClassRules() []linter.Rule {
	return []linter.Rule{
		&ClassDefinitionsOrder{},
		&SubClassBeforeParentClass{},
	}
}

// GetAllRules returns all available linting rules
func GetAllRules() []linter.Rule {
	return GetDefaultRules()
}

// GetRuleByName returns a rule by its name
func GetRuleByName(name string) linter.Rule {
	rules := GetAllRules()
	for _, rule := range rules {
		if rule.Name() == name {
			return rule
		}
	}
	return nil
}
