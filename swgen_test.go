package swgen

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	root := "/Users/larry/diary/programming/"
	t.Log("hello")
	c, err := NewDirPage(root, root)
	if err != nil {
		t.Error(err)
	}
	t.Log(c)
	t.Log(c.ContentHTML())
}
