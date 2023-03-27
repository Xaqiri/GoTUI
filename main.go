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
// TODO: Copy code for changing modes from GoEditor

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
	p.init(0, 0, t.w, t.h-5, "Panel")
	p.text[0] = "*"
	for i := 1; i < 20; i++ {
		if i%5 == 0 {
			p.text = append(p.text, "*")
		} else {
			p.text = append(p.text, "#")
		}
	}
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
	botPanel.init(term.h-3, 0, term.w, 1, "Info")
	for {
		term.activePanel = &term.panels[term.activePanelIndex]
		term.cursor.hideCursor()
		term.cursor.clear()
		for _, p := range term.panels {
			p.drawContent()
			p.draw(&term, false)
		}
		botPanel.text[0] = "Cursor X: " + strconv.Itoa(term.activePanel.cursor.cx) + " Cursor Y:" + strconv.Itoa(term.activePanel.cursor.cy) + " Col: " + strconv.Itoa(term.activePanel.col) + " Row: " + strconv.Itoa(term.activePanel.row)
		botPanel.text[0] += " Y-Offset: " + strconv.Itoa(term.activePanel.yoffset)
		term.activePanel.line = term.activePanel.text[term.activePanel.row]
		// botPanel.text[0] = term.activePanel.line
		// botPanel.drawContent()
		// botPanel.draw(&term, false)
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
			term.activePanel.h++
			if term.activePanel.yoffset > 0 {
				term.activePanel.yoffset--
			}
		case minus:
			term.activePanel.w--
			term.activePanel.h--
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
}

func splitPanelVertically(term *Terminal) {
	p := term.activePanel
	if p.w/2 > len(p.title) {
		term.cursor.clear()
		newPanel := Panel{}
		if (p.w/2)%2 == 0 {
			term.activePanel.w = p.w / 2
		} else {
			term.activePanel.w = (p.w - 1) / 2
		}
		p = term.activePanel
		newPanel.init(p.t, p.l+p.w+2, p.w, p.h, "Panel")
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
		newPanel.init(p.t+p.h+2, p.l, p.w, p.h, "Panel")
		term.panels = append(term.panels, newPanel)
		term.activePanelIndex++
	}
}

func createHelpPanel(t Terminal) Panel {
	p := Panel{}
	p.init(0, 0, t.w/2, 0, "Help")
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
	p.h = len(p.text)
	p.t = t.h / 3
	if t.w/2 > len(p.text[len(p.text)-1]) {
		p.w = t.w / 2
		p.l = t.w / 3
	} else {
		p.w = len(p.text[len(p.text)-1])
		p.l = 0
	}
	p.row = 0
	p.line = p.text[0]
	return p
}

func (p *Panel) draw(t *Terminal, help bool) {
	c := &t.cursor
	if p.w >= t.w {
		p.w = t.w
	}
	if p.h >= t.h {
		p.h = t.h
	}
	// Draw top bar
	c.move(p.l, p.t)
	if help {
		drawThinCorner("left-t")
	} else {
		drawThinCorner("top-left")
	}
	drawHorizontalLine(p.w - len(p.title) - 2)
	fmt.Printf("%v", p.title)
	if help {
		drawThinCorner("right-t")
	} else {
		drawThinCorner("top-right")
	}
	// Draw left bar
	c.move(p.l, p.t+1)
	drawLeftVerticalLine(p.h)
	// Draw right bar
	c.move(p.l+p.w+1, p.t+1)
	// if p.l+p.w >= t.w {
	// drawRightVerticalLine(p.h)
	// } else {
	drawLeftVerticalLine(p.h)
	// }
	// Draw bottom bar
	c.move(p.l, p.t+p.h+1)
	drawThinCorner("bottom-left")
	drawHorizontalLine(p.w - 1)
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
