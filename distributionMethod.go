package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"time"
)

type Direction uint

const (
	DirectionUP Direction = iota
	DirectionDOWN
	DirectionLEFT
	DirectionRIGHT
)

func getAvailableDirections(emptyCell Cell, matrix Matrix) []Direction {
	var directions []Direction
	var adjacentCell *Cell

	if emptyCell.col != 0 {
		adjacentCell = findNextCorner(emptyCell, DirectionLEFT, matrix)
		if adjacentCell != nil && !adjacentCell.visited {
			directions = append(directions, DirectionLEFT)
		}
	}

	if emptyCell.col != width-1 { //max width
		adjacentCell = findNextCorner(emptyCell, DirectionRIGHT, matrix)
		if adjacentCell != nil && !adjacentCell.visited {
			directions = append(directions, DirectionRIGHT)
		}
	}

	if emptyCell.row != 0 {
		adjacentCell = findNextCorner(emptyCell, DirectionUP, matrix)
		if adjacentCell != nil && !adjacentCell.visited {
			directions = append(directions, DirectionUP)
		}
	}

	if emptyCell.row != height-1 { //max height
		adjacentCell = findNextCorner(emptyCell, DirectionDOWN, matrix)
		if adjacentCell != nil && !adjacentCell.visited {
			directions = append(directions, DirectionDOWN)
		}
	}

	return directions
}

func getEmptyCells(matrix Matrix) []Cell {
	var emptyCells []Cell
	for _, cell := range matrix {
		if cell.value == 0 && !cell.isBaseZero {
			emptyCells = append(emptyCells, cell)
		}
	}
	return emptyCells
}

func findLoop(base, cell Cell, matrix Matrix, loop *[]Cell) bool {
	if len(*loop) >= 4 && len(*loop) % 2 == 0 {
		endCell := (*loop)[len(*loop)-1]
		if endCell.col == base.col || endCell.row == base.row {
			return true
		}
	}

	directions := getAvailableDirections(cell, matrix)
	neighbours := neighbouringCells{base: base}

	for _, dir := range directions {
		neighbour := findNextCorner(cell, dir, matrix)
		if neighbour != nil {
			neighbours.cells = append(neighbours.cells, *neighbour)
		}
	}

	if len(neighbours.cells) == 0 {
		return false
	}

	sort.Stable(neighbours) //preserves the order of equal adjacent elements

	for _, neighbour := range neighbours.cells {
		matrix[neighbour.row*width+neighbour.col].visited = true
		*loop = append(*loop, neighbour)

		if findLoop(base, neighbour, matrix, loop) {
			return true
		}

		matrix[neighbour.row*width+neighbour.col].visited = false
		*loop = (*loop)[:len(*loop)-1]
	}

	return false
}

type neighbouringCells struct {
	cells []Cell
	base  Cell
}

func (m neighbouringCells) Len() int {
	return len(m.cells)
}

func (m neighbouringCells) Less(i, j int) bool {
	return distanceBetweenCells(m.cells[i], m.base) < distanceBetweenCells(m.cells[j], m.base)
}

func (m neighbouringCells) Swap(i, j int) {
	m.cells[i], m.cells[j] = m.cells[j], m.cells[i]
}

func distanceBetweenCells(x, y Cell) float64 {
	return math.Sqrt(math.Pow(float64(x.row-y.row), 2) + math.Pow(float64(x.col-y.col), 2))
}

func findCell(row, col int, matrix Matrix) Cell {
	return matrix[row*width+col]
}

func findNextCorner(cell Cell, direction Direction, matrix Matrix) *Cell {
	switch direction {
	case DirectionDOWN:
		if cell.row == 2 {
			return nil
		} else {
			neighbour := findCell(cell.row+1, cell.col, matrix)
			if neighbour.value != 0 || neighbour.isBaseZero {
				return &neighbour
			}
			return findNextCorner(neighbour, direction, matrix)
		}
	case DirectionUP:
		if cell.row == 0 {
			return nil
		} else {
			neighbour := findCell(cell.row-1, cell.col, matrix)
			if neighbour.value != 0 || neighbour.isBaseZero {
				return &neighbour
			}
			return findNextCorner(neighbour, direction, matrix)
		}
	case DirectionRIGHT:
		if cell.col == 3 {
			return nil
		} else {
			neighbour := findCell(cell.row, cell.col+1, matrix)
			if neighbour.value != 0 || neighbour.isBaseZero {
				return &neighbour
			}
			return findNextCorner(neighbour, direction, matrix)
		}
	case DirectionLEFT:
		if cell.col == 0 {
			return nil
		} else {
			neighbour := findCell(cell.row, cell.col-1, matrix)
			if neighbour.value != 0 || neighbour.isBaseZero {
				return &neighbour
			}
			return findNextCorner(neighbour, direction, matrix)
		}
	}
	return nil
}

