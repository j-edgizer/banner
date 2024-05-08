package banner

import (
	"strings"

	"github.com/j-edgizer/banner/fonts"
)

func Inline(input string, font fonts.Font) string {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return ""
	}

	lines := font.Get(rune(input[0])).Lines()
	height := len(lines)
	if len(input) > 1 {
		for _, r := range input[1:] {
			switch r {
			case ' ':
				for i := 0; i < height; i++ {
					lines[i] += "  "
				}
			default:
				letter := font.Get(r).Lines()
				for i := 0; i < height; i++ {
					lines[i] += letter[i]
				}
			}
		}
	}

	for i := 0; i < height; i++ {
		lines[i] = strings.TrimRight(lines[i], " ")
	}
	if lines[height-1] == "" {
		lines = lines[:height-1]
	}
	if lines[0] == "" {
		lines = lines[1:]
	}
	return strings.Join(lines, "\n")
}
