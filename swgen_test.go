package swgen

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	root := "/home/i353434/diary/programming/"
	t.Log("hello")
	c, err := NewDirNode(root, root)
	if err != nil {
		t.Error(err)
	}
	t.Log(c)
}
