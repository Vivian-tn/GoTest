package pkg

import "fmt"

func forXunHuan() {
	tests := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for i, test := range tests {
		if i < 3 {
			fmt.Println("default")
		} else {
			fmt.Println(test)
		}
	}
}
