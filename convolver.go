package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

func newimg(width int, height int) *image.RGBA {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	return img
}
func shift(x int, y int, img image.Image) image.Image {
	nimg := newimg(img.Bounds().Max.X, img.Bounds().Max.Y)
	for i := 0; i < img.Bounds().Max.X; i++ {
		for j := 0; j < img.Bounds().Max.Y; j++ {
			r, _, _, _ := img.At((i+x)%img.Bounds().Max.X, (j+y)%img.Bounds().Max.Y).RGBA()
			r2 := uint8(r / 257)
			v := color.RGBA{r2, r2, r2, 0xff}
			nimg.Set(i, j, v)

		}
	}
	return nimg
}
func mask(m image.Image, s image.Image) image.Image {
	nimg := newimg(m.Bounds().Max.X, m.Bounds().Max.Y)
	for i := 0; i < m.Bounds().Max.X; i++ {
		for j := 0; j < m.Bounds().Max.Y; j++ {
			r, _, _, _ := s.At(i, j).RGBA()
			s, _, _, _ := m.At(i, j).RGBA()
			t := r * s / (16777216)
			nimg.Set(i, j, color.RGBA{uint8(t), uint8(t), uint8(t), 0xff})

		}
	}
	return nimg
}
func albedo(img image.Image) int {
	w := 0
	for i := 0; i < img.Bounds().Max.X; i++ {
		for j := 0; j < img.Bounds().Max.Y; j++ {
			r, _, _, _ := img.At(i, j).RGBA()
			w += int(r)
		}
	}
	return w
}

func main() {
	mn := "square"
	sn := "star"
	mf, _ := os.Open(mn + ".png")
	defer mf.Close()
	sf, _ := os.Open(sn + ".png")
	defer sf.Close()

	m, _, _ := image.Decode(mf)
	s, _, _ := image.Decode(sf)
	maxx := m.Bounds().Max.X
	maxy := m.Bounds().Max.Y
	var albedos [100][100]int
	wmin := 1000000000
	wmax := 0
	for i := 0; i < maxx; i++ {
		for j := 0; j < maxy; j++ {
			w := albedo(mask(m, shift(i, j, s)))
			albedos[i][j] = w

			if w > wmax {
				wmax = w
			}
			if w < wmin {
				wmin = w
			}
		}
	}
	fmt.Println(wmax, wmin)
	output := newimg(maxx, maxy)
	for i := 0; i < maxx; i++ {
		for j := 0; j < maxy; j++ {
			c := albedos[i][j]
			d := uint8(255*c/(wmax-wmin) - 255*wmin/(wmax-wmin))
			v2 := color.RGBA{d, d, d, 0xff}
			output.Set(i, j, v2)
		}
	}
	fmt.Println("bs")
	f, _ := os.Create(mn + "-" + sn + ".png")
	png.Encode(f, output)
}
