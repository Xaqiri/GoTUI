package main

import "fmt"

type Cell struct {
	icon int
	fg   int
	bg   int
}

func (c *Cell) draw() {
	fmt.Printf("\033[38;5;%dm", c.fg)
	fmt.Printf("\033[48;5;%dm", c.bg)
	fmt.Printf("%c", c.icon)
	fmt.Printf("\033[0m")
}
