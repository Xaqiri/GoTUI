package main

import "fmt"

type BorderThickness int

const (
	thin = iota
	thick
)

type Panel struct {
	t, l, w, h int
	title      string
	border     int
	content    [][]*Cell
}

func (p *Panel) init(t, l, w, h int, title string, border int) {
	if t < 1 {
		t = 1
	}
	if l < 1 {
		l = 1
	}
	p.t, p.l, p.w, p.h = t, l, w, h
	p.border = border
	p.title = title

	p.content = make([][]*Cell, p.h)
	for y := 0; y < p.h; y++ {
		p.content[y] = make([]*Cell, p.w)
		for x := 0; x < p.w; x++ {
			if p.border == 1 {
				if y == 0 {
					switch x {
					case 0:
						p.content[y][x] = &Cell{tlCorner, white, black}
					case p.w - 1:
						p.content[y][x] = &Cell{trCorner, white, black}
					default:
						p.content[y][x] = &Cell{hzLine, white, black}
					}
				} else if y == p.h-1 {
					switch x {
					case 0:
						p.content[y][x] = &Cell{blCorner, white, black}
					case p.w - 1:
						p.content[y][x] = &Cell{brCorner, white, black}
					default:
						p.content[y][x] = &Cell{hzLine, white, black}
					}
				} else {
					switch x {
					case 0:
						p.content[y][x] = &Cell{vtLine, white, black}
					case p.w - 1:
						p.content[y][x] = &Cell{vtLine, white, black}
					default:
						p.content[y][x] = &Cell{block, black, white}
					}
				}
			} else {
				p.content[y][x] = &Cell{block, black, white}
			}
		}
	}
	start := p.w - len(p.title) - p.border
	for i := 0; i < len(p.title); i++ {
		p.content[0][start+i] = &Cell{int(p.title[i]), cyan, black}

	}
}

func (p *Panel) addContent(content [][]Cell) {
	// Add fix for content being larger than the panel
	w, h := len(content[0]), len(content)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			p.content[y+p.border][x+p.border] = &content[y][x]
		}
	}
}

func (p *Panel) draw(t *Terminal) {
	t.cursor.move(p.l, p.t)
	for y := 0; y < p.h; y++ {
		for x := 0; x < p.w; x++ {
			p.content[y][x].draw()
		}
		fmt.Printf("\033[1E")
	}
}
