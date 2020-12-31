package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
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

func runProgram(prog program, input []int64) (program, []int64) {
	output := []int64{}
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
			output = append(output, prog.get(param1Mode, 1))
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
	return prog, output
}

func main() {
	part1()
	part2()
}

func draw(view []string) {
	for _, line := range view {
		fmt.Println(line)
	}
}

func view(program program) (program, []string) {
	program, output := runProgram(program, []int64{})
	view := []string{}
	line := ""
	for _, char := range output {
		if char == 10 {
			view = append(view, line)
			line = ""
		} else {
			line += string(rune(char))
		}
	}
	return program, view
}

func get(view []string, x int, y int) string {
	if x < 0 || x >= len(view[y]) || y < 0 || y >= len(view) {
		return ""
	}
	return string(view[y][x])
}

func part1() {
	program := readProgram("17")
	_, view := view(program)
	sum := 0
	for y := 0; y < len(view); y++ {
		for x := 0; x < len(view[y]); x++ {
			if get(view, x, y) == "#" &&
				get(view, x+1, y) == "#" &&
				get(view, x-1, y) == "#" &&
				get(view, x, y+1) == "#" &&
				get(view, x, y-1) == "#" {
				sum += x * y
			}
		}
	}
	fmt.Println(sum)
}

func printASCII(output [][]rune) {
	for _, line := range output {
		for _, char := range line {
			fmt.Printf("%c", char)
		}
		fmt.Println()
	}
}

func transformIntoRunes(output []int64) (result [][]rune) {
	for i := len(output) - 1; i >= 0; i-- {
		if output[i] == 10 {
			output[i] = 0
			output = output[:i]
		} else {
			break
		}
	}
	result = append(result, []rune{})
	for _, char := range output {
		switch char {
		case 10:
			result = append(result, []rune{})
		default:
			result[len(result)-1] = append(result[len(result)-1], rune(int(char)))
		}
	}
	return
}

var directions = []rune{'n', 's', 'e', 'w'}

func changeDirection(from rune, to rune) rune {
	if (from == 'n' && to == 'e') ||
		(from == 'e' && to == 's') ||
		(from == 's' && to == 'w') ||
		(from == 'w' && to == 'n') {
		return 'R'
	}
	return 'L'
}

func getNextDirection(direction rune) (int, int) {
	switch direction {
	case 'n':
		return 0, -1
	case 's':
		return 0, 1
	case 'e':
		return 1, 0
	case 'w':
		return -1, 0
	}
	return 0, 0
}

func oppositeDirection(direction rune) rune {
	switch direction {
	case 'n':
		return 's'
	case 's':
		return 'n'
	case 'e':
		return 'w'
	case 'w':
		return 'e'
	}
	return direction
}

func nextStep(paths [][]rune, x, y int, direction rune) (rune, bool) {
	alternatives := []rune{}
	for _, nextDirection := range directions {
		offsetX, offsetY := getNextDirection(nextDirection)
		if y+offsetY < 0 || y+offsetY >= len(paths) || x+offsetX < 0 || x+offsetX >= len(paths[0]) {
			continue
		}
		if paths[y+offsetY][x+offsetX] == '#' && nextDirection != oppositeDirection(direction) {
			alternatives = append(alternatives, nextDirection)
		}
	}
	if len(alternatives) == 1 {
		return alternatives[0], false
	}
	if len(alternatives) == 0 {
		return direction, true
	}
	return direction, false
}

func findStart(paths [][]rune) (int, int) {
	for y, line := range paths {
		for x, cell := range line {
			if cell == '^' {
				return x, y
			}
		}
	}
	return -1, -1
}

func findCommands(paths [][]rune, x, y int) string {
	steps := 1
	commands := ""
	currentDirection := 'n'
	for {
		nextDirection, theEnd := nextStep(paths, x, y, currentDirection)
		if theEnd {
			commands = fmt.Sprint(commands, ",", steps)
			break
		}
		if nextDirection != currentDirection {
			turn := changeDirection(currentDirection, nextDirection)
			if steps > 1 {
				commands = fmt.Sprint(commands, ",", steps)
			}
			commands = fmt.Sprint(commands, ",", string(turn))
			steps = 1
		} else {
			steps++
		}
		currentDirection = nextDirection
		offsetX, offsetY := getNextDirection(currentDirection)
		x += offsetX
		y += offsetY
	}
	return commands[1:]
}

func generateSubcommands(commands string, start int) []string {
	parts := strings.Split(commands, ",")
	result := []string{}
	for i := len(parts); i >= 2; i -= 2 {
		item := parts[:i]
		result = append(result, strings.Join(item, ","))
	}
	return result
}

var regexStart = regexp.MustCompile(`^,`)
var regexEnd = regexp.MustCompile(`,$`)

func subCommands(commands string, subList []string) ([]string, string) {
	if commands == "" {
		return subList, commands
	}
	parts := generateSubcommands(commands, 0)
	for _, part := range parts {
		processed := strings.ReplaceAll(commands, part, "")
		processed = strings.ReplaceAll(processed, ",,", ",")
		processed = regexStart.ReplaceAllString(processed, "")
		processed = regexEnd.ReplaceAllString(processed, "")
		result, rest := subCommands(processed, append(subList, part))
		if len(result) == 3 &&
			rest == "" &&
			len(result[0]) <= 20 &&
			len(result[1]) <= 20 &&
			len(result[2]) <= 20 {
			return result, rest
		}
	}
	return []string{}, commands
}

func part2() {
	program := readProgram("17")
	_, output := runProgram(program, []int64{})
	paths := transformIntoRunes(output)
	printASCII(paths)
	x, y := findStart(paths)
	commands := findCommands(paths, x, y)
	fmt.Println(commands)
	subCommands, _ := subCommands(commands, []string{})
	subCalls := strings.ReplaceAll(commands, subCommands[0], "A")
	subCalls = strings.ReplaceAll(subCalls, subCommands[1], "B")
	subCalls = strings.ReplaceAll(subCalls, subCommands[2], "C")
	input := []int64{}
	for _, item := range subCalls {
		input = append(input, int64(item))
	}
	input = append(input, 10)
	for _, subC := range subCommands {
		for _, item := range subC {
			input = append(input, int64(item))
		}
		input = append(input, 10)
	}
	input = append(input, int64('n'), 10)
	fmt.Println(input)
	program = readProgram("17/2")
	_, output = runProgram(program, input)
	fmt.Println(output)
}
