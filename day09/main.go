package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type program struct {
	pointer      int64
	relativeBase int64
	code         map[int64]int64
}

func readProgram(day string) program {
	path, _ := os.Getwd()
	file, err := os.Open(path + "/day" + day + "/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()
	values := strings.Split(line, ",")
	code := make(map[int64]int64)
	for index, value := range values {
		number, _ := strconv.ParseInt(value, 10, 64)
		code[int64(index)] = number
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return program{code: code}
}

func (p *program) read(address int64) int64 {
	value, exists := p.code[address]
	if !exists {
		return 0
	}
	return value
}

func (p *program) readWithMode(mode int64, offset int64) int64 {
	switch mode {
	case 1:
		return p.read(p.pointer + offset)
	case 2:
		return p.read(p.read(p.pointer+offset) + p.relativeBase)
	default:
		return p.read(p.read(p.pointer + offset))
	}
}

func (p *program) writeWithMode(mode int64, offset int64, value int64) {
	switch mode {
	case 2:
		p.code[p.read(p.pointer+offset)+p.relativeBase] = value
	default:
		p.code[p.read(p.pointer+offset)] = value
	}
}

func runProgram(program program, input []int64) program {
	for program.read(program.pointer) != 99 {
		opcode := program.read(program.pointer) % 100
		param1Mode := (program.read(program.pointer) / 100) % 10
		param2Mode := (program.read(program.pointer) / 1000) % 10
		param3Mode := (program.read(program.pointer) / 10000) % 10
		log.Println("internal", program.read(program.pointer), param3Mode, param2Mode, param1Mode)
		switch opcode {
		case 1: // sum
			val1 := program.readWithMode(param1Mode, 1)
			val2 := program.readWithMode(param2Mode, 2)

			program.writeWithMode(param3Mode, 3, val1+val2)
			program.pointer += 4
		case 2: // mul
			val1 := program.readWithMode(param1Mode, 1)
			val2 := program.readWithMode(param2Mode, 2)

			program.writeWithMode(param3Mode, 3, val1*val2)
			program.pointer += 4
		case 3: // input
			position := program.readWithMode(param1Mode, 1)
			value := input[0]
			input = input[1:]
			if param3Mode == 2 {
				program.code[position+program.relativeBase] = value
			} else {
				program.code[position] = value
			}
			program.pointer += 2
		case 4: // ouput
			position := program.readWithMode(param1Mode, 1)
			switch param1Mode {
			case 1: // immediate
				log.Println(position)
			case 2: // relative
				log.Println(program.read(position + program.relativeBase))
			default:
				log.Println(program.read(position))
			}
			program.pointer += 2
		case 5: // is true
			val1 := program.readWithMode(param1Mode, 1)
			val2 := program.readWithMode(param2Mode, 2)

			if val1 != 0 {
				program.pointer = val2
			} else {
				program.pointer += 3
			}
		case 6: // is false
			val1 := program.readWithMode(param1Mode, 1)
			val2 := program.readWithMode(param2Mode, 2)

			if val1 == 0 {
				program.pointer = val2
			} else {
				program.pointer += 3
			}
		case 7: // less than
			val1 := program.readWithMode(param1Mode, 1)
			val2 := program.readWithMode(param2Mode, 2)

			program.writeWithMode(param3Mode, 3, 0)
			if val1 < val2 {
				program.writeWithMode(param3Mode, 3, 1)
			}
			program.pointer += 4
		case 8: // equals
			val1 := program.readWithMode(param1Mode, 1)
			val2 := program.readWithMode(param2Mode, 2)

			program.writeWithMode(param3Mode, 3, 0)
			if val1 == val2 {
				program.writeWithMode(param3Mode, 3, 1)
			}
			program.pointer += 4
		case 9: // set relative base
			val1 := program.readWithMode(param1Mode, 1)
			program.relativeBase += val1
			program.pointer += 2
		}
	}
	return program
}

func main() {
	part1()
	part2()
}

func part1() {
	program := readProgram("09")
	runProgram(program, []int64{1})
}

func part2() {}
