package main

import "strconv"

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
// TODO: Use a simple roguelike to test the library

type Map struct {
	w, h   int
	layout [][]Cell
}

type Entity struct {
	cell   Cell
	x, y   int
	blocks bool
	hp     int
}

func (e Entity) String() string {
	return string(e.cell.icon) + " " + strconv.Itoa(e.x) + " " + strconv.Itoa(e.y)
}

func (m *Map) init(w, h int, player Entity) {
	wall := Cell{0x2588, white, black}
	floor := Cell{0x0020, black, black}
	m.w = w
	m.h = h
	m.layout = make([][]Cell, m.h)
	for y := 0; y < m.h; y++ {
		m.layout[y] = make([]Cell, m.w)
		for x := 0; x < m.w; x++ {
			if y == 0 || y == m.h-1 || x == 0 || x == m.w-1 {
				m.layout[y][x] = wall
			} else if x == player.x && y == player.y {
				m.layout[player.y][player.x] = player.cell
			} else {
				m.layout[y][x] = floor
			}
		}
	}
}

func (m *Map) update(player Entity) {
	for y := 0; y < m.h; y++ {
		for x := 0; x < m.w; x++ {
			if x == player.x && y == player.y {
				m.layout[player.y][player.x] = player.cell
			}
		}
	}

}

func (m *Map) move(g *Panel, p *Entity, dx, dy int) {
	if m.layout[p.y+dy][p.x+dx].icon != block {
		m.layout[p.y][p.x] = Cell{0x0020, black, black}
		p.x += dx
		p.y += dy
	}
}

func main() {
	var term Terminal
	term.init()
	defer term.restore()
	term.cursor.clear()
	topPanel := Panel{}
	topPanel.init(1, 1, term.w, 3, term.strToCells("Info"), 1)
	gamePanel := Panel{}
	gamePanel.init(4, 1, term.w, term.h-4, term.strToCells("Game"), 1)
	term.panels = []Panel{topPanel, gamePanel}
	player := Entity{Cell{0x0040, white, black}, 1, 1, false, 10}
	m := Map{}
	m.init(30, 10, player)

	gamePanel.addContent(m.layout)
	topPanel.addContent([][]Cell{term.strToCells(player.String())})
	for {

		m.update(player)
		topPanel.clear([][]Cell{term.strToCells(player.String())})
		gamePanel.clear(m.layout)

		for _, p := range term.panels {

			p.draw(&term)
		}

		inp, _, _ := term.reader.ReadRune()
		switch inp {
		case ctrlQ:
			return
		case ctrlV:
			term.splitPanel(vertical)
		case ctrlS:
			term.splitPanel(horizontal)
		case 'k':
			m.move(&gamePanel, &player, 0, -1)
		case 'j':
			m.move(&gamePanel, &player, 0, 1)
		case 'h':
			m.move(&gamePanel, &player, -1, 0)
		case 'l':
			m.move(&gamePanel, &player, 1, 0)
		case 'w':
			for i := 0; i < 10; i++ {
				m.layout[5][3+i] = Cell{block, white, black}
			}
		}
	}
}
