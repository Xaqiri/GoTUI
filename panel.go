package main

type BorderThickness int

const (
	thin BorderThickness = iota
	thick
)

type Panel struct {
	t, l, w, h int
	title      []Cell
	topBorder,
	bottomBorder,
	leftBorder,
	rightBorder int
	content          [][]Cell
	visualContent    [][]Cell
	fg, bg           color
	cursor           Cursor
	xOffset, yOffset int
}

func (p *Panel) init(t, l, w, h int, title []Cell, fg, bg color) {
	if t < 1 {
		t = 1
	}
	if l < 1 {
		l = 1
	}
	p.t, p.l, p.w, p.h = t, l, w, h
	p.topBorder, p.bottomBorder, p.leftBorder, p.rightBorder = 1, 1, 1, 1
	p.title = title
	p.fg, p.bg = fg, bg
	p.cursor = Cursor{0, 0, false, 0}
	p.visualContent = make([][]Cell, p.h)
	for y := 0; y < p.h; y++ {
		p.visualContent[y] = make([]Cell, p.w)
	}
	p.xOffset, p.yOffset = 0, 0
	p.clear()
}

func (p *Panel) update() {
	// Add fix for content being larger than the panel
	for y := 0; y < len(p.content); y++ {
		for x := 0; x < len(p.content[y]); x++ {
			p.visualContent[y+p.topBorder+p.yOffset][x+p.leftBorder+p.xOffset] = p.content[y][x]
		}
	}
}

func (p *Panel) addContent(content []Cell) {
	p.content = append(p.content, content)
	if len(p.content) > p.h-p.topBorder-p.bottomBorder {
		p.yOffset += len(p.content) - p.h - p.topBorder - p.bottomBorder
	}
}

func (p *Panel) clear() {
	box := map[string]Cell{
		"tlCorner": Cell{tlCorner, p.fg, p.bg},
		"trCorner": Cell{trCorner, p.fg, p.bg},
		"blCorner": Cell{blCorner, p.fg, p.bg},
		"brCorner": Cell{brCorner, p.fg, p.bg},
		"vtLine":   Cell{vtLine, p.fg, p.bg},
		"hzLine":   Cell{hzLine, p.fg, p.bg},
		"space":    Cell{space, p.fg, p.bg},
		"block":    Cell{block, p.fg, p.fg},
	}
	for y := 0; y < p.h; y++ {
		for x := 0; x < p.w; x++ {
			switch y {
			case 0:
				if p.topBorder > 0 {
					p.visualContent[y][x] = box["hzLine"]
					if x == 0 {
						p.visualContent[y][x] = box["tlCorner"]
					} else if x == p.w-1 {
						p.visualContent[y][x] = box["trCorner"]
					} else {
						p.visualContent[y][x] = box["hzLine"]
					}
				} else {
					p.visualContent[y][x] = box["space"]
				}
			case p.h - 1:
				if p.bottomBorder > 0 {
					if x == 0 {
						p.visualContent[y][x] = box["blCorner"]
					} else if x == p.w-1 {
						p.visualContent[y][x] = box["brCorner"]
					} else {
						p.visualContent[y][x] = box["hzLine"]
					}
				} else {
					p.visualContent[y][x] = box["space"]
				}
			default:
				if (p.leftBorder > 0 && x == 0) || (p.rightBorder > 0 && x == p.w-1) {
					p.visualContent[y][x] = box["vtLine"]
				} else {
					p.visualContent[y][x] = box["space"]
				}
			}
		}
	}
	if p.topBorder < 1 {
		if p.leftBorder > 0 {
			p.visualContent[0][0] = box["vtLine"]
		}
		if p.rightBorder > 0 {
			p.visualContent[0][p.w-1] = box["vtLine"]
		}
	} else if p.topBorder > 0 {
		if p.leftBorder < 1 {
			p.visualContent[0][0] = box["hzLine"]
		}
		if p.rightBorder < 1 {
			p.visualContent[0][p.w-1] = box["hzLine"]
		}
	}

	if p.bottomBorder < 1 {
		if p.leftBorder > 0 {
			p.visualContent[p.h-1][0] = box["vtLine"]
		}
		if p.rightBorder > 0 {
			p.visualContent[p.h-1][p.w-1] = box["vtLine"]
		}
	} else if p.bottomBorder > 0 {
		if p.leftBorder < 1 {
			p.visualContent[p.h-1][0] = box["hzLine"]
		}
		if p.rightBorder < 1 {
			p.visualContent[p.h-1][p.w-1] = box["hzLine"]
		}
	}

	if p.topBorder > 0 {
		start := p.w - len(p.title) - p.topBorder
		for i := 0; i < len(p.title); i++ {
			p.visualContent[0][start+i] = p.title[i]
		}
	}
	p.update()
}
