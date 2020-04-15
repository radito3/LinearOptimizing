package main

import "fmt"

type Cell struct {
	value int
	row   int
	col   int

	isBaseZero bool

	visited bool
}

type Matrix []Cell

func (m Matrix) Len() int {
	return len(m)
}

func (m Matrix) Less(i, j int) bool {
	return m[i].value < m[j].value
}

func (m Matrix) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func newMatrix(val *[height][width]int) Matrix {
	var matrix Matrix
	for row, i := range val {
		for col, j := range i {
			var cell Cell
			if j == -1 {
				cell.value = 0
				cell.isBaseZero = true
				val[row][col] = 0
			} else {
				cell.value = j
			}
			cell.row = row
			cell.col = col
			matrix = append(matrix, cell)
		}
	}
	return matrix
}

const height = 3
const width = 4

func main() {
	A := [3]int{20, 15, 25}
	B := [4]int{13, 17, 19, 11}
	C := [3][4]int{{6, 5, 2, 1}, {3, 5, 4, 2}, {5, 3, 6, 3}}

	fmt.Println("North east corner + Distribution methods:")

	X := northEastCornerMethod(A, B)
	fmt.Printf("%v\n", X)

	delta, sum := distributionMethod(&C, &X)
	fmt.Printf("%v\n", delta)
	fmt.Printf("%v\n", X)
	fmt.Printf("%d\n", sum)

	fmt.Println("North east corner + Method of potentials:")

	X = northEastCornerMethod(A, B)
	fmt.Printf("%v\n", X)

	delta, sum = methodOfPotentials(&C, &X)
	fmt.Printf("%v\n", delta)
	fmt.Printf("%v\n", X)
	fmt.Printf("%d\n", sum)

	fmt.Println("Minimal element + Distribution methods:")

	X = minElementMethod(A, B, C)
	fmt.Printf("%v\n", X)

	delta, sum = distributionMethod(&C, &X)
	fmt.Printf("%v\n", delta)
	fmt.Printf("%v\n", X)
	fmt.Printf("%d\n", sum)

	fmt.Println("Minimal element + Method of potentials:")

	X = minElementMethod(A, B, C)
	fmt.Printf("%v\n", X)

	delta, sum = methodOfPotentials(&C, &X)
	fmt.Printf("%v\n", delta)
	fmt.Printf("%v\n", X)
	fmt.Printf("%d\n", sum)
}
