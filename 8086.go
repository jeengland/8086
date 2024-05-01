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

type Instruction struct {
	opcode    string
	operation func(bytes []byte, instruction Instruction) int
}

var opcodes = map[uint8]Instruction{
	0b100010: {"mov", rmToFromReg},
	0b1011:   {"mov", immediateToReg},
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

	i := 0
	for i < len(dat) {
		var count int
		bytes := getBytes(dat, i)
		b1 := bytes[0]

		opcode4b := b1 & 0b11110000 >> 4
		opcode6b := b1 & 0b11111100 >> 2

		// check 4 bit opcodes
		instruction4b, ok := opcodes[opcode4b]
		if ok {
			count = instruction4b.operation(bytes, instruction4b)
		}

		// check 6 bit opcodes
		instruction6b, ok := opcodes[opcode6b]
		if ok {
			count = instruction6b.operation(bytes, instruction6b)
		}

		i += count
	}
}

func getBytes(data []byte, i int) []byte {
	end := i + 6
	if end > len(data) {
		end = len(data)
	}
	return data[i:end]
}

func get16BitValue(lo byte, hi byte) int {
	return int(uint16(hi)<<8 | uint16(lo))
}

func rmToFromReg(bytes []byte, instruction Instruction) int {
	b0 := bytes[0]
	b1 := bytes[1]

	d := b0 & 0b00000010 >> 1
	w := b0 & 0b00000001

	mod := b1 & 0b11000000 >> 6
	reg := b1 & 0b00111000 >> 3
	rm := b1 & 0b00000111

	wide := int(w)

	regD := registers[reg][wide]
	rmD := registers[rm][wide]

	if mod == 0b11 {
		if d == 0b0 {
			fmt.Printf("%s %s, %s\n", instruction.opcode, rmD, regD)
			return 2
		} else {
			fmt.Printf("%s %s, %s\n", instruction.opcode, regD, rmD)
			return 2
		}
	} else {
		panic(3)
	}
}

func immediateToReg(bytes []byte, instruction Instruction) int {
	var count int
	b0 := bytes[0]
	b1 := bytes[1]

	w := b0 & 0b00001000 >> 3
	reg := b0 & 0b00000111

	wide := int(w)

	regD := registers[reg][wide]

	var value int

	if w == 0 {
		value = int(b1)
		count = 2
	} else {
		b2 := bytes[2]
		value = get16BitValue(b1, b2)
		count = 3
	}

	fmt.Printf("%s %s, %d\n", instruction.opcode, regD, value)
	return count
}
