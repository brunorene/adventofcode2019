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

func runProgram(program program) program {
	for program.code[program.pointer] != 99 {
		switch program.code[program.pointer] {
		case 1:
			val1 := program.code[program.code[program.pointer+1]]
			val2 := program.code[program.code[program.pointer+2]]
			program.code[program.code[program.pointer+3]] = val1 + val2
			program.pointer += 4
		case 2:
			val1 := program.code[program.code[program.pointer+1]]
			val2 := program.code[program.code[program.pointer+2]]
			program.code[program.code[program.pointer+3]] = val1 * val2
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
	program := readProgram("02")
	program.code[1] = 12
	program.code[2] = 2
	program = runProgram(program)
	log.Println(program.code[0])
}

func part2() {
	for noun := 0; noun < 100; noun++ {
		for verb := 0; verb < 100; verb++ {
			program := readProgram("02")
			program.code[1] = noun
			program.code[2] = verb
			program = runProgram(program)
			if program.code[0] == 19690720 {
				log.Println(noun*100 + verb)
				return
			}
		}
	}
}
