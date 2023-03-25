package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/term"
)

// TODO: Add movement between panels with ^(hjkl)
//       Might need to make Terminal.panels a linked list
//       where each panel points left, right, up, and down
// TODO: Add button to create a custom panel rather than splitting current
//       in half
// TODO: Fix drawing logic to stop screen flickering
// TODO: Make separate files for Terminal and Panel
// TODO: Copy cursor code from GoEditor

type PanelType int

const (
	text = iota
	menu
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
	t, l, w, h       int
	xoffset, yoffset int
	col, row         int
	title            string
	text             []string
	cursor           Cursor
	line             string
	panelType        PanelType
}

func (p *Panel) init(t, l, w, h int, title string) {
	var cursor Cursor
	p.cursor = cursor
	p.cursor.init(l, t)
	p.t, p.l, p.w, p.h = t, l, w, h
	p.title = title
	p.col, p.row = 0, 0
	p.text = []string{""}
	p.line = p.text[0]
	p.panelType = text
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
	p.init(0, 0, 10, 10, "Panel")
	t.panels = []Panel{p}
	t.activePanelIndex = 0
}

func (t *Terminal) restore() {
	t.cursor.showCursor()
	term.Restore(0, t.initialState)
}

func main() {
	botPanel := Panel{}
	var term Terminal
	term.init()
	defer term.restore()
	help := false
	h := createHelpPanel(term)
	botPanel.init(term.h-3, 0, term.w, 3, "Info")
	for {
		term.activePanel = &term.panels[term.activePanelIndex]
		term.cursor.hideCursor()
		term.cursor.clear()
		for _, p := range term.panels {
			p.drawContent()
			p.draw(&term, false)
		}
		botPanel.text[0] = "Cursor X: " + strconv.Itoa(term.activePanel.cursor.cx) + " Cursor Y:" + strconv.Itoa(term.activePanel.cursor.cy) + " Col: " + strconv.Itoa(term.activePanel.col) + " Row: " + strconv.Itoa(term.activePanel.row)
		term.activePanel.line = term.activePanel.text[term.activePanel.row]
		// botPanel.text[0] = term.activePanel.line
		botPanel.drawContent()
		botPanel.draw(&term, false)
		if help {
			term.cursor.hideCursor()
			h.drawContent()
			h.draw(&term, false)
		}
		if !help {
			term.cursor.showCursor()
		}
		term.cursor.move(term.activePanel.cursor.cx, term.activePanel.cursor.cy)
		inp, _, _ := term.reader.ReadRune()
		switch inp {
		case esc:
			help = false
		case ctrlQ:
			term.cursor.move(0, 0)
			term.cursor.clear()
			return
		case ctrlV:
			splitPanelVertically(&term)
		case ctrlS:
			splitPanelHorizontally(&term)
		case ctrlH:
			help = !help
		case plus:
			term.activePanel.w++
			// term.activePanel.h++
		case minus:
			term.activePanel.w--
			// term.activePanel.h--
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
		case del:
			p := term.activePanel
			p.updateCursorPosition(-1, 0)
			p.line = p.line[:len(p.line)-1]
		default:
			if !help {
				if inp == 'j' {
					term.activePanel.updateCursorPosition(0, 1)
				}
				if inp == 'k' {
					term.activePanel.updateCursorPosition(0, -1)
				}
				if inp == 'l' {
					term.activePanel.updateCursorPosition(1, 0)
				}
				if inp == 'h' {
					term.activePanel.updateCursorPosition(-1, 0)
				}
			} else {
				if inp == 'j' {
					h.row++
					if h.row > len(h.text)-1 {
						h.row = 0
					}
				}
				if inp == 'k' {
					h.row--
					if h.row < 0 {
						h.row = len(h.text) - 1
					}
				}
			}
			// if !help {
			// 	p := term.activePanel
			// 	p.text[p.row] += string(inp)
			// 	p.cursor.cx++
			// 	p.col++
			// } else {
			// 	if inp == 'j' {
			// 		h.t += 1
			// 	} else if inp == 'k' {
			// 		h.t -= 1
			// 	} else if inp == 'h' {
			// 		h.l -= 1
			// 	} else if inp == 'l' {
			// 		h.l += 1
			// 	}
			// }
		}
	}
}

func debug(t *Terminal) {
	t.cursor.move(1, t.h-5)
	fmt.Print("hello")
	// for i := 0; i < t.w; i++ {
	// 	if i%10 == 0 {
	// 		fmt.Print(string('*'))
	// 	} else {
	// 		fmt.Print(string('|'))
	// 	}
	// }
}

func splitPanelVertically(term *Terminal) {
	p := term.activePanel
	if p.w/2 > len(p.title) {
		term.cursor.clear()
		newPanel := Panel{}
		term.activePanel.w = p.w / 2
		newPanel.init(p.t, p.l+p.w, p.w, p.h, "Panel")
		term.panels = append(term.panels, newPanel)
		term.activePanelIndex = len(term.panels) - 1
	}
}

func splitPanelHorizontally(term *Terminal) {
	p := term.activePanel
	if p.h/2 > 5 {
		term.cursor.clear()
		newPanel := Panel{}
		term.activePanel.h = p.h / 2
		newPanel.init(p.t+p.h, p.l, p.w, p.h, "Panel")
		term.panels = append(term.panels, newPanel)
		term.activePanelIndex++
	}
}

func createHelpPanel(t Terminal) Panel {
	p := Panel{}
	p.init(10, 20, t.w/3, 0, "Help")
	p.panelType = menu
	p.text = []string{
		"Escape: Close this menu",
		"^Q: Quit",
		"^S: Split panel horizontally",
		"^V: Split panel vertically",
		"^H: Open help",
		"Tab: Move to the next panel",
		"Shift-tab: Move to the previous panel",
	}
	p.h = len(p.text) + 2
	p.row = 0
	p.line = p.text[0]
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
		if p.panelType == menu && i == p.row {
			fmt.Printf(reverseColors)
			fmt.Print(line)
			fmt.Printf(resetColors)
		} else {
			fmt.Print(line)
		}
		p.cursor.clearLine()
	}
	p.cursor.cx, p.cursor.cy = x, y
}

func (p *Panel) updateCursorPosition(x, y int) {
	if p.col+x < 0 {
		p.col = 0
	} else if p.col+x > p.w-3 {
		p.col = p.w - 3
	} else {
		p.col += x
		p.cursor.cx += x
	}
	if p.row+y < 0 {
		p.row = 0
	} else if p.row+y > p.h-3 {
		p.row = p.h - 3
	} else {
		p.row += y
		p.cursor.cy += y
	}
}
