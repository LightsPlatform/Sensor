/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 02-02-2018
 * |
 * | File Name:     generators/generator.go
 * +===============================================
 */

package generators

// Generator provides a general way for generating load
// Generate calls in each iteration and returns number
// of packets that must be generated.
type Generator interface {
	Generate() <-chan int
}
