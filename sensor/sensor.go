/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 02-02-2018
 * |
 * | File Name:     sensor/sensor.go
 * +===============================================
 */

package sensor

import "time"

// Sensor represents virtual sensor that
// only generate random data with given generator
type Sensor struct {
	id     int
	Name   string
	Buffer chan Data
	// TODO Generator
}

// Data represents sensor data that contains
// time and value
type Data struct {
	Time  time.Time
	Value string
}

// New creates new sensor and store its user given script
func New(name string, script []byte) *Sensor {
	return &Sensor{}
}

// Run runs sensor, running sensor generate data using
// user given script.
// it is a blocking function so run it in new thread
func (s *Sensor) Run() {
}
