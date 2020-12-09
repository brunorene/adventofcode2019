package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type chemical struct {
	quantity    int
	name        string
	ingredients map[string]chemical
}

func createChemical(info string) chemical {
	quantityAndName := strings.Split(info, " ")
	quantity, _ := strconv.Atoi(quantityAndName[0])
	return chemical{quantity: quantity, name: quantityAndName[1], ingredients: make(map[string]chemical)}
}

func main() {
	path, _ := os.Getwd()
	file, err := os.Open(path + "/day14/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	chemicals := make(map[string]chemical)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " => ")
		ingredients := strings.Split(parts[0], ", ")
		current, exists := chemicals[createChemical(parts[1]).name]
		if !exists {
			current = createChemical(parts[1])
			chemicals[current.name] = current
		}
		for _, ingredient := range ingredients {
			child := createChemical(ingredient)
			current.ingredients[child.name] = child
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	part1(chemicals)
	part2(chemicals)
}

func oreNeeded(expected chemical, chemicals map[string]chemical, leftovers map[string]int) (int, int) {
	sum := 0
	multiplier := int(math.Ceil(float64(expected.quantity) / float64(chemicals[expected.name].quantity)))
	for _, ingredient := range chemicals[expected.name].ingredients {
		if ingredient.name == "ORE" {
			sum += multiplier * ingredient.quantity
		} else {
			ingredient.quantity *= multiplier
			chemicalNeeded := ingredient.quantity
			leftover, exists := leftovers[ingredient.name]
			if exists {
				if leftover > chemicalNeeded {
					chemicalNeeded = 0
					leftovers[ingredient.name] = leftover - ingredient.quantity

				} else {
					chemicalNeeded -= leftover
					delete(leftovers, ingredient.name)
				}
				ingredient.quantity = chemicalNeeded
			}
			chemicalGenerated, oreQuantity := oreNeeded(ingredient, chemicals, leftovers)
			sum += oreQuantity
			leftovers[ingredient.name] += chemicalGenerated - chemicalNeeded
			if leftovers[ingredient.name] == 0 {
				delete(leftovers, ingredient.name)
			}
		}
	}
	return chemicals[expected.name].quantity * multiplier, sum
}

func part1(chemicals map[string]chemical) {
	log.Println(oreNeeded(chemicals["FUEL"], chemicals, make(map[string]int)))
}

func part2(chemicals map[string]chemical) {
	ore := 0
	wanted := 2
	for ore < 1000000000000 {
		chem := createChemical(fmt.Sprint(wanted, " FUEL"))
		_, ore = oreNeeded(chem, chemicals, make(map[string]int))
		wanted += 100
	}
	for ore > 1000000000000 {
		chem := createChemical(fmt.Sprint(wanted, " FUEL"))
		_, ore = oreNeeded(chem, chemicals, make(map[string]int))
		log.Println(wanted, ore)
		wanted--
	}
}
