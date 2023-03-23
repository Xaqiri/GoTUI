package main

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

// TODO: Add movement between panels with ^(hjkl)
//       Might need to make Terminal.panels a linked list
//       where each panel points left, right, up, and down
// TODO: Add button to create a custom panel rather than splitting current
//       in half

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
	t, l, w, h int
	offset     int
	col, row   int
	title      string
	text       []string
	cursor     Cursor
}

func (p *Panel) init(t, l, w, h int, title string) {
	var cursor Cursor
	p.cursor = cursor
	p.cursor.init(l, t)
	p.t, p.l, p.w, p.h = t, l, w, h
	p.title = title
	p.offset = 0
	p.col, p.row = 0, 0
	p.text = []string{""}
}

func (t *Terminal) init() {
	var cursor Cursor
	t.w, t.h, _ = term.GetSize(0)
	t.cursor = cursor
	t.cursor.init(0, 0)
	t.reader = bufio.NewReader(os.Stdin)
	t.writer = bufio.NewWriter(os.Stdout)
	t.initialState, _ = term.MakeRaw(0)
	p := Panel{}
	p.init(0, 0, t.w, t.h, "Panel")
	t.panels = []Panel{p}
	t.activePanelIndex = 0
}

func (t *Terminal) restore() {
	t.cursor.showCursor()
	term.Restore(0, t.initialState)
}

func main() {
	var term Terminal
	term.init()
	defer term.restore()
	help := false
	for {
		h := createHelpPanel(term)
		term.activePanel = &term.panels[term.activePanelIndex]
		term.cursor.hideCursor()
		term.cursor.clear()
		for _, p := range term.panels {
			p.drawContent()
			p.draw(&term, false)
		}
		if help {
			term.cursor.hideCursor()
			h.drawContent()
			h.draw(&term, true)
		}
		if !help {
			term.cursor.showCursor()
		}
		term.cursor.move(term.activePanel.cursor.cx, term.activePanel.cursor.cy)
		//       if help {
		//          term.cursor.hideCursor()
		//         h := createHelpPanel(term)
		//        h.draw(&term)
		//       h.drawContent()
		//  }
		inp, _, _ := term.reader.ReadRune()
		switch inp {
		case esc:
			help = !help
		case ctrlQ:
			term.cursor.move(0, 0)
			term.cursor.clear()
			return
		case ctrlV:
			p := term.activePanel
			if p.w/2 > len(p.title) {
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
		case ctrlH:
			help = !help
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
		case cr:
			t := term.activePanel
			t.cursor.move(t.l, t.cursor.cy)
			t.text = append(t.text, "")
			t.col = 0
			t.row++
		default:
			if !help {
				p := term.activePanel
				p.text[p.row] += string(inp)
				p.cursor.cx++
				p.col++
			}
		}
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

func createHelpPanel(t Terminal) Panel {
	p := Panel{}
	p.init(0, 0, t.w, 0, "Help")
	p.text = []string{
		"Escape: Close this menu",
		"^Q: Quit",
		"^S: Split panel horizontally",
		"^V: Split panel vertically",
		"^H: Open help",
		"Tab: Move to the next panel",
		"Shift-tab: Move to the previous panel",
	}
	p.t = t.h - len(p.text) - 2
	p.h = len(p.text) + 2
	return p
}

func (p *Panel) draw(t *Terminal, help bool) {
	c := &t.cursor
	// Draw top bar
	c.move(p.l, p.t)
	if help {
		drawThinCorner("left-t")
	} else {
		drawThinCorner("top-left")
	}
	fmt.Printf(p.title)
	drawHorizontalLine(p.w - len(p.title) - 2)
	if help {
		drawThinCorner("right-t")
	} else {
		drawThinCorner("top-right")
	}
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
}

func (p *Panel) drawContent() {
	x, y := p.cursor.cx, p.cursor.cy
	for i, line := range p.text {
		p.cursor.move(p.l+1, p.t+i+1)
		fmt.Print(line)
		p.cursor.clearLine()
	}
	p.cursor.cx, p.cursor.cy = x, y
}
