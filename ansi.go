package main

import "fmt"

type Cursor struct {
	cx, cy int
	hidden bool
}

func (c *Cursor) init(x, y int) {
	c.move(x, y)
	c.hidden = false
}

func (c *Cursor) hideCursor() {
	fmt.Printf(hide)
}

func (c *Cursor) showCursor() {
	fmt.Printf(show)
}

func (c *Cursor) clear() {
	fmt.Printf("%s[2J", string(esc))
}

func (c *Cursor) clearLine() {
	fmt.Printf("%s[K", string(esc))
}

func (c *Cursor) move(x, y int) {
	c.cx, c.cy = x+1, y+1
	fmt.Printf("%s[%d;%dH", string(esc), y+1, x+1)
}

func (c *Cursor) left(num int) {
	c.move(c.cx-num, c.cy)
	// c.cx -= num
}

func (c *Cursor) down(num int) {
	c.cy += num
}
