package gitconfig

import (
	"testing"
)

func TestNew(t *testing.T) {
	g, err := Combined()
	if err != nil {
		t.Error(err)
	}

	//fmt.Printf("GitConfig: %#v\n", g)

	if g == nil {
		t.Error("wut?")
	}
}

func TestRead(t *testing.T) {
	g, err := Read("fixtures/gitconfig")
	if err != nil {
		t.Error(err)
	}

	if g["alias.st"] != "status" {
		t.Errorf("alias.st == %s, want %s", g["alias.st"], "status")
	}
}

func TestGitConfig_Save(t *testing.T) {
	g, err := Read("fixtures/gitconfig")
	if err != nil {
		t.Error(err)
	}

	err = g.Save("fixtures/test")
	if err != nil {
		t.Error(err)
	}

	g2, err := Read("fixtures/test")
	if err != nil {
		t.Error(err)
	}

	table := []string{
		"alias.st",
		"bitbucket.keychain",
		"bitbucket.user",
		"bitbucket.url",
	}
	for _, s := range table {
		if g2[s] != g[s] {
			t.Errorf("g2 %s == %s, want %s", s, g2[s], g[s])
		}
	}
}
