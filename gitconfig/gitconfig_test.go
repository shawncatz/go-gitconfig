package gitconfig

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	g, err := New()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("GitConfig: %#v\n", g)

	if g == nil {
		t.Error("wut?")
	}
}
