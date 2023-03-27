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
// TODO: Clean up panel drawing code
// TODO: Make separate files for Terminal
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

func (t *Terminal) init() {
	var cursor Cursor

	t.w, t.h, _ = term.GetSize(0)
	t.cursor = cursor
	t.cursor.init(1, 1)
	t.reader = bufio.NewReader(os.Stdin)
	t.writer = bufio.NewWriter(os.Stdout)
	t.initialState, _ = term.MakeRaw(0)
	p := Panel{}
	p.init(1, 1, t.w, t.h-5, "Panel")
	p.border = true
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
	botPanel.init(term.h-3, 1, term.w, 1, "Info")
	for {
		term.activePanel = &term.panels[term.activePanelIndex]
		term.cursor.hideCursor()
		term.cursor.clear()
		for _, p := range term.panels {
			p.drawContent()
			if p.border {
				p.draw(&term)
			}
		}
		botPanel.text[0] = "Cursor X: " + strconv.Itoa(term.activePanel.cursor.cx) + " Cursor Y:" + strconv.Itoa(term.activePanel.cursor.cy) + " Col: " + strconv.Itoa(term.activePanel.col) + " Row: " + strconv.Itoa(term.activePanel.row)
		// botPanel.text[0] += " Y-Offset: " + strconv.Itoa(term.activePanel.yoffset)
		term.activePanel.line = term.activePanel.text[term.activePanel.row]
		botPanel.text[0] += " Height: " + strconv.Itoa(term.h) + " Width: " + strconv.Itoa(term.w)
		botPanel.drawContent()
		// botPanel.draw(&term, false)
		if help {
			term.cursor.hideCursor()
			h.drawContent()
			h.draw(&term)
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
	p.init(1, 1, t.w/2, 0, "Help")
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
