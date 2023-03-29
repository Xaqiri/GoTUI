package main

// TODO: Add movement between panels with ^(hjkl)
//       Might need to make Terminal.panels a linked list
//       where each panel points left, right, up, and down
// TODO: Add button to create a custom panel rather than splitting current
//       in half
// TODO: Fix drawing logic to stop screen flickering
// TODO: Clean up panel drawing code
// TODO: Copy cursor code from GoEditor

func (t *Terminal) createTestUi() {
	h := createHelpPanel(*t)
	topPanel, midPanel, botPanel := Panel{}, Panel{}, Panel{}
	topPanel.init(1, 1, t.w, 3, "Title Bar")
	topPanel.panelType = menu
	topPanel.text = []string{"Help"}
	topPanel.menuItems = map[string]Panel{topPanel.text[0]: h}
	topPanel.border = 1
	midPanel.init(4, 1, t.w, t.h-6, "Panel 0")
	midPanel.border = 1
	botPanel.init(t.h-2, 1, t.w, 3, "Info")
	botPanel.text = []string{"", "", ""}
	botPanel.border = 1
	t.panels = []Panel{topPanel, midPanel, botPanel}
	t.activePanelIndex = 2
}

func main() {
	help := false
	var term Terminal
	term.init()
	defer term.restore()
	term.createTestUi()
	for {
		term.cursor.hideCursor()

		if help && term.selection == "" {
			term.activePanel = &term.panels[0]
		}
		if !help {
			term.activePanel = &term.panels[1]
			term.selection = ""
			term.cursor.clear()
		}

		// term.activePanel = &term.panels[term.activePanelIndex]
		term.panels[len(term.panels)-1].text[0] = "Panels: " + term.activePanel.title
		term.panels[len(term.panels)-1].text[1] = "Selection: " + term.selection

		for _, p := range term.panels {
			p.draw(&term)
			if term.selection != "" {
				old := term.activePanel
				p := old.menuItems[term.selection]
				term.activePanel = &p
				term.activePanel.draw(&term)

			}
		}
		// if term.selection != "" {
		// old := term.activePanel
		// p := old.menuItems[term.selection]
		// term.activePanel = &p
		// term.activePanel.draw(&term)
		// term.activePanel = old
		// }

		if term.activePanel.panelType != menu {
			term.cursor.showCursor()
		}

		// term.panels[len(term.panels)-1].text[0] = "Cursor X: " + strconv.Itoa(term.activePanel.cursor.cx) + " Cursor Y:" + strconv.Itoa(term.activePanel.cursor.cy) + " Col: " + strconv.Itoa(term.activePanel.col) + " Row: " + strconv.Itoa(term.activePanel.row)
		// botPanel.text[0] = "Width: " + strconv.Itoa(term.activePanel.w) + " Height:" + strconv.Itoa(term.activePanel.h)
		// term.panels[len(term.panels)-1].text[0] += " Y-Offset: " + strconv.Itoa(term.activePanel.yoffset)
		// term.activePanel.line = term.activePanel.text[term.activePanel.row]
		// botPanel.text[0] += " T-Width: " + strconv.Itoa(term.w) + " T-Height: " + strconv.Itoa(term.h)
		term.cursor.move(term.activePanel.cursor.cx, term.activePanel.cursor.cy)
		inp, _, _ := term.reader.ReadRune()
		term.panels[len(term.panels)-1].text[2] = "Key: " + string(inp)

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
			if t.panelType == text {
				t.cursor.move(t.l, t.cursor.cy)
				t.text = append(t.text, "")
				t.col = 0
				t.row++
			} else {
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
	p.w = len(p.text[len(p.text)-1])
	p.row = 0
	p.line = p.text[0]
	return p
}
