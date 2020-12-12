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

func (p *program) get(mode int64, offset int64) int64 {
	switch mode {
	case 1:
		return p.read(p.pointer + offset)
	case 2:
		return p.read(p.read(p.pointer+offset) + p.relativeBase)
	default:
		return p.read(p.read(p.pointer + offset))
	}
}

func (p *program) set(mode int64, offset int64, value int64) {
	switch mode {
	case 2:
		p.code[p.read(p.pointer+offset)+p.relativeBase] = value
	default:
		p.code[p.read(p.pointer+offset)] = value
	}
}

func runProgram(prog program, input []int64) program {
	for prog.read(prog.pointer) != 99 {
		opcode := prog.read(prog.pointer) % 100
		param1Mode := (prog.read(prog.pointer) / 100) % 10
		param2Mode := (prog.read(prog.pointer) / 1000) % 10
		param3Mode := (prog.read(prog.pointer) / 10000) % 10
		switch opcode {
		case 1: // sum
			prog.set(param3Mode, 3, prog.get(param1Mode, 1)+prog.get(param2Mode, 2))
			prog.pointer += 4
		case 2: // mul
			prog.set(param3Mode, 3, prog.get(param1Mode, 1)*prog.get(param2Mode, 2))
			prog.pointer += 4
		case 3: // input
			prog.set(param1Mode, 1, input[0])
			input = input[1:]
			prog.pointer += 2
		case 4: // ouput
			log.Println(prog.get(param1Mode, 1))
			prog.pointer += 2
		case 5: // is true
			if prog.get(param1Mode, 1) != 0 {
				prog.pointer = prog.get(param2Mode, 2)
			} else {
				prog.pointer += 3
			}
		case 6: // is false
			if prog.get(param1Mode, 1) == 0 {
				prog.pointer = prog.get(param2Mode, 2)
			} else {
				prog.pointer += 3
			}
		case 7: // less than
			if prog.get(param1Mode, 1) < prog.get(param2Mode, 2) {
				prog.set(param3Mode, 3, 1)
			} else {
				prog.set(param3Mode, 3, 0)
			}
			prog.pointer += 4
		case 8: // equals
			if prog.get(param1Mode, 1) == prog.get(param2Mode, 2) {
				prog.set(param3Mode, 3, 1)
			} else {
				prog.set(param3Mode, 3, 0)
			}
			prog.pointer += 4
		case 9: // set relative base
			prog.relativeBase += prog.get(param1Mode, 1)
			prog.pointer += 2
		}
	}
	return prog
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
