package main

// import "golang.org/x/tour/pic"

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/fogleman/gg"
)

func Fibonacci(x uint16) (fibs []uint16) {
	// returns x fibonacci numbers
	if x < 3 {
		return fibs
	}

	// create slice if x is valid
	fibs = make([]uint16, x)
	fibs[0] = 1
	fibs[1] = 1

	for i := 2; i < cap(fibs); i++ {
		fibs[i] = fibs[i-1] + fibs[i-2]
	}
	//fmt.Println(fibs)
	return fibs
}

func DrawArc(centerX, centerY, radius uint16, startAngle, endAngle float64, imgArr [][]uint8, color uint8) {
	height := len(imgArr)
	width := len(imgArr[0])

	// Ensure angles are in the range [0, 2*pi)
	if startAngle > endAngle {
		endAngle += 2 * math.Pi
	}

	angle_resolution := 0.3 / float64(radius)

	// Iterate over possible points in the circle
	for theta := startAngle; theta <= endAngle; theta += angle_resolution { // Increment angle for finer detail
		x := int(float64(centerX) + float64(radius)*math.Cos(theta))
		y := int(float64(centerY) + float64(radius)*math.Sin(theta))

		if x >= 0 && x < width && y >= 0 && y < height {
			imgArr[y][x] = color
		}
	}
}

func CreateSpiralArr(img_height, img_width, fibonacci_depth uint16) [][]uint8 {
	// returns a slice of length dy, each element of which is a slice of
	// dx uint8

	// DrawArc := func(x, y, r int, arc_length float64) {
	// 	fmt.Printf("%v %v %v %v\n", x, y, r, arc_length)
	// }

	// const angle_resolution = 0.1 // rads
	// starting point
	center_x := img_width / 2
	center_y := img_height / 2

	fibs := Fibonacci(fibonacci_depth)

	var s [][]uint8 = make([][]uint8, img_height)
	for i := 0; i < int(img_height); i++ {
		s[i] = make([]uint8, img_width)
	}
	cur_angle := math.Pi
	DrawArc(center_x, center_y, fibs[0], math.Pi/2, math.Pi, s, 255)
	for cur_fib := 1; cur_fib < len(fibs); cur_fib++ {
		if fibs[cur_fib] > img_height {
			break
		}
		start := time.Now()
		DrawArc(center_x, center_y, fibs[cur_fib], cur_angle, cur_angle+math.Pi*2, s, 255)
		elapsed := time.Since(start)
		fmt.Printf("Circle: rad = %v time = %v\n", fibs[cur_fib], elapsed)
		cur_angle += math.Pi / 2
		// fmt.Println(int(math.Cos((math.Pi * float64(cur_fib)) / 4)))
		// Move the center point of the next circle/arc
		center_x += uint16((float64(fibs[cur_fib-1])) * math.Cos((math.Pi*float64(cur_fib))/2))
		center_y += uint16((float64(fibs[cur_fib-1])) * math.Sin((math.Pi*float64(cur_fib))/2))
		//fmt.Println(math.Cos((math.Pi * float64(cur_fib)) / 2))
	}

	return s
}

func CreateImg(img_arr [][]uint8, width int, height int, filename string) string {
	context := gg.NewContext(width, height)

	// Iterate over the grayscale data and set pixel colors
	for y, row := range img_arr {
		for x, value := range row {
			// Create a grayscale color
			gray := color.Gray{Y: value}
			context.SetColor(gray)
			context.SetPixel(x, y)
		}
	}

	// Save the image to a file
	if err := context.SavePNG(filename); err != nil {
		log.Fatal(err)
	}

	log.Printf("Image saved as %v", filename)
	return filename
}

func main() {
	fibonacci_depth := 16
	img_height := 1024
	img_width := 1024
	img_arr := CreateSpiralArr(uint16(img_height), uint16(img_width), uint16(fibonacci_depth))
	output_img_name := CreateImg(img_arr, len(img_arr[0]), len(img_arr), "output1.png")
	fmt.Println(output_img_name)
}
