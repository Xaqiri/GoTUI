package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/term"
)

type Terminal struct {
	w, h             int
	reader           *bufio.Reader
	writer           *bufio.Writer
	cursor           *Cursor
	panels           []Panel
	activePanel      *Panel
	activePanelIndex int
	initialState     *term.State
	selection        string
}

func (t *Terminal) init() {
	var cursor Cursor

	t.w, t.h, _ = term.GetSize(0)
	t.cursor = &cursor
	t.cursor.init(1, 1)
	t.reader = bufio.NewReader(os.Stdin)
	t.writer = bufio.NewWriter(os.Stdout)
	t.initialState, _ = term.MakeRaw(0)
	t.activePanelIndex = 0
	t.cursor.hideCursor()
}

func (t *Terminal) restore() {
	t.cursor.move(1, 1)
	t.cursor.clear()
	t.cursor.showCursor()
	term.Restore(0, t.initialState)
	fmt.Println("")
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

func (t *Terminal) splitPanel(direction dir) {
	newPanel := Panel{}
	if direction == vertical {
		if t.activePanel.w/2 > 10 {
			var l, w int
			t.cursor.clear()
			if t.activePanel.w%2 == 0 {
				t.activePanel.w /= 2
				l = t.activePanel.l + t.activePanel.w
				w = t.activePanel.w
			} else {
				t.activePanel.w = (t.activePanel.w / 2) + 1
				l = t.activePanel.l + t.activePanel.w
				w = t.activePanel.w - 1
			}
			p := t.activePanel
			newPanel.init(p.t, l, w, p.h, t.strToCells("Panel "+strconv.Itoa(len(t.panels))), p.border)
			t.panels = append(t.panels, newPanel)
			t.activePanelIndex++
		}
	} else if direction == horizontal {
		if t.activePanel.h/2 > 10 {
			var top, h int
			t.cursor.clear()
			if t.activePanel.h%2 == 0 {
				t.activePanel.h /= 2
				top = t.activePanel.t + t.activePanel.h
				h = t.activePanel.h
			} else {
				t.activePanel.h = (t.activePanel.h / 2) + 1
				top = t.activePanel.t + t.activePanel.h
				h = t.activePanel.h - 1
			}
			p := t.activePanel
			newPanel.init(top, p.l, p.w, h, t.strToCells("Panel "+strconv.Itoa(len(t.panels))), p.border)
			t.panels = append(t.panels, newPanel)
			t.activePanelIndex++
		}
	}
}
