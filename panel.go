package main

import "fmt"

type PanelType int

const (
	text = iota
	menu
)

type Panel struct {
	t, l, w, h       int
	xoffset, yoffset int
	col, row         int
	title            string
	text             []string
	cursor           *Cursor
	line             string
	panelType        PanelType
	border           int
	menuItems        map[string]Panel
}

func (p *Panel) init(t, l, w, h int, title string) {
	var cursor Cursor
	if t < 1 {
		t = 1
	}
	if l < 1 {
		l = 1
	}
	p.cursor = &cursor
	p.cursor.init(l, t)
	p.t, p.l, p.w, p.h = t, l, w, h
	p.border = 0
	p.title = title
	p.col, p.row = 0, 0
	p.text = []string{""}
	p.line = p.text[0]
	p.panelType = text
}

func (p *Panel) draw(t *Terminal) {
	p.drawContent(t)
	if p.border != 0 {
		// Draw top bar
		t.cursor.move(p.l, p.t)
		drawThinCorner("top-left")
		drawHorizontalLine(p.w - len(p.title) - p.border*2)
		fmt.Printf("%v", p.title)
		drawThinCorner("top-right")
		// Draw left bar
		t.cursor.move(p.l, p.t+p.border)
		drawLeftVerticalLine(p.h - p.border*2)
		// Draw right bar
		t.cursor.move(p.l+p.w-p.border, p.t+p.border)
		if p.l+p.w > t.w {
			drawRightVerticalLine(p.h - p.border*2)
		} else {
			drawLeftVerticalLine(p.h - p.border*2)
		}
		// Draw bottom bar
		t.cursor.move(p.l, p.t+p.h-p.border)
		drawThinCorner("bottom-left")
		drawHorizontalLine(p.w - p.border*2)
		drawThinCorner("bottom-right")
	}
}

func (p *Panel) drawContent(t *Terminal) {
	x, y := p.cursor.cx, p.cursor.cy
	active := (t.activePanel.title == p.title)
	for i := 0; i < p.h; i++ {
		p.cursor.move(p.l+p.border, p.t+i+p.border)
		if i > len(p.text)-1 {
			break
		}
		if p.panelType == menu && active == true && i+p.yoffset == p.row {
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