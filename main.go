/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 17-01-2018
 * |
 * | File Name:     main.go
 * +===============================================
 */

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/LightsPlatform/vSensor/sensor"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var sensors map[string]*sensor.Sensor

// init initiates global variables
func init() {
	sensors = make(map[string]*sensor.Sensor)
}

// handle registers apis and create http handler
func handle() http.Handler {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/about", aboutHandler)

		api.POST("/sensor/:id", sensorHandler)
		api.GET("/sensor/:id/data", dataHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "404 Not Found"})
	})

	return r
}

func main() {
	fmt.Println("vSensor Light @ 2018")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handle(),
	}

	go func() {
		fmt.Printf("vSensor Listen: %s\n", srv.Addr)
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("Listen Error:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("vSensor Shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Shutdown Error:", err)
	}
}

func aboutHandler(c *gin.Context) {
	c.String(http.StatusOK, "18.20 is leaving us")
}

func sensorHandler(c *gin.Context) {
	id := c.Param("id")
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sensor, err := sensor.New(id, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if _, ok := sensors[id]; !ok {
		go sensor.Run()
	}
	sensors[id] = sensor

	c.String(http.StatusOK, id)
}

func dataHandler(c *gin.Context) {
	id := c.Param("id")
	data := make([]sensor.Data, 0)

	sensor, ok := sensors[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Sensor %s was not found on vSensor", id)})
		return
	}

	exists := true
	for exists {
		select {
		case d := <-sensor.Buffer:
			data = append(data, d)
		default:
			exists = false
		}
	}

	c.JSON(http.StatusOK, data)
}
