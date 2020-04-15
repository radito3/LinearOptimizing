package main

import "sort"

func removeRow(matrix Matrix, row int) Matrix {
	var result Matrix
	for _, cell := range matrix {
		if cell.row != row {
			result = append(result, cell)
		}
	}
	return result
}

func removeColumn(matrix Matrix, column int) Matrix {
	var result Matrix
	for _, cell := range matrix {
		if cell.col != column {
			result = append(result, cell)
		}
	}
	return result
}

func minElementMethod(A [height]int, B [width]int, C [height][width]int) [height][width]int {
	X := [height][width]int{} //the default value of int in Go is 0

	matrix := newMatrix(&C)
	sort.Sort(matrix)

	for steps := 0; steps < height+width-1; steps++ {
		min := matrix[0]
		i, j := min.row, min.col

		if A[i] < B[j] {
			X[i][j] = A[i]
			B[j] -= A[i]
			matrix = removeRow(matrix, i)
		} else if A[i] > B[j] {
			X[i][j] = B[j]
			A[i] -= B[j]
			matrix = removeColumn(matrix, j)
		} else {
			X[i][j] = A[i]
			matrix = removeRow(matrix, i)
			matrix = removeColumn(matrix, j)
		}
	}

	return X
}
