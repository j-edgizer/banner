// this annotation keeps the file out of the built binary
///go:build manual_test
/// +build manual_test

package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"unicode"
)

var asciiTranslation = map[rune]rune{
	0x7F: 'Ä',
	0x80: 'Ö',
	0x81: 'Ü',
	0x82: 'ä',
	0x83: 'ö',
	0x84: 'ü',
	0x85: 'ß',
}

type letterDescription struct {
	character   rune
	lines       []string
	description string
}

func main() {
	fontName := flag.String("name", "", "a font name from the figlet database")
	keepLeadingSpaces := flag.Bool("spaces", true, "do not strip leading spaces from letters")

	flag.Parse()

	if *fontName == "" {
		fmt.Println("please specify a font name with the -name switch")

		return
	}

	fontUrl := "http://www.figlet.org/fonts/" + *fontName + ".flf"
	fmt.Println("downloading font from", fontUrl)

	resp, err := http.Get(fontUrl)
	if err != nil {
		fmt.Println("could not download font", err)

		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("error downloading font, status code: %d\n", resp.StatusCode)

		return
	}

	scanner := bufio.NewScanner(resp.Body)

	// did we finish the file header?
	lettersStarted := false

	// did the currently processed letter end?
	letterEnded := true
	letterWidth := 0

	// the current letter
	letter := letterDescription{}
	// letters stored so far
	var letters []letterDescription

	characterAsciiCounter := rune(0x20) // figlet characters start at space character, ASCII 0x20 / 32
	characterAscii := rune(0)

	for scanner.Scan() {
		line := scanner.Text()

		if letterEnded {
			if characterAscii == 0 {
				// at first, all characters follow the ascii numbering
				characterAscii = characterAsciiCounter
			}

			if strings.HasSuffix(line, "@") {

				// a regular character starts
				fmt.Printf("%d %c\n", characterAscii, characterAscii)

				// first line of letter
				letterWidth = len(line) - 1
				letterEnded = false
				lettersStarted = true

				// some ascii characters are not in order
				translatedCharacterAscii, found := asciiTranslation[characterAscii]
				if found {
					fmt.Printf("%d %c -> %d %c\n", characterAscii, characterAscii, translatedCharacterAscii, translatedCharacterAscii)
					characterAscii = translatedCharacterAscii
				}
			} else {
				// lines starting with a number denote an ascii character number definition
				match := regexp.MustCompile(`^(\d+)\s+(.*)$`).FindStringSubmatch(line)
				if len(match) != 0 {
					asciiNumberString := match[1]
					letter.description = match[2]

					asciiNumber, err := strconv.Atoi(strings.TrimSpace(asciiNumberString))
					if err != nil {
						fmt.Printf("found invalid ascii description in line: %s", line)

						return
					}

					fmt.Println(line)
					fmt.Printf("-> %d %c\n", asciiNumber, asciiNumber)

					characterAscii = rune(asciiNumber)
				}

				match = regexp.MustCompile(`^(0x[\da-fA-F]+)\s+(.*)$`).FindStringSubmatch(line)
				if len(match) != 0 {
					asciiHexNumberString := match[1]
					letter.description = match[2]

					asciiNumber, err := strconv.ParseInt(asciiHexNumberString, 0, 32)
					if err != nil {
						fmt.Printf("found invalid ascii description '%s' in line: %s", asciiHexNumberString, line)

						return
					}

					fmt.Println(line)
					fmt.Printf("-> %d %c\n", asciiNumber, asciiNumber)

					characterAscii = rune(asciiNumber)
				}

				continue
			}
		} else if len(line) == letterWidth+2 && strings.HasSuffix(line, "@@") {
			// last line of letter
			letterEnded = true

			// trim second @
			line = line[:len(line)-1]

			letter.character = characterAscii

			letter.lines = append(letter.lines, line)

			// clear lines
			if !*keepLeadingSpaces {
				letter.lines = stripLeadingSpaces(letter.lines)
			}
			letter.lines = stripTrailingDollars(letter.lines)

			exists := false
			for _, l := range letters {
				if l.character == characterAscii {
					exists = true
					break
				}
			}

			if !exists {
				letters = append(letters, letter)
			}

			// clear current letter
			letter = letterDescription{}
			characterAscii = rune(0)

			// continue with next letter
			characterAsciiCounter++

			continue
		}

		// skip file header
		if !lettersStarted {
			continue
		}

		// skip lines not containing character information
		if !strings.HasSuffix(line, "@") {
			continue
		}

		letter.lines = append(letter.lines, line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("could not read downloaded font file", err)

		return
	}

	err = writeFile(letters, *fontName, path.Join(relativePath(), "..", "font_"+*fontName+".go"))
	if err != nil {
		fmt.Println("could not write font Go file", err)

		return
	}
}

func writeFile(letters []letterDescription, fontName string, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create font file: %s", err)
	}

	defer file.Close()

	prefix := `package fonts

var ` + capitalizeFirstLetter(fontName) + ` = Font{
`
	postfix := `}
`

	_, err = fmt.Fprint(file, prefix)
	if err != nil {
		return fmt.Errorf("failed to write prefix to font file: %s", err)
	}

	for _, letter := range letters {
		descriptionComment := letter.description
		if letter.description != "" {
			descriptionComment = "/* " + descriptionComment + " */"
		}

		_, err := fmt.Fprintf(file, "\t'%s': "+descriptionComment+" `\n", escapeForSingleQuote(string(letter.character)))
		if err != nil {
			return fmt.Errorf("failed to write letter to font file: %s", err)
		}

		for _, line := range letter.lines {
			_, err = fmt.Fprintln(file, escapeForBackticks(line))
			if err != nil {
				return fmt.Errorf("failed to write line to font file: %s", err)
			}
		}

		_, err = fmt.Fprint(file, "`,\n")
		if err != nil {
			return fmt.Errorf("failed to write letter end to font file: %s", err)
		}
	}

	_, err = fmt.Fprint(file, postfix)
	if err != nil {
		return fmt.Errorf("failed to write postfix to file: %s", err)
	}

	return nil
}

