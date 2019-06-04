package swgen

import (
	"time"
)

func Generate(input, output, root string, ignore Ignore) (err error) {
	source := &Source{input, input, ignore}
	dir, err := NewDir(source)
	if err != nil {
		return
	}

	target := &Target{root: output, urlRoot: root}
	navigator, err := dir.generateNavigator(target)
	if err != nil {
		return
	}

	metadata := &Metadata{navigator, time.Now(), ""}
	err = dir.Flush(metadata, target)
	if err != nil {
		return
	}

	return
}
