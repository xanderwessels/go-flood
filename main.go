package main

import (
	"bufio"
	"fmt"
	"github.com/gookit/color"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type Cell struct {
	color      int
	x          int
	y          int
	discovered bool
}

func main() {
	// Get width and height from the command line, or default to 5x10.
	height, err := strconv.Atoi(os.Args[1])
	if err != nil {
		height = 5
		log.Fatal(err)
	}
	width, err := strconv.Atoi(os.Args[2])
	if err != nil {
		width = 10
		log.Fatal(err)
	}

	//Setup a grid and the (empty) current selection.
	fullGrid := make([][]Cell, height)
	currentSelection := make([]Cell, 0, width*height)
	for i := range fullGrid {
		fullGrid[i] = make([]Cell, width)
	}

	//Fill the grid with random colors.
	rand.Seed(time.Now().Unix())
	for i := range fullGrid {
		for j := range fullGrid[i] {
			fullGrid[i][j] = Cell{rand.Intn(5), i, j, false}
		}
	}

	//Add (0,0) to the current selection.
	fullGrid[0][0].discovered = true
	currentSelection = append(currentSelection, fullGrid[0][0])

	reader := bufio.NewReader(os.Stdin)

	tries := 0
	for {
		wipeConsole()
		printGrid(fullGrid)

		//Check if we have filled the full grid.
		if len(currentSelection) == width*height {
			fmt.Printf("Complete!, you did it in %d tries", tries)
			break
		}

		//Get user input.
		fmt.Print("Select color: ")
		text, _ := reader.ReadString('\n')
		sel, _ := strconv.Atoi(string(text[0]))

		tries++

		//Change the color of the current selection and the grid.
		for index, cs := range currentSelection {
			currentSelection[index].color = sel
			fullGrid[cs.x][cs.y].color = sel
		}

		//Find the new selection.
		currentSelection = append(currentSelection, findNewSelection(currentSelection, fullGrid)...)
	}
}

func findNewSelection(currentSelection []Cell, fullGrid [][]Cell) []Cell {
	//Width and height parameters are needed for bounds checking.
	height := len(fullGrid)
	width := len(fullGrid[0])

	//Find surrounding cells with same color that haven't been discovered yet.
	newSelection := make([]Cell, 0, width*height)
	for _, cs := range currentSelection {
		//Check the top of the cell
		if cs.x > 0 {
			topN := &fullGrid[cs.x-1][cs.y]
			if cs.color == topN.color && !topN.discovered {
				topN.discovered = true
				newSelection = append(newSelection, *topN)
			}
		}
		//Check to down of the cell
		if cs.x < height-1 {
			downN := &fullGrid[cs.x+1][cs.y]
			if cs.color == downN.color && !downN.discovered {
				downN.discovered = true
				newSelection = append(newSelection, *downN)
			}
		}
		//Check the left of the cell
		if cs.y > 0 {
			leftN := &fullGrid[cs.x][cs.y-1]
			if cs.color == leftN.color && !leftN.discovered {
				leftN.discovered = true
				newSelection = append(newSelection, *leftN)
			}
		}
		//Check the right of the cell
		if cs.y < width-1 {
			rightN := &fullGrid[cs.x][cs.y+1]
			if cs.color == rightN.color && !rightN.discovered {
				rightN.discovered = true
				newSelection = append(newSelection, *rightN)
			}
		}
	}

	//If we haven't found any new cells, we return.
	if len(newSelection) == 0 {
		return nil
	}

	//If we have found new surrounding cells with the same color, we recurse to find any cells that surround that cell
	//with the same color.
	return append(newSelection, findNewSelection(newSelection, fullGrid)...)
}

func printGrid(grid [][]Cell) {

	for i := range grid {
		for j := range grid[i] {
			c := color.BgBlue
			switch grid[i][j].color {
			case 0:
				c = color.BgGreen
			case 1:
				c = color.BgRed
			case 2:
				c = color.BgBlue
			case 3:
				c = color.BgYellow
			case 4:
				c = color.BgGray
			}
			c.Printf(" ")
			c.Printf(strconv.Itoa(grid[i][j].color))
			c.Printf(" ")
		}
		fmt.Println()
	}
}

func wipeConsole() {
	cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}