func findLoopDelta(loop []Cell, C [height][width]int) int {
	delta := 0
	isMinus := true
	for _, cell := range loop {
		if isMinus {
			delta -= C[cell.row][cell.col]
			isMinus = false
		} else {
			delta += C[cell.row][cell.col]
			isMinus = true
		}
	}
	return delta
}

func maxKey(deltas map[int][]Cell) int {
	maxNumber := math.MinInt32
	for n := range deltas {
		if n > maxNumber {
			maxNumber = n
		}
	}
	return maxNumber
}

//FIXME maybe incorrect?
func recalculateCells(invalidDeltas map[int][]Cell, X *[height][width]int) {
	maxDeltaLoop := invalidDeltas[maxKey(invalidDeltas)]
	var arr []Cell
	for i := 1; i < len(maxDeltaLoop); i += 2 {
		arr = append(arr, maxDeltaLoop[i])
	}
	temp := Matrix(arr)
	sort.Stable(temp)
	minElem := temp[0]
	row, col := maxDeltaLoop[0].row, maxDeltaLoop[0].col
	X[row][col] = minElem.value
	numEmptyCells := 0

	isMinus := true
	for _, cell := range maxDeltaLoop[1:] {
		if isMinus {
			X[cell.row][cell.col] -= minElem.value
			isMinus = false
			if X[cell.row][cell.col] == 0 {
				numEmptyCells++
			}
		} else {
			X[cell.row][cell.col] += minElem.value
			isMinus = true
		}
	}

	//DEBUG
	fmt.Printf("%v\n", X)
	time.Sleep(1*time.Second)

	if numEmptyCells > 1 {
		for _, cell := range maxDeltaLoop[1:] {
			if X[cell.row][cell.col] == 0 {
				X[cell.row][cell.col] = -1 //become base zero
				break
			}
		}
	}
}

//find cycles for every empty cell (X[i][j] == 0) find the loop starting and ending at the empty cell that
//             only turns at full cells (goes through empty cells) and is shortest (closest to beginning cell))
//then for each cycle do:
//mark C's with -, then +, and so on, beginning from empty cell, and calculate sum of C's with corresponding signs
//    if > 0 -> add to list
//           find max element in list
//           for the cycle that begins with that element:
//              mark cells values (X's) starting with +, then -, and so on
//              find the min value (X[i][j]) from the ones with -
//              write that value in the beginning of the cycle
//              recalculate the rest of the cells' values by either adding or subtracting that value from their value
//           after recalculation, if there are cells whose value is 0 (the number of empty cells in the matrix
//           	is no longer height + width - 1), one of them must become a "base 0" and be treated as a full cell
//           then start from beginning of algorithm
//calculate the sum of the values of the full cells multiplied by it's C (X[i][j] * C[i][j])

func distributionMethod(C, X *[height][width]int) ([]int, int) {
	matrix := newMatrix(X)
	emptyCells := getEmptyCells(matrix)
	invalidDeltas := make(map[int][]Cell)
	var validDeltas []int

	for _, cell := range emptyCells {
		loop := []Cell{cell}
		hasLoop := findLoop(cell, cell, matrix, &loop)
		if !hasLoop {
			//retry with different loops? and/or different base zero
			fmt.Fprintf(os.Stderr, "\ncycle for cell(%d,%d) not found\n", cell.row+1, cell.col+1)
		}

		clearVisitedCells(&matrix)

		delta := findLoopDelta(loop, *C)
		if delta > 0 {
			invalidDeltas[delta] = loop
		} else {
			validDeltas = append(validDeltas, delta)
		}
	}

	if len(invalidDeltas) != 0 {
		recalculateCells(invalidDeltas, X)
		return distributionMethod(C, X)
	}

	sum := 0
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if X[i][j] != 0 {
				sum += X[i][j] * C[i][j]
			}
		}
	}
	return validDeltas, sum
}

func clearVisitedCells(matrix *Matrix) {
	for _, cell := range *matrix {
		if cell.visited {
			(*matrix)[cell.row*width+cell.col].visited = false
		}
	}
}
