package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"

	"gocv.io/x/gocv"
)

const MinimumArea = 40000

func main() {
	if len(os.Args) < 2 {
		fmt.Println("How to run:\n\tmotion-detect [camera ID]")
		return
	}

	// parse args
	deviceID := os.Args[1]

	webcam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		log.Fatal("Error opening video capture device: %v\n", deviceID)
		return
	}
	defer webcam.Close()

	window := gocv.NewWindow("Motion Window")
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	delta := gocv.NewMat()
	defer delta.Close()

	deltaThreshold := gocv.NewMat()
	defer deltaThreshold.Close()

	mog2 := gocv.NewBackgroundSubtractorMOG2()
	defer mog2.Close()

	log.Printf("Start reading device: %v\n", deviceID)
	carcounter := 0
	cooldown := 0
	for {
		if ok := webcam.Read(&img); !ok {
			log.Printf("Device closed: %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}
		//remove background
		mog2.Apply(img, &delta)
		//Threshold object map to reduce noise
		gocv.Threshold(delta, &deltaThreshold, 250, 255, gocv.ThresholdBinary)

		//counter box coordinates + dimensions
		boxwidth := 40
		boxheight := 900
		box_x := 1700
		box_y := 100

		contours := gocv.FindContours(deltaThreshold, gocv.RetrievalExternal, gocv.ChainApproxSimple)
		for i := 0; i < contours.Size(); i++ {
			area := gocv.ContourArea(contours.At(i))

			//discard anything too small to be counted as a car
			if area < MinimumArea {
				continue
			}

			//create bounding rect for car object
			rect := gocv.BoundingRect(contours.At(i))
			// draw bounding rect into display image
			gocv.Rectangle(&img, rect, color.RGBA{255, 255, 255, 0}, 4)
			// is the rightmost x coordinate inside the counter box?
			if math.Abs(float64(rect.Max.X)-float64(box_x)) < float64(boxwidth) {
				// ensure that we are not counting too often, the same car may be detected several times otherwise
				if cooldown == 0 {
					cooldown = 5
					carcounter++
					log.Printf("Car: %d Area: %f Coordinates: %d, %f\n", carcounter, area, rect.Max.X, math.Abs(float64(rect.Max.X)-float64(box_x)))
				}
			}
		}
		if cooldown > 0 {
			cooldown--
			//log.Printf("Cool %d\n", cooldown)
		}
		// draw counter box into every frame
		gocv.Rectangle(&img, image.Rect(box_x, box_y, box_x+boxwidth, box_y+boxheight), color.RGBA{255, 255, 255, 0}, 10)
		// draw counter into every frame
		gocv.PutText(&img, fmt.Sprint(carcounter), image.Pt(40, 140), gocv.FontHersheyPlain, 12, color.RGBA{255, 255, 255, 0}, 2)
		// show input image (with stuff drawn on it)
		window.IMShow(img)

		// ESC terminates
		if window.WaitKey(1) == 27 {
			break
		}
	}
}
