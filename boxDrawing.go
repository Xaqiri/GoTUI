package main

import "fmt"

func drawThinCorner(corner string) {
	switch corner {
	case "left-t":
		fmt.Printf("\u251C")
	case "top-left":
		fmt.Printf("\u256D")
	case "right-t":
		fmt.Printf("\u2524")
	case "top-right":
		fmt.Printf("\u256E")
	case "bottom-left":
		fmt.Printf("\u2570")
	case "bottom-right":
		fmt.Printf("\u256F")
	}
}

func drawLeftVerticalLine(length int) {
	for i := 0; i < length; i++ {
		fmt.Printf("\u2502")
		fmt.Printf("%s[B%s[D", string(esc), string(esc))
	}
}

func drawRightVerticalLine(length int) {
	for i := 0; i < length; i++ {
		fmt.Printf("\u2502")
		fmt.Printf("%s[B", string(esc))
	}
}

func drawHorizontalLine(length int) {
	for i := 0; i < length; i++ {
		fmt.Printf("\u2500")
	}
}
