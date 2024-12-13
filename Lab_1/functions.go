// Name: Nathnael Getachew
// ID:   UGR/8932/13
package main
import "fmt"
// Function to add two numbers
func add(a int, b int) int {
 return a + b
}
func main() {
 result := add(5, 3)
 fmt.Println("Sum:", result)
 // Control flow example: if-else statement
 if result > 5 {
 fmt.Println("The result is greater than 5")
 } else {
 fmt.Println("The result is 5 or less")
 }
// Loop example: for loop
 for i := 0; i < 5; i++ {
 fmt.Println("Loop iteration:", i)
 }
}