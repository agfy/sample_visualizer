package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/jpeg"
	"log"
	"os"
	"strconv"
)

type MyRGB struct {
	R uint32
	G uint32
	B uint32
}

func (m MyRGB) toString() string {
	return strconv.FormatUint(uint64(m.R), 10) + " " + strconv.FormatUint(uint64(m.G), 10) + " " +
		strconv.FormatUint(uint64(m.B), 10)
}

func ExistInColors(colors []MyRGB, color MyRGB) bool {
	for _, c := range colors {
		if c == color {
			return true
		}
	}
	return false
}

var (
	greens = []MyRGB{{26201, 52491, 26211}, {26108, 59062, 26437}, {26560, 52308, 26211},
		{26108, 59150, 25983}, {25749, 59156, 26890}, {27185, 58689, 25529},
		{26467, 58879, 26437}, {25842, 52497, 27118}, {26467, 58791, 26890},
		{25749, 59333, 25983}, {26826, 58696, 26437}, {26099, 53019, 26014},
		{25842, 52762, 25757}, {26919, 51861, 27572}, {26560, 52132, 27118},
		{24765, 53134, 26664}, {26201, 52579, 25757}, {27544, 58154, 27344}}
	greys = []MyRGB{{26214, 26214, 26214}, {26985, 26985, 26985}, {27499, 27499, 27499},
		{27446, 26538, 25821}, {25855, 26485, 25760}, {26471, 26471, 26471},
		{27397, 27763, 28663}, {26985, 27073, 26531}, {27344, 26714, 27439},
		{26728, 26728, 26728}}
	blacks = []MyRGB{{0, 0, 0}, {0, 88, 0}, {0, 270, 0}, {257, 257, 257},
		{257, 446, 1164}, {0, 446, 1164}}
	reds  = []MyRGB{{204 << 8, 0, 0}, {255 << 8, 0, 0}, {52184, 112, 2325}, {51311, 309, 0}}
	blues = []MyRGB{{735, 0, 64807}, {13461, 0, 65290}, {13102, 0, 64836},
		{13820, 0, 65535}, {26560, 52220, 26664}, {376, 0, 64807}, {13461, 0, 63929},
		{17, 0, 64807}, {0, 1340, 0}, {1094, 0, 62085}}
	whites = []MyRGB{{61229, 65535, 65535}, {64459, 65535, 64175}}
)

func main() {
	samples := 8
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	picturePath := pwd + "/output.jpeg"
	f, err := os.Open(picturePath)
	if err != nil {
		fmt.Println(err)
	}
	fi, _ := f.Stat()
	fmt.Println(fi.Name())
	//defer f.Close()sss
	img, format, err := image.Decode(f)
	if err != nil {
		fmt.Println("Decoding error:", err.Error(), img.ColorModel(), format)
	}

	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y

	mostFrequentColors := make([]MyRGB, 0)
	for y := 0; y < height; y += 4 * samples {
		for x := 0; x < width; x += 6 * samples {
			rgbMap := make(map[color.Color]int)
			for localY := 0; localY < 4*samples; localY++ {
				for localX := 0; localX < 6*samples; localX++ {
					if x+localX >= width || y+localY >= height {
						continue
					}
					rgbColor := img.At(x+localX, y+localY)
					rgbMap[rgbColor]++
				}
			}

			mostFrequentColor := img.At(0, 0)
			maxValue := 0
			for k, v := range rgbMap {
				if v > maxValue {
					maxValue = v
					mostFrequentColor = k
				}
			}

			r, g, b, _ := mostFrequentColor.RGBA()
			mostFrequentColors = append(mostFrequentColors, MyRGB{R: r, G: g, B: b})
		}
	}

	bytes := make([]byte, 0, width*height*3)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			localX := x / (6 * samples)
			localY := y / (4 * samples)
			var localWidth int
			if width%(6*samples) == 0 {
				localWidth = width / (6 * samples)
			} else {
				localWidth = width/(6*samples) + 1
			}

			mostFrequentColor := mostFrequentColors[localX+localY*localWidth]
			bytes = append(bytes, byte(mostFrequentColor.R>>8), byte(mostFrequentColor.G>>8), byte(mostFrequentColor.B>>8))
		}
	}

	unsupportedCount := 0
	white_bytes := make([]byte, 0, width*height*3)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			localX := x / (6 * samples)
			localY := y / (4 * samples)
			var localWidth int
			if width%(6*samples) == 0 {
				localWidth = width / (6 * samples)
			} else {
				localWidth = width/(6*samples) + 1
			}

			mostFrequentColor := mostFrequentColors[localX+localY*localWidth]
			if ExistInColors(greens, mostFrequentColor) || ExistInColors(blacks, mostFrequentColor) ||
				ExistInColors(reds, mostFrequentColor) || ExistInColors(whites, mostFrequentColor) ||
				ExistInColors(blues, mostFrequentColor) {
				white_bytes = append(white_bytes, byte(0), byte(0), byte(0))
			} else if ExistInColors(greys, mostFrequentColor) {
				white_bytes = append(white_bytes, byte(255), byte(255), byte(255))
			} else {
				println("unsupported color", mostFrequentColor.toString(), strconv.Itoa(unsupportedCount))
				unsupportedCount++
			}
		}
	}

	//creates a picture with mesh
	bytes = make([]byte, 0, width*height*3)
	for y := img.Bounds().Min.Y; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if (x%(6*samples) == 0 || y%(samples*4) == 0) && (y < 350) {
				r, g, b = 0, 0, 0
			}
			bytes = append(bytes, byte(b>>8), byte(g>>8), byte(r>>8))
		}
	}

	mat, err := gocv.NewMatFromBytes(height, width, gocv.MatTypeCV8UC3, bytes)
	if err != nil {
		log.Fatal(err)
	}

	img2, err := mat.ToImage()
	if err != nil {
		log.Fatal(err)
	}
	w, _ := os.Create("2.jpg")
	defer w.Close()
	jpeg.Encode(w, img2, &jpeg.Options{Quality: 90})

	white_mat, err := gocv.NewMatFromBytes(height, width, gocv.MatTypeCV8UC3, white_bytes)
	if err != nil {
		log.Fatal(err)
	}

	img3, err := white_mat.ToImage()
	if err != nil {
		log.Fatal(err)
	}
	w, _ = os.Create("3.jpg")
	defer w.Close()
	jpeg.Encode(w, img3, &jpeg.Options{Quality: 90})
}
