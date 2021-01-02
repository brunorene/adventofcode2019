package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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
	reference := make(map[byte]node)
	for scanner.Scan() {
		line := scanner.Text()
		maze = append(maze, []byte{})
		x := 0
		for _, cell := range line {
			maze[len(maze)-1] = append(maze[len(maze)-1], byte(cell))
			if cell == '@' || (cell >= 'a' && cell <= 'z') {
				reference[byte(cell)] = node{string(cell), coords{x, y}}
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
	x, y int
}

type node struct {
	name   string
	coords coords
}

type edge struct {
	start, end node
}
type path struct {
	edge   edge
	length int
}

func nextNodes(maze [][]byte, start node, history map[coords]bool, keys map[byte]bool) []node {
	result := []node{}
	for _, offsetCoords := range []coords{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
		coords := coords{start.coords.x + offsetCoords.x, start.coords.y + offsetCoords.y}
		hasKeyForDoor := maze[coords.y][coords.x] < 65 ||
			maze[coords.y][coords.x] > 90 ||
			keys[maze[coords.y][coords.x]+('a'-'A')]
		if hasKeyForDoor && maze[coords.y][coords.x] != '#' && !history[coords] {
			result = append(result, node{string(maze[coords.y][coords.x]), coords})
		}
	}
	return result
}

func copyMapCoords(m map[coords]bool) map[coords]bool {
	copy := make(map[coords]bool)
	for k, v := range m {
		copy[k] = v
	}
	return copy
}

func copyMapByte(m map[byte]bool) map[byte]bool {
	copy := make(map[byte]bool)
	for k, v := range m {
		copy[k] = v
	}
	return copy
}

func generateGraph(maze [][]byte, start node, history map[coords]bool, keys map[byte]bool, level string) []path {
	if start.name[0] >= 'a' && start.name[0] <= 'z' && !keys[start.name[0]] {
		return []path{{edge{start, start}, 0}}
	}
	nextNodes := nextNodes(maze, start, history, keys)
	if len(nextNodes) == 0 {
		return []path{}
	}
	result := []path{}
	for _, node := range nextNodes {
		newHistory := copyMapCoords(history)
		newHistory[node.coords] = true
		paths := generateGraph(maze, node, newHistory, keys, level+" ")
		for _, p := range paths {
			if start.name != p.edge.end.name {
				result = append(result, path{edge{start, p.edge.end}, p.length + 1})
			}
		}
	}
	return result
}

type queueItem struct {
	node node
	keys map[byte]bool
}

func minPaths(paths []path) []path {
	mem := make(map[edge]int)
	for _, p := range paths {
		length, exists := mem[p.edge]
		if !exists {
			mem[p.edge] = p.length
		} else if length > p.length {
			mem[p.edge] = p.length
		}
	}
	result := []path{}
	for edge, length := range mem {
		result = append(result, path{edge, length})
	}
	return result
}

func lowerCostNode(table map[string]int, visited map[string]bool) (string, int) {
	nonVisited := make(map[string]int)
	for name, cost := range table {
		if !visited[name] {
			nonVisited[name] = cost
		}
	}

	minName := "-"
	minCost := infinity
	for name, cost := range nonVisited {
		if cost < minCost {
			minName = name
			minCost = cost
		}
	}
	return minName, minCost
}

const infinity = int(^uint(0) >> 1)

func part1(maze [][]byte, reference map[byte]node) {
	queue := []queueItem{{reference['@'], make(map[byte]bool)}}
	foundPaths := make(map[path]bool)
	neighbours := make(map[byte][]path)
	for len(queue) > 0 {
		next := queue[0]
		queue = queue[1:]
		paths := minPaths(generateGraph(maze, next.node, map[coords]bool{next.node.coords: true}, next.keys, ""))
		for _, p := range paths {
			if !foundPaths[p] {
				foundPaths[p] = true
				_, exists := neighbours[p.edge.start.name[0]]
				if !exists {
					neighbours[p.edge.start.name[0]] = []path{}
				}
				neighbours[p.edge.start.name[0]] = append(neighbours[p.edge.start.name[0]], p)
				keys := copyMapByte(next.keys)
				keys[p.edge.end.name[0]] = true
				queue = append(queue, queueItem{p.edge.end, keys})
			}
		}
	}

	// Dijkstra
	costTable := make(map[string]int)
	for name := range reference {
		costTable[string(name)] = infinity
	}
	costTable["@"] = 0

	visited := make(map[string]bool)

	for len(visited) != len(reference) {
		next, _ := lowerCostNode(costTable, visited)

		visited[next] = true

		for _, path := range neighbours[next[0]] {
			distanceNeighbour := path.length + costTable[next]
			if distanceNeighbour < costTable[path.edge.end.name] {
				costTable[path.edge.end.name] = distanceNeighbour
			}
		}
	}

	fmt.Println(costTable)

}

func part2() {

}
