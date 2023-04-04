package main

// TODO: Add movement between panels with ^(hjkl)
//       Might need to make Terminal.panels a linked list
//       where each panel points left, right, up, and down
// TODO: Add button to create a custom panel rather than splitting current
//       in half
// TODO: Copy cursor code from GoEditor
// TODO: Look up how to implement a timer here
// TODO: Add support for vertical and horizontal menu bars

type Entity struct {
	c      Cell
	size   int
	cursor Cursor
}

func activePanel(t Terminal, e *Entity) {
	x, y := e.cursor.cx, e.cursor.cy
	for i, p := range t.panels {
		if (x >= p.l && x <= p.l+p.w) && (y >= p.t && y <= p.t+p.h) {
			t.panelIndex = i
			e.cursor.move(p.l+p.border, p.t+p.border)
		}
	}
}

func main() {
	var term Terminal
	term.cursor.clear()
	term.init()
	defer term.restore()
	newPanel := false
	t := Panel{}
	m := Panel{}
	r := Panel{}
	n := Panel{}

	t.init(1, 1, term.w, 3, term.strToCells("Panel"), 1)
	m.init(4, 1, term.w-10, term.h-3, term.strToCells("Middle"), 1)
	r.init(4, m.l+m.w, 10, term.h-3, term.strToCells("Right"), 1)
	n.init(m.t+10, m.l+5, 10, 5, term.strToCells("New"), 1)

	e := Entity{Cell{block, cyan, white}, 5, Cursor{1, 1, false}}
	for {
		if newPanel {
			term.panels = []Panel{t, m, r, n}
		} else {
			term.panels = []Panel{t, m, r}
		}

		for _, p := range term.panels {
			term.addPanel(p)
		}
		if newPanel {
			term.addPanel(n)
		}

		term.content[e.cursor.cy-1][e.cursor.cx-1] = e.c
		term.draw()

		inp, _, _ := term.reader.ReadRune()
		switch inp {
		case ctrlQ:
			return
		case 'h':
			e.cursor.left(1)
		case 'j':
			e.cursor.down(1)
		case 'k':
			e.cursor.up(1)
		case 'l':
			e.cursor.right(1)
		case 'i':
			activePanel(term, &e)
		case 'n':
			newPanel = !newPanel
			// case '-':
			// 	p.size -= 1
			// case '+':
			// 	p.size += 1
			// case '\u0020':
			// 	for y := p.cursor.cy; y < p.cursor.cy+p.size; y++ {
			// 		for x := p.cursor.cx; x < p.cursor.cx+p.size; x++ {
			// 			term.content[y][x] = Cell{block, white, white}
			// 		}
			// 	}
		}
	}
}
