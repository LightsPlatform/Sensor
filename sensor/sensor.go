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

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	log "github.com/sirupsen/logrus"
)

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
func New(name string, script []byte) (*Sensor, error) {
	// Store user script
	f, err := os.Create(fmt.Sprintf("/tmp/sensor-%s.py", name))
	if err != nil {
		return nil, err
	}
	f.Write(script)

	return &Sensor{
		Name:   name,
		Buffer: make(chan Data, 1024),
	}, nil
}

// Run runs sensor, running sensor generate data using
// user given script.
// it is a blocking function so run it in new thread
func (s *Sensor) Run() {
	t := time.Tick(1 * time.Second)
	for {
		select {
		case <-t:
			cmd := exec.Command("runtime.py", "--job", "rule", fmt.Sprintf("/tmp/sensor-%s.py", s.Name))

			// run
			if _, err := cmd.Output(); err != nil {
				if err, ok := err.(*exec.ExitError); ok {
					log.Errorf("%s: %s", err.Error(), err.Stderr)
				}
			}
		}
	}
}