func capitalizeFirstLetter(str string) string {
	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func escapeForBackticks(input string) string {
	return strings.ReplaceAll(input, "`", "` + \"`\" + `")
}

func escapeForSingleQuote(input string) string {
	input = strings.ReplaceAll(input, `\`, `\\`)
	input = strings.ReplaceAll(input, `'`, `\'`)
	return input
}

func stripLeadingSpaces(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}

	// count the leading spaces in the first line
	leadingSpaces := strings.IndexFunc(lines[0], func(c rune) bool {
		return c != ' '
	})

	// no leading spaces found?
	if leadingSpaces == -1 {
		return lines
	}

	// letters which consist of spaces only should be preserved
	allSpaces := true

	for _, line := range lines {
		currentSpaces := strings.IndexFunc(line, func(c rune) bool {
			return c != ' '
		})

		if currentSpaces == -1 {
			leadingSpaces = 0
			break
		}

		if currentSpaces < leadingSpaces {
			leadingSpaces = currentSpaces
		}

		// line has characters other than spaces
		if allSpaces && !regexp.MustCompile(`^\s*\$?@@?$`).MatchString(line) {
			allSpaces = false
		}
	}

	if allSpaces {
		return lines
	}

	// strip the common leading spaces from all lines
	for i, line := range lines {
		lines[i] = line[leadingSpaces:]
	}

	return lines
}

func stripTrailingDollars(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}

	strippedLines := make([]string, len(lines))

	for i, line := range lines {
		if !strings.HasSuffix(line, "$@") {
			// if there is a line that does not end on a dollar sign, it is not a spacer
			return lines
		}
		strippedLines[i] = line[:len(line)-2] + "@"
	}

	return strippedLines
}

func relativePath() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("No caller information")
	}
	return filepath.Dir(filename)
}
