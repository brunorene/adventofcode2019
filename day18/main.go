package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
)

func main() {
	path, _ := os.Getwd()
	file, err := os.Open(path + "/day18/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	maze := [][]byte{}
	scanner := bufio.NewScanner(file)
	y := 0
	reference := make(map[int]node)
	for scanner.Scan() {
		line := scanner.Text()
		maze = append(maze, []byte{})
		x := 0
		for _, cell := range line {
			maze[len(maze)-1] = append(maze[len(maze)-1], byte(cell))
			if cell == '@' || (cell >= 'a' && cell <= 'z') {
				reference[int(cell)-'a'] = node{int(cell) - 'a', coords{x, y, 0}}
			}
			x++
		}
		y++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	part1(maze, reference)

}

type coords struct {
	x, y, keys int
}

func (c *coords) hasKey(key int) bool {
	return (c.keys & int(math.Pow(2, float64(key)))) > 0
}

func (c *coords) addKey(key int) {
	c.keys = c.keys | int(math.Pow(2, float64(key)))
}

func (c *coords) removeKey(key int) {
	c.keys = c.keys &^ int(math.Pow(2, float64(key)))
}

type node struct {
	name   int
	coords coords
}
type path struct {
	start, end node
	length     int
}

func nextNodes(maze [][]byte, start node, history map[coords]bool) []node {
	result := []node{}
	for _, offsetCoords := range []coords{{0, 1, 0}, {0, -1, 0}, {1, 0, 0}, {-1, 0, 0}} {
		coords := coords{start.coords.x + offsetCoords.x, start.coords.y + offsetCoords.y, start.coords.keys}
		opensDoor := maze[coords.y][coords.x] < 'A' || maze[coords.y][coords.x] > 'Z' || start.coords.hasKey(int(maze[coords.y][coords.x])-'A')
		if opensDoor && maze[coords.y][coords.x] != '#' && !history[coords] {
			result = append(result, node{int(maze[coords.y][coords.x]) - 'a', coords})
		}
	}
	return result
}

func generateGraph(maze [][]byte, start node, history map[coords]bool, level string) []path {
	if start.name >= 0 && start.name <= 25 && !start.coords.hasKey(start.name) {
		return []path{{start, start, 0}}
	}
	nextNodes := nextNodes(maze, start, history)
	if len(nextNodes) == 0 {
		return []path{}
	}
	result := []path{}
	for _, node := range nextNodes {
		history[node.coords] = true
		paths := generateGraph(maze, node, history, level+" ")
		for _, p := range paths {
			p.end.coords.addKey(p.end.name)
			if start.name != p.end.name {
				result = append(result, path{start, p.end, p.length + 1})
			}
		}
	}
	return result
}

const infinity = int(^uint(0) >> 1)

func lowerCostNode(table map[node]int, visited map[node]bool) (node, int) {
	nonVisited := make(map[node]int)
	for n, cost := range table {
		if !visited[n] {
			nonVisited[n] = cost
		}
	}

	minName := node{}
	minCost := infinity
	for name, cost := range nonVisited {
		if cost < minCost {
			minName = name
			minCost = cost
		}
	}
	return minName, minCost
}

type costNode struct {
	node   node
	length int
}

func part1(maze [][]byte, reference map[int]node) {
	queue := []node{reference['@'-'a']}
	foundPaths := make(map[path]bool)
	neighbours := make(map[node][]path)
	for len(queue) > 0 {
		next := queue[0]
		queue = queue[1:]
		paths := generateGraph(maze, next, map[coords]bool{next.coords: true}, "")
		for _, p := range paths {
			if !foundPaths[p] {
				foundPaths[p] = true
				_, exists := neighbours[p.start]
				if !exists {
					neighbours[p.start] = []path{}
				}
				neighbours[p.start] = append(neighbours[p.start], p)
				p.end.coords.addKey(p.end.name)
				queue = append(queue, p.end)
			}
		}
	}
	// for p := range foundPaths {
	// 	fmt.Println(p)
	// }

	// Dijkstra
	costTable := make(map[node]int)
	for p := range foundPaths {
		costTable[p.start] = infinity
		costTable[p.end] = infinity
	}
	costTable[reference['@'-'a']] = 0

	visited := make(map[node]bool)

	for len(visited) != len(costTable) {
		next, _ := lowerCostNode(costTable, visited)

		visited[next] = true

		for _, path := range neighbours[next] {
			distanceNeighbour := path.length + costTable[next]
			if distanceNeighbour < costTable[path.end] {
				costTable[path.end] = distanceNeighbour
			}
		}
	}

	fmt.Println("---//---")
	results := []costNode{}
	for k, v := range costTable {
		results = append(results, costNode{k, v})
	}
	sort.Slice(results, func(i, j int) bool {
		if results[i].node.coords.keys == results[j].node.coords.keys {
			return results[i].length > results[j].length
		}
		return results[i].node.coords.keys < results[j].node.coords.keys
	})
	for _, res := range results {
		fmt.Println(res)
	}
}

func part2() {

}
