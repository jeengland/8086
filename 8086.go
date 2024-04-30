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

var opcodes = map[uint8]string{
	0b100010: "mov",
}

var registers = map[uint8][2]string{
	0b000: {"al", "ax"},
	0b001: {"cl", "cx"},
	0b010: {"dl", "dx"},
	0b011: {"bl", "bx"},
	0b100: {"ah", "sp"},
	0b101: {"ch", "bp"},
	0b110: {"dh", "si"},
	0b111: {"bh", "di"},
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

		wide := int(w)

		opcodeD := opcodes[opcode]
		regD := registers[reg][wide]
		rmD := registers[rm][wide]

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
