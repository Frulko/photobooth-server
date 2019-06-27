package main

import (
	"flag"
	"image/color"
	"time"

	"github.com/mcuadros/go-rpi-ws281x"
)


func main() {
	// create a new canvas with the given width and height, and the config, in this
	// case the configuration is for a Unicorn pHAT (8x4 pixels matrix) with the
	// default configuration
	c, _ := ws281x.NewCanvas(1, 2, &ws281x.DefaultConfig)


	// initialize the canvas and the matrix
	c.Initialize()

	// since ws281x implements image.Image any function like draw.Draw from the std
	// library may be used with it.
	// 
	// now we copy a white image into the ws281x.Canvas, this turn on all the leds
	// to white
	draw.Draw(c, c.Bounds(), image.NewUniform(color.White), image.ZP, draw.Over)

	// render and sleep to see the leds on
	c.Render()
	time.Sleep(time.Second * 5)

	// don't forget close the canvas, if not you leds may remain on
	c.Close()
}
