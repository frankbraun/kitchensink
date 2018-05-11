// 6bitgen generates prints all possible 6-bit ASCII in ascending order.
package main

import (
	"fmt"
)

func main() {
	for i := 64; i <= 95; i++ {
		fmt.Printf("%c", i)
	}
	fmt.Println()
	for i := 32; i <= 63; i++ {
		fmt.Printf("%c", i)
	}
	fmt.Println()
}
