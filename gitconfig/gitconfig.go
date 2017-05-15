package gitconfig

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
)

var (
	reSection = regexp.MustCompile(`\[(?P<name>.+)\]`)
	reSetting = regexp.MustCompile(`\s*(?P<Key>[^=\s]+)\s*=\s*(?P<Value>.*)\s*`)
	reComment = regexp.MustCompile(`\s*#`)
)

type GitConfig map[string]string

func (g GitConfig) Save(to string) error {
	file, err := os.Create(to)
	if err != nil {
		return err
	}
	defer file.Close()

	keys := []string{}
	for k := range g {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	section := ""
	for _, k := range keys {
		list := strings.SplitN(k, ".", 2)

		if section != list[0] {
			fmt.Fprintf(file, "\n[%s]\n", list[0])
			section = list[0]
		}

		fmt.Fprintf(file, "%s = %s\n", list[1], g[k])
	}

	return nil
}

// New returns a GitConfig object containing
// the all of the configurations found in the
// directory structure.
func Combined() (GitConfig, error) {
	g := GitConfig{}

	list := files()
	for i := len(list) - 1; i >= 0; i-- {
		f := list[i]

		err := read(g, f)
		if err != nil {
			return nil, err
		}
	}

	return g, nil
}

// Read returns a GitConfig object of the given file.
func Read(file string) (GitConfig, error) {
	g := GitConfig{}

	err := read(g, file)
	if err != nil {
		return nil, err
	}

	return g, nil
}

// files finds all of the git configuration files
// it walks up the directory structure and looks for
// .git/config or .gitconfig until reaching the home
// directory or root path.
func files() []string {
	p := os.Getenv("PWD")
	out := ascend([]string{}, p)
	return out
}

// read updates the given GitConfig with the configuration
// found in the given file.
func read(g GitConfig, file string) error {
	var section string
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	str := string(bytes)

	lines := strings.Split(str, "\n")

	for _, l := range lines {
		m := reSection.FindStringSubmatch(l)
		if m != nil {
			section = m[1]
			continue
		}

		if reComment.FindStringSubmatch(l) != nil {
			continue
		}

		r := reSetting.FindStringSubmatch(l)
		if r == nil {
			continue
		}

		k := section + "." + r[1]
		v := r[2]

		g[k] = v
	}

	return nil
}

// ascend walks up the directory structure searching for
// git configuration files.
func ascend(list []string, dir string) []string {
	if _, err := os.Stat(dir + "/.git/config"); err == nil {
		list = append(list, dir+"/.git/config")
	}

	if _, err := os.Stat(dir + "/.gitconfig"); err == nil {
		list = append(list, dir+"/.gitconfig")
	}

	if dir == os.Getenv("HOME") || dir == "/" {
		return list
	}

	return ascend(list, path.Dir(dir))
}
