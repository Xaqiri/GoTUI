package main

import "fmt"

type Cell struct {
	icon int
	fg   color
	bg   color
}

func (c *Cell) draw() {
	fmt.Printf("\033[38;5;%dm", c.fg)
	fmt.Printf("\033[48;5;%dm", c.bg)
	fmt.Printf("%c", c.icon)
	fmt.Printf("\033[0m")
}
