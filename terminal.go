package main

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

type Terminal struct {
	t, l         int
	w, h         int
	reader       *bufio.Reader
	writer       *bufio.Writer
	cursor       *Cursor
	panels       []Panel
	content      [][]Cell
	initialState *term.State
	panelIndex   int
}

func (t *Terminal) init() {
	var cursor Cursor
	t.t, t.l = 1, 1
	t.w, t.h, _ = term.GetSize(0)
	t.cursor = &cursor
	t.cursor.init(t.l, t.t)
	t.reader = bufio.NewReader(os.Stdin)
	t.writer = bufio.NewWriter(os.Stdout)
	t.initialState, _ = term.MakeRaw(0)
	t.cursor.hideCursor()

	t.content = make([][]Cell, t.h)
	for y := 0; y < t.h; y++ {
		t.content[y] = make([]Cell, t.w)
		for x := 0; x < t.w; x++ {
			t.content[y][x] = Cell{block, black, black}
		}
	}
}

func (t *Terminal) restore() {
	term.Restore(0, t.initialState)
	t.cursor.move(0, 0)
	t.cursor.clear()
	t.cursor.showCursor()
	fmt.Println("")
}

func (t *Terminal) clear() {
	for y := 0; y < t.h; y++ {
		for x := 0; x < t.w; x++ {
			t.content[y][x] = Cell{block, black, black}
		}
	}
}

func (t *Terminal) draw() {
	t.cursor.move(t.l, t.t)
	for y := 0; y < t.h; y++ {
		for x := 0; x < t.w; x++ {
			t.content[y][x].draw()
		}
		t.cursor.move(t.l, t.t+y+1)
	}
}

func (t *Terminal) getSize() (int, int, error) {
	return term.GetSize(0)
}

func (t *Terminal) strToCells(str string) []Cell {
	cells := make([]Cell, len(str))
	for i := 0; i < len(str); i++ {
		cells[i] = Cell{int(str[i]), white, black}
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
	for y := 0; y < p.h; y++ {
		for x := 0; x < p.w; x++ {
			t.content[y+p.t-p.border][x+p.l-p.border] = p.visualContent[y][x]
		}
	}
}

// func (t *Terminal) splitPanel(direction dir) {
// 	newPanel := Panel{}
// 	if direction == vertical {
// 		if t.activePanel.w/2 > 10 {
// 			var l, w int
// 			t.cursor.clear()
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
// 			t.cursor.clear()
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
