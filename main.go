package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gookit/color"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

type Cell struct {
	color      int
	x          int
	y          int
	discovered bool
}

var clear map[string]func() //create a map for storing clear funcs

func main() {
	initConsoleWipe()

	// Get width and height from the command line, or default to 5x10.
	var height int
	var width int
	flag.IntVar(&height, "h", 5, "specify the height of the flood, defaults to 5")
	flag.IntVar(&width, "w", 10, "specify the width of the flood, defaults to 10")
	flag.Parse()

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
			fmt.Printf("Complete!, you did it in %d tries\n", tries)
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

//Wipes the console on linux or windows.
//(https://stackoverflow.com/questions/22891644/how-can-i-clear-the-terminal-screen-in-go)
func wipeConsole() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func initConsoleWipe() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
