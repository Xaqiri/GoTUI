package main

import (
	"bufio"
	"os"

	"golang.org/x/term"
)

type Terminal struct {
	t, l               int
	w, h               int
	reader             *bufio.Reader
	writer             *bufio.Writer
	csr                *Cursor
	drawCursor         *Cursor
	panels             []*Panel
	content            [][]Cell
	initialState       *term.State
	panelIndex         int
	colors             Colors
	fg, bg, brFG, brBG color
}

func (t *Terminal) init() {
	var drawCursor Cursor
	var csr Cursor
	t.t, t.l = 1, 1
	t.w, t.h, _ = term.GetSize(0)
	t.drawCursor = &drawCursor
	t.drawCursor.init(t.l, t.t)
	t.csr = &csr
	t.csr.init(t.l, t.t)
	t.reader = bufio.NewReader(os.Stdin)
	t.writer = bufio.NewWriter(os.Stdout)
	t.initialState, _ = term.MakeRaw(0)
	t.drawCursor.hideCursor()
	t.colors.init()
	t.fg = t.colors.white
	t.bg = t.colors.black
	t.brFG = t.colors.brWhite
	t.brBG = t.colors.brBlack

	t.content = make([][]Cell, t.h)
	for y := 0; y < t.h; y++ {
		t.content[y] = make([]Cell, t.w)
		for x := 0; x < t.w; x++ {
			t.content[y][x] = Cell{space, t.fg, t.bg}
		}
	}
}

func (t *Terminal) restore() {
	term.Restore(0, t.initialState)
	t.drawCursor.move(1, 1)
	t.drawCursor.clear()
	t.drawCursor.showCursor()
}

func (t *Terminal) clear() {
	for y := 0; y < t.h; y++ {
		for x := 0; x < t.w; x++ {
			t.content[y][x] = Cell{space, t.fg, t.bg}
		}
	}
}

func (t *Terminal) update() {
	for _, p := range t.panels {
		p.clear()
		t.addPanel(*p)
	}
}

func (t *Terminal) draw() {
	t.drawCursor.hideCursor()
	t.clear()
	t.update()
	t.drawCursor.move(t.l, t.t)
	for y := 0; y < t.h; y++ {
		for x := 0; x < t.w; x++ {
			if y == t.csr.cy-1 {
				t.content[y][x].bg = t.brBG
				if t.content[y][x].icon == block {
					t.content[y][x] = Cell{block, t.brBG, t.brBG}
				}
			}
			t.content[y][x].draw()
		}
		t.drawCursor.move(t.l, t.t+y+1)
	}
	t.drawCursor.move(t.csr.cx, t.csr.cy)
	t.drawCursor.showCursor()
}

func (t *Terminal) getSize() (int, int, error) {
	return term.GetSize(0)
}

func (t *Terminal) strToCells(str string) []Cell {
	cells := make([]Cell, len(str))
	for i := 0; i < len(str); i++ {
		cells[i] = Cell{int(str[i]), t.fg, t.bg}
	}
	return cells
}

func (t *Terminal) addText(x, y int, str string) {
	text := t.strToCells(str)
	for i := 0; i < len(text); i++ {
		t.content[y][x+i] = text[i]
	}
}

func (t *Terminal) addPanel(p Panel) {
	for y := p.t; y < p.t+p.h; y++ {
		for x := p.l; x < p.l+p.w; x++ {
			t.content[y-1][x-1] = p.visualContent[y-p.t][x-p.l]
		}
	}
}

// func (t *Terminal) splitPanel(direction dir) {
// 	newPanel := Panel{}
// 	if direction == vertical {
// 		if t.activePanel.w/2 > 10 {
// 			var l, w int
// 			t.drawCursor.clear()
// 			if t.activePanel.w%2 == 0 {
// 				t.activePanel.w /= 2
// 				l = t.activePanel.l + t.activePanel.w
// 				w = t.activePanel.w
// 			} else {
// 				t.activePanel.w = (t.activePanel.w / 2) + 1
// 				l = t.activePanel.l + t.activePanel.w
// 				w = t.activePanel.w - 1
// 			}
// 			p := t.activePanel
// 			newPanel.init(p.t, l, w, p.h, t.strToCells("Panel "+strconv.Itoa(len(t.panels))), p.border)
// 			t.panels = append(t.panels, newPanel)
// 			t.activePanelIndex++
// 		}
// 	} else if direction == horizontal {
// 		if t.activePanel.h/2 > 10 {
// 			var top, h int
// 			t.drawCursor.clear()
// 			if t.activePanel.h%2 == 0 {
// 				t.activePanel.h /= 2
// 				top = t.activePanel.t + t.activePanel.h
// 				h = t.activePanel.h
// 			} else {
// 				t.activePanel.h = (t.activePanel.h / 2) + 1
// 				top = t.activePanel.t + t.activePanel.h
// 				h = t.activePanel.h - 1
// 			}
// 			p := t.activePanel
// 			newPanel.init(top, p.l, p.w, h, t.strToCells("Panel "+strconv.Itoa(len(t.panels))), p.border)
// 			t.panels = append(t.panels, newPanel)
// 			t.activePanelIndex++
// 		}
// 	}
// }
