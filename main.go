package main

import "fmt"

type opStack[T any] []T

func main() {
	fmt.Println("Test")

	t := make(opStack[int], 0) // You must initialize data type here
	t = append(t, 0)
	fmt.Println(t[0])
}
