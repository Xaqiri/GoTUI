package main

import "fmt"

type Cursor struct {
	cx, cy int
	hidden bool
	style  int
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

func (c *Cursor) savePos() {
	fmt.Printf("\u001b[s")
}

func (c *Cursor) loadPos() {
	fmt.Printf("\u001b[u")
}

func (c *Cursor) clear() {
	fmt.Printf("%s[2J", string(esc))
}

func (c *Cursor) clearLine() {
	fmt.Printf("%s[K", string(esc))
}

func (c *Cursor) move(x, y int) {
	if x <= 1 {
		c.cx = 1
	} else {
		c.cx = x
	}
	if y <= 1 {
		c.cy = 1
	} else {
		c.cy = y
	}
	fmt.Printf("%s[%d;%dH", string(esc), c.cy, c.cx)
}

func (c *Cursor) up(num int) {
	c.move(c.cx, c.cy-num)
}

func (c *Cursor) down(num int) {
	c.move(c.cx, c.cy+num)
}

func (c *Cursor) left(num int) {
	c.move(c.cx-num, c.cy)
}

func (c *Cursor) right(num int) {
	c.move(c.cx+num, c.cy)
}
