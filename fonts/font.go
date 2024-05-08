package fonts

import (
	"strings"
)

type Font map[rune]Letter

func (f Font) Get(key rune) Letter {
	letter, found := f[key]
	if found {
		return letter
	}
	return f['?']
}

type Letter string

func (l Letter) String() string {
	return strings.Join(l.Lines(), "\n")
}

func (l Letter) Lines() []string {
	trim := string(l[1 : len(l)-2])
	return strings.Split(trim, "@\n")
}
