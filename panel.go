package main

import "fmt"

type Panel struct {
	t, l, w, h       int
	xoffset, yoffset int
	col, row         int
	title            string
	text             []string
	cursor           Cursor
	line             string
	panelType        PanelType
	border           bool
}

func (p *Panel) init(t, l, w, h int, title string) {
	var cursor Cursor
	if t == 0 {
		t = 1
	}
	if l == 0 {
		l = 1
	}
	p.cursor = cursor
	p.cursor.init(l, t)
	p.t, p.l, p.w, p.h = t, l, w, h
	p.title = title
	p.col, p.row = 0, 0
	p.text = []string{""}
	p.line = p.text[0]
	p.panelType = text
}

func (p *Panel) draw(t *Terminal) {
	c := &t.cursor
	padding := 2
	if p.w > t.w {
		p.w = t.w
	}
	if p.h > t.h {
		p.h = t.h
	}
	// Draw top bar
	c.move(p.l, p.t)
	drawThinCorner("top-left")
	drawHorizontalLine(p.w - len(p.title) - padding)
	fmt.Printf("%v", p.title)
	drawThinCorner("top-right")
	// Draw left bar
	c.move(p.l, p.t+1)
	drawLeftVerticalLine(p.h)
	// Draw right bar
	c.move(p.l+p.w-1, p.t+1)
	if p.l+p.w-1 >= t.w {
		drawRightVerticalLine(p.h)
	} else {
		drawLeftVerticalLine(p.h)
	}
	// Draw bottom bar
	c.move(p.l, p.t+p.h+1)
	drawThinCorner("bottom-left")
	drawHorizontalLine(p.w - padding)
	drawThinCorner("bottom-right")
}

func (p *Panel) drawContent() {
	x, y := p.cursor.cx, p.cursor.cy
	for i := 0; i < p.h; i++ {
		p.cursor.move(p.l+1, p.t+i+1)
		if i > len(p.text)-1 {
			break
		}
		if p.panelType == menu && i+p.yoffset == p.row {
			fmt.Printf(reverseColors)
			fmt.Print(p.text[i+p.yoffset])
			fmt.Printf(resetColors)
		} else {
			fmt.Print(p.text[i+p.yoffset])
		}
		p.cursor.clearLine()
	}
	p.cursor.cx, p.cursor.cy = x, y
}

func (p *Panel) updateCursorPosition(x, y int) {
	if p.col+x < 0 {
		p.col = 0
	} else if p.col+x >= p.w {
		p.col = p.w - 1
	} else {
		p.col += x
		p.cursor.cx += x
	}

	if p.row+y < 0 {
		p.cursor.cy = 0
		p.row += y
		p.yoffset += y
		if p.yoffset < 0 {
			p.yoffset = 0
		}
		if p.row < 0 {
			p.row = 0
		}
	}
	if p.row+y >= p.h {
		if p.yoffset < len(p.text)-p.h {
			p.yoffset += y
		}
	} else {
		p.row += y
		p.cursor.cy += y
	}
}
