package main

import "fmt"

const (
	COLOR_BEG = iota
	GREEN
	RED
	YELLOW
	YELLOW_UNDERLINE
	DEFAULT
	COLOR_END
)

var COLOR_MAP = map[int]string{
	GREEN:   				"\033[32m",
	RED:     				"\033[31m",
	YELLOW :				"\033[33m",
	YELLOW_UNDERLINE :		"\033[4;33m",
	DEFAULT: "\033[0m",
}

func PrintColoredText(text string, color int) {
	if color <= COLOR_BEG || color >= COLOR_END {
		panic(fmt.Sprintf("Color %d not exist", color))
	}
	fmt.Printf("%s%s%s", COLOR_MAP[color], text, COLOR_MAP[DEFAULT])
}
