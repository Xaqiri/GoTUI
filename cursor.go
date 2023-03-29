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
	c.cx, c.cy = x, y
	fmt.Printf("%s[%d;%dH", string(esc), y, x)
}

func (c *Cursor) left(num int) {
	c.move(c.cx-num, c.cy)
}

func (c *Cursor) down(num int) {
	c.move(c.cx, c.cy+num)
}
