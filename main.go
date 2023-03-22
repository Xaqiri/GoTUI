package main

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

type Terminal struct {
	w, h             int
	reader           *bufio.Reader
	writer           *bufio.Writer
	cursor           Cursor
	panels           []Panel
	activePanel      *Panel
	activePanelIndex int
	initialState     *term.State
}

type Panel struct {
	t, l, w, h     int
	cx, cy, offset int
	title          string
	text           []string
}

func (p *Panel) init(t, l, w, h int, title string) {
	p.t, p.l, p.w, p.h = t, l, w, h
	p.title = title
	p.offset = 0
	p.cx, p.cy = l+1, t+1
	p.text = []string{}
}

func (t *Terminal) init() {
	var cursor Cursor
	t.w, t.h, _ = term.GetSize(0)
	t.cursor = cursor
	t.cursor.init(0, 0)
	t.reader = bufio.NewReader(os.Stdin)
	t.writer = bufio.NewWriter(os.Stdout)
	t.initialState, _ = term.MakeRaw(0)
	t.panels = []Panel{}
	t.activePanelIndex = 0
}

func (t *Terminal) restore() {
	t.cursor.showCursor()
	term.Restore(0, t.initialState)
}

func main() {
	var term Terminal
	term.init()
	term.cursor.clear()
	defer term.restore()
	p := Panel{}
	p.init(0, 0, term.w, term.h, "Panel")
	term.panels = []Panel{p}
	for {
		term.activePanel = &term.panels[term.activePanelIndex]
		for _, p := range term.panels {
			p.draw(&term)
		}
		term.cursor.move(term.activePanel.cx, term.activePanel.cy)
		inp, _, _ := term.reader.ReadRune()
		switch inp {
		case ctrlQ:
			term.cursor.move(0, 0)
			term.cursor.clear()
			return
		case ctrlV:
			p := term.activePanel
			if p.w/2 > 10 {
				term.cursor.clear()
				newPanel := Panel{}
				term.activePanel.w = p.w / 2
				newPanel.init(p.t, p.l+p.w, p.w, p.h, "Panel")
				term.panels = append(term.panels, newPanel)
				term.activePanelIndex = len(term.panels) - 1
			}
		case ctrlS:
			p := term.activePanel
			if p.h/2 > 5 {
				term.cursor.clear()
				newPanel := Panel{}
				term.activePanel.h = p.h / 2
				newPanel.init(p.t+p.h, p.l, p.w, p.h, "Panel")
				term.panels = append(term.panels, newPanel)
				term.activePanelIndex++
			}
		case tab:
			term.activePanelIndex++
			if term.activePanelIndex >= len(term.panels) {
				term.activePanelIndex = 0
			}
		case shiftTab:
			term.activePanelIndex--
			if term.activePanelIndex < 0 {
				term.activePanelIndex = len(term.panels) - 1
			}
		default:
			fmt.Printf("%U", inp)
		}
		// fmt.Printf("\033[1E\033[%dG", term.activePanel.l+1)
	}
}

func debug(t *Terminal) {
	t.cursor.move(0, 20)
	for i := 0; i < t.w; i++ {
		if i%10 == 0 {
			fmt.Print(string('*'))
		} else {
			fmt.Print(string('|'))
		}
	}
}

func (p *Panel) draw(t *Terminal) {
	c := &t.cursor
	// Draw top bar
	c.move(p.l, p.t)
	drawThinCorner("top-left")
	fmt.Printf(p.title)
	drawHorizontalLine(p.w - len(p.title) - 2)
	drawThinCorner("top-right")
	// Draw left bar
	c.move(p.l, p.t+1)
	drawLeftVerticalLine(p.h - 2)
	// Draw right bar
	c.move(p.l+p.w-1, p.t+1)
	if p.l+p.w >= t.w {
		drawRightVerticalLine(p.h - 2)
	} else {
		drawLeftVerticalLine(p.h - 2)
	}
	// Draw bottom bar
	c.move(p.l, p.t+p.h-1)
	drawThinCorner("bottom-left")
	drawHorizontalLine(p.w - 2)
	drawThinCorner("bottom-right")
	c.move(p.l+2, p.t+2)
}
