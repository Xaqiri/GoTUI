package main

import "fmt"

func drawCorner(corner string, thickness BorderThickness) {
	if thickness == thin {
		switch corner {
		case "top-t":
			fmt.Printf("\u252C")
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
	} else if thickness == thick {
		switch corner {
		case "top-t":
			fmt.Printf("\u2566")
		case "left-t":
			fmt.Printf("\u2560")
		case "top-left":
			fmt.Printf("\u2554")
		case "right-t":
			fmt.Printf("\u2563")
		case "top-right":
			fmt.Printf("\u2557")
		case "bottom-left":
			fmt.Printf("\u255A")
		case "bottom-right":
			fmt.Printf("\u255D")
		}

	}
}

func drawLeftVerticalLine(length int, thickness BorderThickness) {
	for i := 0; i < length; i++ {
		if thickness == thin {

			fmt.Printf("\u2502")
		} else if thickness == thick {
			fmt.Printf("\u2551")

		}
		fmt.Printf("%s[B%s[D", string(esc), string(esc))
	}
}

func drawRightVerticalLine(length int, thickness BorderThickness) {
	for i := 0; i < length; i++ {
		if thickness == thin {
			fmt.Printf("\u2502")
		} else if thickness == thick {

			fmt.Printf("\u2551")
		}
		fmt.Printf("%s[B", string(esc))
	}
}

func drawHorizontalLine(length int, thickness BorderThickness) {
	for i := 0; i < length; i++ {
		if thickness == thin {

			fmt.Printf("\u2500")
		} else if thickness == thick {

			fmt.Printf("\u2550")
		}
	}
}
