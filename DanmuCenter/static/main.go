package static

import (
	"bufio"
	"embed"
	"log"
)

//go:embed chars.txt
var f embed.FS
var ReplaceDict = make(map[rune]rune)

func init() {
	chars, err := f.Open("chars.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(chars)
	for scanner.Scan() {
		line := []rune(scanner.Text())
		if line[0] == '#' {
			continue
		}
		if len(line) > 1 {
			for _, value := range line[1:] {
				ReplaceDict[value] = line[0]
			}
		}
	}
}
