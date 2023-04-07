package main

import "strconv"

type color int
type Colors struct {
	black,
	red,
	green,
	yellow,
	blue,
	magenta,
	cyan,
	white,
	brBlack,
	brRed,
	brGreen,
	brYellow,
	brBlue,
	brMagenta,
	brCyan,
	brWhite color
}

func (c *Colors) init() {
	c.black = 0
	c.red = 1
	c.green = 2
	c.yellow = 3
	c.blue = 4
	c.magenta = 5
	c.cyan = 6
	c.white = 7
	c.brBlack = 8
	c.brRed = 9
	c.brGreen = 10
	c.brYellow = 11
	c.brBlue = 12
	c.brMagenta = 13
	c.brCyan = 14
	c.brWhite = 15
}

func setRGB(color string) []int64 {
	if len(color) != 8 {
		panic("color must be a hex string")
	}

	rgb := []int64{0, 0, 0}
	r := color[2:4]
	g := color[4:6]
	b := color[6:]
	var err error
	rgb[0], err = strconv.ParseInt(r, 10, 64)
	rgb[1], err = strconv.ParseInt(g, 10, 64)
	rgb[2], err = strconv.ParseInt(b, 10, 64)
	if err != nil {
		panic("color not a hex value")
	}
	return rgb
}
