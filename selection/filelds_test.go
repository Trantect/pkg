package selection

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func matches(t *testing.T, ls Set, want string) {
	if ls.String() != want {
		t.Errorf("Expected '%s', but got '%s'", want, ls.String())
	}
}

func TestSetString(t *testing.T) {
	assert.Equal(t, Set{"x": "y"}.String(), "x=y")
	assert.Equal(t, Set{"foo": "bar"}.String(), "foo=bar")
	assert.Equal(t, Set{"foo": "bar", "baz": "qup"}.String(), "baz=qup,foo=bar")
}

func TestFieldHas(t *testing.T) {
	fieldHasTests := []struct {
		Ls  Fields
		Key string
		Has bool
	}{
		{Set{"x": "y"}, "x", true},
		{Set{"x": ""}, "x", true},
		{Set{"x": "y"}, "foo", false},
	}
	for _, lh := range fieldHasTests {
		assert.Equal(t, lh.Has, lh.Ls.Has(lh.Key))
	}
}

func TestFieldGet(t *testing.T) {
	ls := Set{"x": "y"}
	assert.Equal(t, "y", ls.Get("x"))
}
