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

import "time"

// Generator provides a general way for generating load
// Generate runs new go routine that passes number
// of packets that must be generated into channel.
type Generator interface {
	Generate() <-chan int
}

// UniformGenerator generates traffic uniformaly based on time
type UniformGenerator struct {
	Timeslot time.Duration
}

// Generate runs new go routine that passes number
// of packets that must be generated into channel.
func (u UniformGenerator) Generate() <-chan int {
	g := make(chan int, 0)
	t := time.Tick(u.Timeslot)

	go func() {
		for {
			select {
			case <-t:
				g <- 1
			}
		}
	}()

	return g
}
