package swgen

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestBasicIgnore(t *testing.T) {
	r := strings.NewReader(`
	# comment
	**~
	.git/**
	`)

	bi := NewBasicIgnore(r)
	cases := []struct {
		expect bool
		path   string
	}{
		{true, ".git/config"},
		{true, ".git/hook/post-commit"},
		{true, "hello.org~"},
		{true, "foo/bar.org~"},
		{false, "index.html"},
	}

	for _, c := range cases {
		assert.Equal(t, c.expect, bi.Ignore(c.path), c)
	}
}
