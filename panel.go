package main

import "fmt"

type PanelType int
type BorderThickness int

const (
	text = iota
	menu
)

const (
	thin = iota
	thick
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
	borderStyle      BorderThickness
	menuItems        map[string]Panel
	orientation      dir
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
	p.borderStyle = thin
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
		drawCorner("top-left", p.borderStyle)
		drawHorizontalLine(p.w-len(p.title)-p.border*2, p.borderStyle)
		fmt.Printf("%v", p.title)
		drawCorner("top-right", p.borderStyle)
		// Draw left bar
		t.cursor.move(p.l, p.t+p.border)
		drawLeftVerticalLine(p.h-p.border*2, p.borderStyle)
		// Draw right bar
		t.cursor.move(p.l+p.w-p.border, p.t+p.border)
		if p.l+p.w > t.w {
			drawRightVerticalLine(p.h-p.border*2, p.borderStyle)
		} else {
			drawLeftVerticalLine(p.h-p.border*2, p.borderStyle)
		}
		// Draw bottom bar
		t.cursor.move(p.l, p.t+p.h-p.border)
		drawCorner("bottom-left", p.borderStyle)
		drawHorizontalLine(p.w-p.border*2, p.borderStyle)
		drawCorner("bottom-right", p.borderStyle)
	}
}

// func (p *Panel) drawHelp(t *Terminal) {
// 	p.drawContent(t)
// 	if p.border != 0 {
// 		// Draw top bar
// 		t.cursor.move(p.l, p.t)
// 		drawCorner("left-t")
// 		drawHorizontalLine(p.w - len(p.title) - p.border*2)
// 		fmt.Printf("%v", p.title)
// 		drawCorner("top-t")
// 		// Draw left bar
// 		t.cursor.move(p.l, p.t+p.border)
// 		drawLeftVerticalLine(p.h - p.border*2)
// 		// Draw right bar
// 		t.cursor.move(p.l+p.w-p.border, p.t+p.border)
// 		if p.l+p.w > t.w {
// 			drawRightVerticalLine(p.h - p.border*2)
// 		} else {
// 			drawLeftVerticalLine(p.h - p.border*2)
// 		}
// 		// Draw bottom bar
// 		t.cursor.move(p.l, p.t+p.h-p.border)
// 		drawCorner("bottom-left")
// 		drawHorizontalLine(p.w - p.border*2)
// 		drawCorner("bottom-right")
// 	}
// }

func (p *Panel) drawContent(t *Terminal) {
	x, y := p.cursor.cx, p.cursor.cy
	active := (t.activePanel.title == p.title)
	if p.orientation == vertical {
		for i := 0; i < p.h; i++ {
			p.cursor.move(p.l+p.border, p.t+i+p.border)
			if i > len(p.text)-1 {
				break
			}
			if p.panelType == menu && active && i+p.yoffset == p.row {
				fmt.Printf(reverseColors)
				fmt.Print(p.text[i+p.yoffset])
				fmt.Printf(resetColors)
			} else {
				fmt.Print(p.text[i+p.yoffset])
			}
			p.cursor.clearLine()
		}
	} else if p.orientation == horizontal {
		// TBI
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
	// ydif := 0
	// if p.cursor.cy+y <= p.h && p.cursor.cy+y <= len(p.text)-p.yoffset {
	// ydif = p.cursor.cy + y
	// e.moveDocCursor(p.cursor.cx, p.cursor.cy+n)
	// } else if p.cursor.cy+y > p.h && len(p.text) > p.h {
	// p.yoffset += y
	// if p.yoffset+p.h > len(p.text) {
	// p.yoffset = len(p.text) - p.h
	// }
	// ydif = p.cursor.cy
	// e.moveDocCursor(e.cx, e.cy)
	// }
	// p.row = ydif + p.yoffset - p.t
	// p.cursor.cy = ydif
	if p.row+y < 0 {
		p.cursor.cy = 0
		p.row += y
		p.yoffset -= y
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
	} else if p.row+y < p.h-p.border*2 {
		p.row += y
		p.cursor.cy += y
	}
}
