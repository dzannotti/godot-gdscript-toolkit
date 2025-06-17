package integration

import (
	"testing"

	"github.com/dzannotti/gdtoolkit/internal/testutil"
)

func TestLinterOnAllValidFiles(t *testing.T) {
	fixtures := testutil.GetTestFixtures()

	// Test valid scripts don't crash the linter
	t.Run("ValidScripts", func(t *testing.T) {
		testutil.TestLinterOnValidFiles(t, fixtures.ValidScripts)
	})
}

// TestBasicLinterRules tests the basic linting rules that we've ported from Python
func TestBasicLinterRules(t *testing.T) {
	// Test expression-not-assigned rule
	t.Run("ExpressionNotAssigned", func(t *testing.T) {
		// Valid cases (should pass)
		testutil.SimpleOKCheck(t, `
func foo():
    var x
    x = 1
`)

		testutil.SimpleOKCheck(t, `
func foo():
    bar()
`)

		testutil.SimpleOKCheck(t, `
func foo():
    x.bar()
`)

		testutil.SimpleOKCheck(t, `
func foo():
    for x in [1]: break
`)

		testutil.SimpleOKCheck(t, `
func foo():
    for x in [1]: continue
`)

		testutil.SimpleOKCheck(t, `
func foo():
    '''docstr'''
`)

		testutil.SimpleOKCheck(t, `
func foo():
    await get_tree().create_timer(2.0).timeout
`)

		// Invalid cases (should fail with expression-not-assigned)
		testutil.SimpleNOKCheck(t, `func foo():
    1 + 1
`, "expression-not-assigned", 2)

		testutil.SimpleNOKCheck(t, `func foo():
    true
`, "expression-not-assigned", 2)

		testutil.SimpleNOKCheck(t, `func foo():
    (true)
`, "expression-not-assigned", 2)
	})

	// Test unnecessary-pass rule
	t.Run("UnnecessaryPass", func(t *testing.T) {
		// Valid cases (should pass)
		testutil.SimpleOKCheck(t, `
func foo():
    pass
`)

		testutil.SimpleOKCheck(t, `
func foo():
    var x = true
    if x:
        pass
`)

		// Invalid cases (should fail with unnecessary-pass)
		testutil.SimpleNOKCheck(t, `func foo():
    pass
    1 + 1
`, "unnecessary-pass", 2, "expression-not-assigned")

		testutil.SimpleNOKCheck(t, `func foo():
    if x: pass; 1+1
`, "unnecessary-pass", 2, "expression-not-assigned")
	})

	// Test duplicated-load rule
	t.Run("DuplicatedLoad", func(t *testing.T) {
		// Valid cases (should pass)
		testutil.SimpleOKCheck(t, `
const B = preload('b')
var A = load('a')
func foo():
    var X = load('c')
    var Y = preload('d')
`)

		// Invalid cases (should fail with duplicated-load)
		testutil.SimpleNOKCheck(t, `
const B = preload('b')
var A = load('a')
func foo():
    var X = load('a')
`, "duplicated-load", 5)

		testutil.SimpleNOKCheck(t, `
const B = preload('b')
var A = load('a')
func foo():
    var X = preload('a')
`, "duplicated-load", 5)
	})

	// Test unused-argument rule
	t.Run("UnusedArgument", func(t *testing.T) {
		// Valid cases (should pass)
		testutil.SimpleOKCheck(t, `
func foo(x):
    print(x)
`)

		testutil.SimpleOKCheck(t, `
func foo(_x):
    pass
`)

		// Invalid cases (should fail with unused-argument)
		testutil.SimpleNOKCheck(t, `
func foo(x):
    pass
`, "unused-argument", 2)
	})

	// Test comparison-with-itself rule
	t.Run("ComparisonWithItself", func(t *testing.T) {
		// Valid cases (should pass)
		testutil.SimpleOKCheck(t, `
func foo():
    var x = 1
    if 1 == x:
        return 1
    return 0
`)

		// Invalid cases (should fail with comparison-with-itself)
		testutil.SimpleNOKCheck(t, `func foo():
    if 1 == 1:
        return 1
    return 0
`, "comparison-with-itself", 2)

		testutil.SimpleNOKCheck(t, `func foo(x):
    if x == x:
        return 1
    return 0
`, "comparison-with-itself", 2)

		testutil.SimpleNOKCheck(t, `func foo():
    if "a" == "a":
        return 1
    return 0
`, "comparison-with-itself", 2)

		testutil.SimpleNOKCheck(t, `func foo():
    if (x + 1) == (x + 1):
        return 1
    return 0
`, "comparison-with-itself", 2)
	})
}
