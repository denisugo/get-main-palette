package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
	"regexp"
	"strings"

	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

func main() {
	var pwd = getPwd()
	var filenames = []string{strings.Join([]string{pwd, "res", "test1.jpeg"}, "/"), strings.Join([]string{pwd, "res", "test2.jpeg"}, "/")}

	for _, filename := range filenames {
		var img, err = getImageFromFile(filename)
		if err != nil {
			fmt.Println("Unable to read file: ", filename)
			fmt.Println(err)
			os.Exit(1)
		}

		var colors = make([]color.RGBA, 0) //[]color.RGBA{{0, 0, 0, 255}} // list is filled by selected colors
		var maxX, maxY = img.Bounds().Max.X, img.Bounds().Max.Y
		for x := 0; x < maxX; x++ {
			for y := 0; y < maxY; y++ {
				var c = img.At(x, y)

				var r, g, b, _ = c.RGBA()
				var include = true
				if len(colors) == 0 {
					if math.Abs(float64(255+255+255)-(bigToFloat(r)+bigToFloat(g)+bigToFloat(b))) < float64(deviation) {
						include = false
					}
				} else {
					for _, rgba := range colors {
						if math.Abs((float64(rgba.R)+float64(rgba.G)+float64(rgba.B))-(bigToFloat(r)+bigToFloat(g)+bigToFloat(b))) < float64(deviation) {
							include = false
						}
					}
				}
				if include {
					colors = append(colors, color.RGBA{bigToSmall(r), bigToSmall(g), bigToSmall(b), 255})
				}

			}
		}
		fmt.Println(filename)
		fmt.Println(colors)
		drawSlice(colors, filename)
	}
}

const deviation uint8 = 30

func getPwd() string {
	var dir, err = os.Getwd()
	if err != nil {
		return ""
	}
	return dir

}

func getImageFromFile(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, err := jpeg.Decode(f)
	return image, err
}

func bigToSmall(in uint32) uint8 {
	const max32 float32 = 0xFFFF // according to documentation of color library
	const max8 float32 = 255
	var d = float32(in) / (max32) // between 0 and 1
	return uint8(d * max8)
}
func bigToFloat(in uint32) float64 {
	const max32 float64 = 0xFFFF // according to documentation of color library
	const max8 float64 = 255
	var d = float64(in) / (max32) // between 0 and 1
	return (d * max8)
}

// type SortByColor []color.RGBA

// func (a SortByColor) Len() int      { return len(a) }
// func (a SortByColor) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
// func (a SortByColor) Less(i, j int) bool {
// 	return (a[i].A + a[i].R + a[i].B + a[i].G) < (a[j].A + a[j].R + a[j].B + a[j].G) // TODO: tba
// }

func drawSlice(list []color.RGBA, filename string) {
	var size = len(list)
	var stripe float64 = 40
	c := canvas.New(stripe*float64(size), 200)
	ctx := canvas.NewContext(c)
	ctx.SetFillColor(canvas.White)
	ctx.DrawPath(0, 0, canvas.Rectangle(c.W, c.H))
	// sort.Sort(SortByColor(list))
	for i, v := range list {
		ctx.SetFillColor(v)
		ctx.DrawPath(stripe*float64(i), 0, canvas.Rectangle(stripe, c.H))
	}

	r, _ := regexp.Compile(`[\w|\d]+\.[\w|\d]+`)
	renderers.Write((strings.Join([]string{getPwd(), "out", strings.Join([]string{r.FindString(filename), "png"}, ".")}, "/")), c, canvas.DPMM(5.0))
}
