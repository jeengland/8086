package main

import (
	"flag"
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func decodeOp(b byte) string {
	switch b {
	case 0b100010:
		return "mov"
	}
	panic(1)
}

func decodeReg(r byte, w byte) string {
	wide := w == 1

	switch r {
	case 0b000:
		if wide {
			return "ax"
		} else {
			return "al"
		}
	case 0b001:
		if wide {
			return "cx"
		} else {
			return "cl"
		}
	case 0b010:
		if wide {
			return "dx"
		} else {
			return "dl"
		}
	case 0b011:
		if wide {
			return "bx"
		} else {
			return "bl"
		}
	case 0b100:
		if wide {
			return "sp"
		} else {
			return "ah"
		}
	case 0b101:
		if wide {
			return "bp"
		} else {
			return "ch"
		}
	case 0b110:
		if wide {
			return "si"
		} else {
			return "dh"
		}
	case 0b111:
		if wide {
			return "di"
		} else {
			return "bh"
		}
	}
	panic(2)
}

func main() {
	file := flag.String("file", "", "bytes to process")
	flag.Parse()

	fmt.Printf("bits 16\n\n")

	dat, err := os.ReadFile(*file)

	check(err)

	for i := 0; i < len(dat); i += 2 {
		b1 := dat[i]
		b2 := dat[i+1]

		opcode := b1 & 0b11111100 >> 2
		d := b1 & 0b00000010 >> 1
		w := b1 & 0b00000001

		mod := b2 & 0b11000000 >> 6
		reg := b2 & 0b00111000 >> 3
		rm := b2 & 0b00000111

		opcodeD := decodeOp(opcode)
		regD := decodeReg(reg, w)
		rmD := decodeReg(rm, w)

		if mod == 0b11 {
			if d == 0b0 {
				fmt.Printf("%s %s, %s\n", opcodeD, rmD, regD)
			} else {
				fmt.Printf("%s %s, %s\n", opcodeD, regD, rmD)
			}
		} else {
			panic(3)
		}
	}
}
