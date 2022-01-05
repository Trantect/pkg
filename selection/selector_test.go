package selection

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitTerms(t *testing.T) {
	testcases := map[string][]string{
		// Simple selectors
		`a`:                            {`a`},
		`a=avalue`:                     {`a=avalue`},
		`a=avalue,b=bvalue`:            {`a=avalue`, `b=bvalue`},
		`a=avalue,b==bvalue,c!=cvalue`: {`a=avalue`, `b==bvalue`, `c!=cvalue`},

		// Empty terms
		``:     nil,
		`a=a,`: {`a=a`, ``},
		`,a=a`: {``, `a=a`},

		// Escaped values
		`k=\,,k2=v2`:   {`k=\,`, `k2=v2`},   // escaped comma in value
		`k=\\,k2=v2`:   {`k=\\`, `k2=v2`},   // escaped backslash, unescaped comma
		`k=\\\,,k2=v2`: {`k=\\\,`, `k2=v2`}, // escaped backslash and comma
		`k=\a\b\`:      {`k=\a\b\`},         // non-escape sequences
		`k=\`:          {`k=\`},             // orphan backslash

		// Multi-byte
		`함=수,목=록`: {`함=수`, `목=록`},
	}

	for selector, expectedTerms := range testcases {
		assert.Equal(t, expectedTerms, splitTerms(selector))
	}
}

func TestSplitTerm(t *testing.T) {
	testcases := map[string]struct {
		lhs string
		op  string
		rhs string
		ok  bool
	}{
		// Simple terms
		`a=value`:  {lhs: `a`, op: `=`, rhs: `value`, ok: true},
		`b==value`: {lhs: `b`, op: `==`, rhs: `value`, ok: true},
		`c!=value`: {lhs: `c`, op: `!=`, rhs: `value`, ok: true},

		// Empty or invalid terms
		``:  {lhs: ``, op: ``, rhs: ``, ok: false},
		`a`: {lhs: ``, op: ``, rhs: ``, ok: false},

		// Escaped values
		`k=\,`:          {lhs: `k`, op: `=`, rhs: `\,`, ok: true},
		`k=\=`:          {lhs: `k`, op: `=`, rhs: `\=`, ok: true},
		`k=\\\a\b\=\,\`: {lhs: `k`, op: `=`, rhs: `\\\a\b\=\,\`, ok: true},

		// Multi-byte
		`함=수`: {lhs: `함`, op: `=`, rhs: `수`, ok: true},
	}

	for term, expected := range testcases {
		lhs, op, rhs, ok := splitTerm(term)
		assert.Equal(t, expected.ok, ok)
		assert.Equal(t, expected.lhs, lhs)
		assert.Equal(t, expected.op, op)
		assert.Equal(t, expected.rhs, rhs)
	}
}

func TestEscapeValue(t *testing.T) {
	// map values to their normalized escaped values
	testcases := map[string]string{
		``:      ``,
		`a`:     `a`,
		`=`:     `\=`,
		`,`:     `\,`,
		`\`:     `\\`,
		`\=\,\`: `\\\=\\\,\\`,
	}

	for unescapedValue, escapedValue := range testcases {
		actualEscaped := EscapeValue(unescapedValue)
		assert.Equal(t, actualEscaped, escapedValue)

		actualUnescaped, err := UnescapeValue(escapedValue)
		assert.Nil(t, err)
		assert.Equal(t, actualUnescaped, unescapedValue)
	}

	// test invalid escape sequences
	invalidTestcases := []string{
		`\`,   // orphan slash is invalid
		`\\\`, // orphan slash is invalid
		`\a`,  // unrecognized escape sequence is invalid
	}
	for _, invalidValue := range invalidTestcases {
		_, err := UnescapeValue(invalidValue)
		if _, ok := err.(InvalidEscapeSequence); !ok || err == nil {
			t.Errorf("UnescapeValue(%s): expected invalid escape sequence error, got %#v", invalidValue, err)
		}
	}
}

func TestSelectorParse(t *testing.T) {
	testGoodStrings := []string{
		"x=a,y=b,z=c",
		"",
		"x!=a,y=b",
		`x=a||y\=b`,
		`x=a\=\=b`,
	}
	testBadStrings := []string{
		"x=a||y=b",
		"x==a==b",
		"x=a,b",
		"x in (a)",
		"x in (a,b,c)",
		"x",
	}
	for _, test := range testGoodStrings {
		lq, err := ParseSelector(test)
		assert.Nil(t, err)
		assert.Equal(t, test, lq.String())
	}
	for _, test := range testBadStrings {
		_, err := ParseSelector(test)
		assert.NotNil(t, err)
	}
}

func TestDeterministicParse(t *testing.T) {
	s1, err := ParseSelector("x=a,a=x")
	s2, err2 := ParseSelector("a=x,x=a")
	if err != nil || err2 != nil {
		t.Errorf("Unexpected parse error")
	}
	assert.Nil(t, err)
	assert.Nil(t, err2)
	assert.Equal(t, s1.String(), s2.String())
}

func expectMatch(t *testing.T, selector string, ls Set) {
	lq, err := ParseSelector(selector)
	assert.Nil(t, err)
	assert.True(t, lq.Matches(ls))
}

func expectNoMatch(t *testing.T, selector string, ls Set) {
	lq, err := ParseSelector(selector)
	assert.Nil(t, err)
	assert.False(t, lq.Matches(ls))
}

func TestEverything(t *testing.T) {
	assert.True(t, Everything().Matches(Set{"x": "y"}))
	assert.True(t, Everything().Empty())
}

func TestSelectorMatches(t *testing.T) {
	expectMatch(t, "", Set{"x": "y"})
	expectMatch(t, "x=y", Set{"x": "y"})
	expectMatch(t, "x=y,z=w", Set{"x": "y", "z": "w"})
	expectMatch(t, "x!=y,z!=w", Set{"x": "z", "z": "a"})
	expectMatch(t, "notin=in", Set{"notin": "in"}) // in and notin in exactMatch
	expectNoMatch(t, "x=y", Set{"x": "z"})
	expectNoMatch(t, "x=y,z=w", Set{"x": "w", "z": "w"})
	expectNoMatch(t, "x!=y,z!=w", Set{"x": "z", "z": "w"})

	fieldset := Set{
		"foo":     "bar",
		"baz":     "blah",
		"complex": `=value\,\`,
	}
	expectMatch(t, "foo=bar", fieldset)
	expectMatch(t, "baz=blah", fieldset)
	expectMatch(t, "foo=bar,baz=blah", fieldset)
	expectMatch(t, `foo=bar,baz=blah,complex=\=value\\\,\\`, fieldset)
	expectNoMatch(t, "foo=blah", fieldset)
	expectNoMatch(t, "baz=bar", fieldset)
	expectNoMatch(t, "foo=bar,foobar=bar,baz=blah", fieldset)
}

func TestOneTermEqualSelector(t *testing.T) {
	assert.True(t, OneTermEqualSelector("x", "y").Matches(Set{"x": "y"}))
	assert.False(t, OneTermEqualSelector("x", "y").Matches(Set{"x": "z"}))
}

func expectMatchDirect(t *testing.T, selector, ls Set) {
	assert.True(t, SelectorFromSet(selector).Matches(ls))
}

func expectNoMatchDirect(t *testing.T, selector, ls Set) {
	assert.False(t, SelectorFromSet(selector).Matches(ls))
}

func TestSetMatches(t *testing.T) {
	labelset := Set{
		"foo": "bar",
		"baz": "blah",
	}
	expectMatchDirect(t, Set{}, labelset)
	expectMatchDirect(t, Set{"foo": "bar"}, labelset)
	expectMatchDirect(t, Set{"baz": "blah"}, labelset)
	expectMatchDirect(t, Set{"foo": "bar", "baz": "blah"}, labelset)
	expectNoMatchDirect(t, Set{"foo": "=blah"}, labelset)
	expectNoMatchDirect(t, Set{"baz": "=bar"}, labelset)
	expectNoMatchDirect(t, Set{"foo": "=bar", "foobar": "bar", "baz": "blah"}, labelset)
}

func TestNilMapIsValid(t *testing.T) {
	selector := Set(nil).AsSelector()
	assert.NotNil(t, selector)
	assert.True(t, selector.Empty())
}

func TestSetIsEmpty(t *testing.T) {
	assert.True(t, (Set{}).AsSelector().Empty())
	assert.True(t, (andTerm(nil)).Empty())
	assert.False(t, (&hasTerm{}).Empty())
	assert.False(t, (&notHasTerm{}).Empty())
	assert.True(t, (andTerm{andTerm{}}).Empty())
	assert.False(t, (andTerm{&hasTerm{"a", "b"}}).Empty())
}

func TestRequiresExactMatch(t *testing.T) {
	testCases := map[string]struct {
		S     Selector
		Label string
		Value string
		Found bool
	}{
		"empty set":                 {Set{}.AsSelector(), "test", "", false},
		"empty hasTerm":             {&hasTerm{}, "test", "", false},
		"skipped hasTerm":           {&hasTerm{"a", "b"}, "test", "", false},
		"valid hasTerm":             {&hasTerm{"test", "b"}, "test", "b", true},
		"valid hasTerm no value":    {&hasTerm{"test", ""}, "test", "", true},
		"valid notHasTerm":          {&notHasTerm{"test", "b"}, "test", "", false},
		"valid notHasTerm no value": {&notHasTerm{"test", ""}, "test", "", false},
		"nil andTerm":               {andTerm(nil), "test", "", false},
		"empty andTerm":             {andTerm{}, "test", "", false},
		"nested andTerm":            {andTerm{andTerm{}}, "test", "", false},
		"nested andTerm matches":    {andTerm{&hasTerm{"test", "b"}}, "test", "b", true},
		"andTerm with non-match":    {andTerm{&hasTerm{}, &hasTerm{"test", "b"}}, "test", "b", true},
	}
	for _, v := range testCases {
		value, found := v.S.RequiresExactMatch(v.Label)
		assert.Equal(t, v.Value, value)
		assert.Equal(t, v.Found, found)
	}
}

func TestTransform(t *testing.T) {
	testCases := []struct {
		name      string
		selector  string
		transform func(field, value string) (string, string, error)
		result    string
		isEmpty   bool
	}{
		{
			name:      "empty selector",
			selector:  "",
			transform: func(field, value string) (string, string, error) { return field, value, nil },
			result:    "",
			isEmpty:   true,
		},
		{
			name:      "no-op transform",
			selector:  "a=b,c=d",
			transform: func(field, value string) (string, string, error) { return field, value, nil },
			result:    "a=b,c=d",
			isEmpty:   false,
		},
		{
			name:     "transform one field",
			selector: "a=b,c=d",
			transform: func(field, value string) (string, string, error) {
				if field == "a" {
					return "e", "f", nil
				}
				return field, value, nil
			},
			result:  "e=f,c=d",
			isEmpty: false,
		},
		{
			name:      "remove field to make empty",
			selector:  "a=b",
			transform: func(field, value string) (string, string, error) { return "", "", nil },
			result:    "",
			isEmpty:   true,
		},
		{
			name:     "remove only one field",
			selector: "a=b,c=d,e=f",
			transform: func(field, value string) (string, string, error) {
				if field == "c" {
					return "", "", nil
				}
				return field, value, nil
			},
			result:  "a=b,e=f",
			isEmpty: false,
		},
	}

	for i, tc := range testCases {
		result, err := ParseAndTransformSelector(tc.selector, tc.transform)
		assert.Nilf(t, err, "case [%d]: unexpected error: %v", i, err)
		assert.Equalf(t, tc.isEmpty, result.Empty(), "[%d] expected empty: %t, got: %t", i, tc.isEmpty, result.Empty())
		assert.Equalf(t, tc.result, result.String(), "[%d] expected result: %s, got: %s", i, tc.result, result.String())
	}
}
