/*

	A devdraw listener. Maybe I've named this wrong.

*/

package main

import (
	"fmt"
	"image"
	"log"
	"math/rand"
	"time"

	"9fans.net/go/draw"
)

func watcher() {
	log.Println("invoked watcher function")
}


// redraw repaints the world. This is the view function.
func redraw(d *draw.Display, resized bool) {
	if resized {
		if err := d.Attach(draw.Refmesg); err != nil {
			log.Fatalf("can't reattach to window: %v", err)
		}
	}

	// draw coloured rects at mouse positions
	// first param is the clip rectangle. which can be 0. meaning no clip?
	var clipr image.Rectangle
	fmt.Printf("empty clip? %v\n", clipr)
	d.ScreenImage.Draw(clipr, d.White, nil, image.ZP)
	d.Flush()
}


func main() {
	log.Println("hello from devdraw")

	// Make the window.
	d, err := draw.Init(nil, "", "experiment1", "")
	if err != nil {
		log.Fatal(err)
	}

	// make some colors
	back, _ := d.AllocImage(image.Rect(0, 0, 1, 1), d.ScreenImage.Pix, true, 0xDADBDAff)
	red, _ := d.AllocImage(image.Rect(0, 0, 1, 1), d.ScreenImage.Pix, true, 0xFF0000ff)

	fmt.Printf("background colour: %v\n ", back)

	// get mouse positions
	mousectl := d.InitMouse()
	redraw(d, false)

	go func() {
		rect := d.Image.R
		log.Println(rect)
		// make a random variable...
		start := time.Now()
		outerreps := 2000
		innerreps := 10
		for j := 0; j < outerreps; j++ {
			for i := 0 ; i < innerreps; i++ {
				X := rand.Intn(rect.Max.X)
				Y := rand.Intn(rect.Max.Y)
				d.ScreenImage.Draw(image.Rect(X, Y, X+22, Y+26), back, nil, image.ZP)
			}
			d.Flush()
		}
		elapsed := time.Since(start)
		log.Println( "time per quad: ", (elapsed.Nanoseconds() / int64(outerreps * innerreps)) / 1000, "µs")
	} ()

	for {
		select {
		case <-mousectl.Resize:
			redraw(d, true)
		case m := <-mousectl.C:
			// fmt.Printf("mouse field %v buttons %d\n", m, m.Buttons)
	
			if (m.Buttons & 1 == 1) {
				// Draws little rectangles for each recorded mouse position.
				d.ScreenImage.Draw(image.Rect(m.X, m.Y, m.X+22, m.Y+26), red, nil, image.ZP)
				d.Flush()
			}
		}
	}

	fmt.Print("bye\n")
}
