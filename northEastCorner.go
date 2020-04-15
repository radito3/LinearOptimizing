package main

func northEastCornerMethod(A [height]int, B [width]int) [height][width]int {
	X := [height][width]int{} //the default value of int in Go is 0
	row, col := 0, 0

	for steps := 0; steps < height+width-1; steps++ {
		if A[row] < B[col] {
			X[row][col] = A[row]
			B[col] -= A[row]
			row++
		} else if A[row] > B[col] {
			X[row][col] = B[col]
			A[row] -= B[col]
			col++
		} else {
			X[row][col] = A[row]
			row++
			col++
		}
	}

	return X
}
