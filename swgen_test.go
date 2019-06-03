package swgen

import (
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	root := "/home/i353434/diary/programming/"
	t.Log("hello")
	ignore := NewBasicIgnore(strings.NewReader(".git/**"))
	c, err := NewDirNode(root, root, ignore)
	if err != nil {
		t.Error(err)
	}
	t.Log(c)
}
