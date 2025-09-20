package main

import "fmt"

func main() {
	var n int
	fmt.Scanln(&n)

	nums := make([]int, n)

	for i := 0; i < n; i++ {
		fmt.Scan(&nums[i])
	}

	maxVal := 0
	for i := 0; i < n; i++ {
		if nums[i] > maxVal {
			maxVal = nums[i]
		}
	}

	fmt.Println(maxVal)
}
