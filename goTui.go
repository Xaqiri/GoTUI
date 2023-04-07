package main

import (
	"strconv"
	"strings"
)

// TODO: Add movement between panels with ^(hjkl)
//       Might need to make Terminal.panels a linked list
//       where each panel points left, right, up, and down
// TODO: Add button to create a custom panel rather than splitting current
//       in half
// TODO: Copy cursor code from GoEditor
// TODO: Look up how to implement a timer here
// TODO: Add support for vertical and horizontal menu bars

func activePanel(t Terminal) {
	x, y := t.csr.cx, t.csr.cy
	for i, p := range t.panels {
		if (x >= p.l && x <= p.l+p.w) && (y >= p.t && y <= p.t+p.h) {
			t.panelIndex = i
			t.csr.move(p.l+p.leftBorder, p.t+p.topBorder)
		}
	}
}

func main() {
	var term Terminal
	term.drawCursor.clear()
	term.init()
	defer term.restore()
	help := false
	l := Panel{}
	m := Panel{}
	h := Panel{}

	l.init(1, 1, 6, term.h, term.strToCells("Line"), term.colors.cyan, term.colors.black)
	m.init(1, l.w+l.leftBorder, term.w-l.w, term.h, term.strToCells("Text"), term.colors.white, term.colors.black)
	h.init(term.h-20, 10, 20, 10, term.strToCells("Help"), term.colors.green, term.colors.black)

	h.topBorder, h.rightBorder, h.leftBorder, h.bottomBorder = 1, 1, 1, 1

	l.topBorder = 0
	l.rightBorder = 1
	l.leftBorder = 0
	l.bottomBorder = 0
	for i := 1; i <= l.h-l.topBorder-l.bottomBorder; i++ {
		num := strconv.Itoa(i)
		if i < 10 {
			num = strings.Repeat(" ", l.w-l.leftBorder-l.rightBorder-1) + num
		} else if i < 100 {
			num = strings.Repeat(" ", l.w-l.leftBorder-l.rightBorder-2) + num
		}
		l.addContent(term.strToCells(num))
	}
	h.addContent(term.strToCells("X: " + strconv.Itoa(term.csr.cx) + " Y: " + strconv.Itoa(term.csr.cy)))
	term.panels = []*Panel{&l, &m}

	term.draw()
	for {

		inp, _, _ := term.reader.ReadRune()
		switch inp {
		case ctrlQ:
			return
		case 'h':
			term.csr.left(1)
		case 'j':
			term.csr.down(1)
		case 'k':
			term.csr.up(1)
		case 'l':
			term.csr.right(1)
		case 'i':
			activePanel(term)
		case 'n':
			help = !help
		}

		if help {
			h.content[0] = term.strToCells("X: " + strconv.Itoa(term.csr.cx) + " Y: " + strconv.Itoa(term.csr.cy))
			term.panels = []*Panel{&l, &m, &h}
		} else {
			term.panels = []*Panel{&l, &m}
		}

		term.draw()
	}
}
