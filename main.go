package main

// TODO: Add movement between panels with ^(hjkl)
//       Might need to make Terminal.panels a linked list
//       where each panel points left, right, up, and down
// TODO: Add button to create a custom panel rather than splitting current
//       in half
// TODO: Fix drawing logic to stop screen flickering
// TODO: Clean up panel drawing code
// TODO: Copy cursor code from GoEditor
//       Cursor code's a mess right now, terminal cursor should be for
//       drawing everything and the panel cursor should only keep
//       track of the inner row/column
// TODO: Look up how to implement a timer here
// TODO: Add support for vertical and horizontal menu bars

func (t *Terminal) createMenuTestUi() {
	h := createHelpPanel(*t)
	topPanel := Panel{}
	topPanel.init(1, 1, t.w, 3, "Title Bar")
	topPanel.panelType = menu
	topPanel.borderStyle = thick
	topPanel.text = []string{"File", "Help"}
	topPanel.menuItems = map[string]Panel{topPanel.text[0]: h}
	topPanel.border = 1
	topPanel.orientation = horizontal
	midPanel := Panel{}
	midPanel.init(4, 1, t.w, t.h-6, "Panel 0")
	midPanel.border = 1
	midPanel.text = []string{"", ""}
	midPanel.borderStyle = thick
	// botPanel.init(t.h-2, 1, t.w, 3, "Info")
	// botPanel.text = []string{"", "", ""}
	// botPanel.border = 1
	t.panels = []Panel{topPanel, midPanel}
	t.activePanelIndex = 1
	t.activePanel = &t.panels[t.activePanelIndex]
}

func main() {
	help, menuBar := false, false
	var term Terminal
	term.init()
	defer term.restore()
	term.createMenuTestUi()

	h := term.panels[0].menuItems["Help"]

	for {
		term.cursor.hideCursor()
		term.cursor.clear()
		if menuBar && !help {
			term.activePanel = &term.panels[0]
		} else if menuBar && help {
			term.activePanel = &h
		} else if !menuBar && !help {
			term.selection = ""
			// term.activePanel = &term.panels[1]
		}

		// term.panels[2].text[0] = strconv.FormatBool(help) + " " + term.activePanel.title + " " + strconv.Itoa(term.activePanel.row)
		// term.panels[2].text[1] = strconv.Itoa(term.activePanel.row)
		for _, p := range term.panels {
			p.draw(&term)
		}

		if help {
			term.activePanel.draw(&term)
		}

		if term.activePanel.panelType != menu {
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
			term.splitPanel(vertical)
		case ctrlS:
			term.splitPanel(horizontal)
		case ctrlH:
			if !help {
				menuBar = !menuBar
			}
			help = false
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
			if t.panelType == text {
				t.cursor.move(t.l, t.cursor.cy)
				t.text = append(t.text, "")
				t.col = 0
				t.row++
			} else {
				if t.text[t.row] == "Help" {
					help = true
				}
				term.selection = t.text[t.row]
			}
		case del:
			p := term.activePanel
			if p.panelType != menu {
				p.updateCursorPosition(-1, 0)
				p.line = p.line[:len(p.line)-1]
			}
		default:
			if term.activePanel.panelType != menu {
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
				t := term.activePanel
				if inp == 'j' {
					// term.activePanel.updateCursorPosition(0, 1)

					// return
					t.row++
					if t.row > len(t.text)-1 {
						t.row = 0
					}
				}
				if inp == 'k' {
					t.row--
					if t.row < 0 {
						t.row = len(t.text) - 1
					}
				}
			}
		}
	}
}

func createHelpPanel(t Terminal) Panel {
	p := Panel{}
	p.init(3, 1, 0, 0, "Help")
	p.border = 1
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
	p.h = len(p.text) + p.border*2
	p.w = len(p.text[len(p.text)-1]) + 2
	p.row = 0
	p.line = p.text[0]
	return p
}
