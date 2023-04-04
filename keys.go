package main

const (
	ctrlA = iota + 1
	ctrlB
	ctrlC
	ctrlD
	ctrlE
	ctrlF
	ctrlG
	ctrlH
	ctrlI
	ctrlJ
	ctrlK
	ctrlL
	ctrlM
	ctrlN
	ctrlO
	ctrlP
	ctrlQ
	ctrlR
	ctrlS
	ctrlT
	ctrlU
	ctrlV
	ctrlW
	ctrlX
	ctrlY
	ctrlZ
)

const (
	esc           = '\u001B'
	tab           = '\u0009'
	shiftTab      = '\u005A'
	at            = '\u0040'
	cr            = '\u000D'
	del           = '\u007F'
	hide          = string(esc) + "[?25l"
	show          = "\u001b[?25h"
	plus          = '\u002B'
	minus         = '\u002D'
	reverseColors = "\u001B[7m"
	resetColors   = "\u001B[m"
)

type dir int

const (
	vertical dir = iota
	horizontal
)
