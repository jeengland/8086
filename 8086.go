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

var effectiveAddresses = map[uint8][2]string{
	0b000: {"bx", "si"},
	0b001: {"bx", "di"},
	0b010: {"bp", "si"},
	0b011: {"bp", "di"},
	0b100: {"si", ""},
	0b101: {"di", ""},
	0b110: {"bp", ""},
	0b111: {"bx", ""},
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
	return int(int16(hi)<<8 | int16(lo))
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

	if mod == 0b11 {
		rmD := registers[rm][wide]
		if d == 0b0 {
			fmt.Printf("%s %s, %s\n", instruction.opcode, rmD, regD)
			return 2
		} else {
			fmt.Printf("%s %s, %s\n", instruction.opcode, regD, rmD)
			return 2
		}
	} else if mod == 0b00 {
		rmD := effectiveAddresses[rm]
		var efAd string
		if rmD[1] != "" {
			efAd = fmt.Sprintf("[%s + %s]", rmD[0], rmD[1])
		} else {
			efAd = fmt.Sprintf("[%s]", rmD[0])
		}
		if d == 0b0 {
			fmt.Printf("%s %s, %s\n", instruction.opcode, efAd, regD)
			return 2
		} else {
			fmt.Printf("%s %s, %s\n", instruction.opcode, regD, efAd)
			return 2
		}
	} else if mod == 0b01 {
		b2 := bytes[2]

		d8 := int(b2)

		rmD := effectiveAddresses[rm]
		var efAd string

		if rmD[1] != "" {
			efAd = fmt.Sprintf("[%s + %s + %d]", rmD[0], rmD[1], d8)
		} else {
			efAd = fmt.Sprintf("[%s + %d]", rmD[0], d8)
		}

		if d == 0b0 {
			fmt.Printf("%s %s, %s\n", instruction.opcode, efAd, regD)
			return 3
		} else {
			fmt.Printf("%s %s, %s\n", instruction.opcode, regD, efAd)
			return 3
		}
	} else if mod == 0b10 {
		b2 := bytes[2]
		b3 := bytes[3]

		d16 := get16BitValue(b2, b3)

		rmD := effectiveAddresses[rm]
		var efAd string

		if rmD[1] != "" {
			efAd = fmt.Sprintf("[%s + %s + %d]", rmD[0], rmD[1], d16)
		} else {
			efAd = fmt.Sprintf("[%s + %d]", rmD[0], d16)
		}

		if d == 0b0 {
			fmt.Printf("%s %s, %s\n", instruction.opcode, efAd, regD)
			return 4
		} else {
			fmt.Printf("%s %s, %s\n", instruction.opcode, regD, efAd)
			return 4
		}
	} else {
		// 100010 1 1 - 01 010 110
		fmt.Println(int(b0), int(b1))
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
