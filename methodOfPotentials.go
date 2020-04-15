package main

import (
	"fmt"
	"os"
)

func findCellDelta(cell Cell, U [height]int, V [width]int, C [height][width]int) int {
	return U[cell.row] + V[cell.col] - C[cell.row][cell.col]
}

func computePotentials(C, X [height][width]int) (U [height]int, V [width]int) {
	U[0] = 0
	//TODO implement
	return
}

//chose either a 'u' (U's are mapped to rows) or a 'v' (V's are mapped to columns) and assign it a value
//    preferably 0 (e.g. u1 = 0)
//for a full cell on the chosen row/column calculate the remaining potential with this
//    uX = C(X,i) - vi where i is the column of a full cell in row X
//    vX = C(i,X) - ui -> X(i,X) != 0
//calculate all U's and V's
//calculate empty cells' D's with this
//    d(i,j) = ui + vj - C(i,j)
//if any d > 0 -> for the biggest d do:
//         find empty cell loop
//         mark cells values (X's) by starting with +, then -, and so on,
//         find the min value (X[i][j]) from the ones with -
//         write that value in the beginning of the cycle
//         recalculate the rest of the cells' values by either adding or subtracting that value from their value
//         after recalculation, if there are cells whose value is 0 (the number of empty cells in the matrix
//             is no longer height + width - 1), one of them must become a "base 0" and be treated as a full cell
//         then start from beginning of algorithm
//calculate the sum of the values of the full cells multiplied by it's C (X[i][j] * C[i][j])

func methodOfPotentials(C, X *[height][width]int) ([]int, int) {
	matrix := newMatrix(X)
	emptyCells := getEmptyCells(matrix)
	invalidDeltas := make(map[int][]Cell)
	var validDeltas []int

	U, V := computePotentials(*C, *X)

	for _, cell := range emptyCells {
		delta := findCellDelta(cell, U, V, *C)
		if delta <= 0 {
			validDeltas = append(validDeltas, delta)
			continue
		}

		loop := []Cell{cell}
		hasLoop := findLoop(cell, cell, matrix, &loop)
		if !hasLoop {
			//retry with different loops? and/or different base zero
			fmt.Fprintf(os.Stderr, "cycle for cell(%d,%d) not found\n", cell.row+1, cell.col+1)
		}

		clearVisitedCells(&matrix)

		invalidDeltas[delta] = loop
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
