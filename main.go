package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
)

func AddPhtoFrame() {
	file, err := os.Open("./res/gophergala.jpg")
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}
	m := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dx()))
	yellow := color.RGBA{255, 255, 0, 255}
	red := color.RGBA{255, 0, 0, 255}
	green := color.RGBA{0, 255, 0, 255}
	blue := color.RGBA{0, 0, 255, 255}

	dx := img.Bounds().Dx()
	dy := img.Bounds().Dy()

	y1 := dy / 8

	y2 := y1 + 20

	y3 := dy - y1 - y1

	y4 := y3 - 20

	x1 := dx / 8

	x2 := x1 + 20

	x3 := dx - x1

	x4 := x3 - 20
	for x := 0; x <= dx; x++ {
		for y := 0; y <= dy; y++ {

			if x < x2 {
				if x < x1 {
					m.Set(x, y, blue)
				} else {
					m.Set(x, y, green)
				}
				continue
			}

			if x > x4 {
				if x > x3 {
					m.Set(x, y, green)
				} else {
					m.Set(x, y, blue)
				}
				continue
			}

			if y < y2 {
				if y < y1 {
					m.Set(x, y, yellow)
				} else {
					m.Set(x, y, red)
				}
				continue
			}

			if y > y4 {
				if y > y3 {
					m.Set(x, y, red)
				} else {
					m.Set(x, y, yellow)
				}
				continue
			}
			m.Set(x, y, img.At(x, y))
		}
	}

	if file, err := os.Create("./res/gophergala_m.jpg"); err != nil {
		panic(err)
	} else {
		jpeg.Encode(file, m, nil)
		defer file.Close()
	}
	//x779
	//y888

}

func main() {
	AddPhtoFrame()
	println(`This final exercise,let's add a photo frame for gala!`)
}
