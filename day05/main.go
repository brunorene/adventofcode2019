package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type program struct {
	pointer int
	code    []int
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
	numbers := []int{}
	for _, value := range values {
		number, _ := strconv.Atoi(value)
		numbers = append(numbers, number)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return program{code: numbers}
}

func readWithMode(program program, mode int, offset int) int {
	if mode == 1 { // immediate
		return program.code[program.pointer+offset]
	}
	return program.code[program.code[program.pointer+offset]]
}

func runProgram(program program, input []int) program {
	for program.code[program.pointer] != 99 {
		opcode := program.code[program.pointer] % 100
		param1Mode := (program.code[program.pointer] / 100) % 10
		param2Mode := (program.code[program.pointer] / 1000) % 10
		// param3Mode := (program.code[program.pointer] / 10000) % 10
		switch opcode {
		case 1:
			val1 := readWithMode(program, param1Mode, 1)
			val2 := readWithMode(program, param2Mode, 2)

			program.code[program.code[program.pointer+3]] = val1 + val2
			program.pointer += 4
		case 2:
			val1 := readWithMode(program, param1Mode, 1)
			val2 := readWithMode(program, param2Mode, 2)

			program.code[program.code[program.pointer+3]] = val1 * val2
			program.pointer += 4
		case 3:
			position := program.code[program.pointer+1]
			value := input[0]
			input = input[1:]
			program.code[position] = value
			program.pointer += 2
		case 4:
			position := program.code[program.pointer+1]
			if param1Mode == 1 { // immediate
				log.Println(position)
			} else {
				log.Println(program.code[position])
			}
			program.pointer += 2
		case 5:
			val1 := readWithMode(program, param1Mode, 1)
			val2 := readWithMode(program, param2Mode, 2)

			if val1 != 0 {
				program.pointer = val2
			} else {
				program.pointer += 3
			}
		case 6:
			val1 := readWithMode(program, param1Mode, 1)
			val2 := readWithMode(program, param2Mode, 2)

			if val1 == 0 {
				program.pointer = val2
			} else {
				program.pointer += 3
			}
		case 7:
			val1 := readWithMode(program, param1Mode, 1)
			val2 := readWithMode(program, param2Mode, 2)

			program.code[program.code[program.pointer+3]] = 0
			if val1 < val2 {
				program.code[program.code[program.pointer+3]] = 1
			}
			program.pointer += 4
		case 8:
			val1 := readWithMode(program, param1Mode, 1)
			val2 := readWithMode(program, param2Mode, 2)

			program.code[program.code[program.pointer+3]] = 0
			if val1 == val2 {
				program.code[program.code[program.pointer+3]] = 1
			}
			program.pointer += 4
		}
	}
	return program
}

func main() {
	part1()
	part2()
}

func part1() {
	program := readProgram("05")
	runProgram(program, []int{1})
}

func part2() {
	program := readProgram("05")
	runProgram(program, []int{5})
}
