package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/fogleman/gg"
	"github.com/fogleman/ln/ln"
)

var colors = [...]color.Color{
	color.RGBA{1, 205, 254, 255},
	color.RGBA{71, 209, 213, 255},
	color.RGBA{134, 153, 209, 255},
	color.RGBA{135, 149, 232, 255},
	color.RGBA{148, 208, 255, 255},
	color.RGBA{167, 228, 174, 255},
	color.RGBA{170, 231, 232, 255},
	color.RGBA{173, 140, 255, 255},
	color.RGBA{191, 255, 230, 255},
	color.RGBA{193, 209, 253, 255},
	color.RGBA{198, 172, 201, 255},
	color.RGBA{199, 116, 232, 255},
	color.RGBA{217, 224, 252, 255},
	color.RGBA{219, 169, 206, 255},
	color.RGBA{232, 210, 255, 255},
	color.RGBA{236, 180, 191, 255},
	color.RGBA{239, 230, 235, 255},
	color.RGBA{250, 192, 170, 255},
	color.RGBA{253, 214, 181, 255},
	color.RGBA{253, 243, 184, 255},
	color.RGBA{254, 128, 254, 255},
	color.RGBA{254, 254, 102, 255},
	color.RGBA{254, 254, 196, 255},
	color.RGBA{255, 106, 213, 255},
	color.RGBA{255, 224, 241, 255},
	color.RGBA{255, 255, 255, 255},
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func randomColor() color.Color {
	return colors[r.Intn(len(colors))]
}

func randomContrastingColor(c1 color.Color) color.Color {
	contrastMinimum := 1.6

	c2 := randomColor()
	contrast := getContrast(c1, c2)

	for contrast <= contrastMinimum {
		fmt.Printf("contrast %v was less than %v\n", contrast, contrastMinimum)

		c2 = randomColor()
		contrast = getContrast(c1, c2)
	}

	fmt.Printf("final contrast %v\n", contrast)

	return c2
}

type Shape struct {
	ln.Mesh
}

func (s *Shape) Paths() ln.Paths {
	var result ln.Paths

	steps := 400

	for i := 0; i <= steps; i++ {
		p := float64(i) / float64(steps)

		// plane := ln.Plane{ln.Vector{0, 0, p*2 - 1}, ln.Vector{0, 0, 1}}
		// result = append(result, plane.IntersectMesh(&s.Mesh)...)

		// plane := ln.Plane{ln.Vector{p*2 - 1, 0, 0}, ln.Vector{1, 0, 0}}
		// result = append(result, plane.IntersectMesh(&s.Mesh)...)

		plane := ln.Plane{ln.Vector{0, p*2 - 1, 0}, ln.Vector{0, 1, 0}}
		result = append(result, plane.IntersectMesh(&s.Mesh)...)
	}

	return result
}

// this is a fun hack to get around JPEG compression on Twitter
func setTransparentPixel(dc *gg.Context) *gg.Context {
	bounds := dc.Image().Bounds()

	i := image.NewRGBA(bounds)

	draw.Draw(i, bounds, dc.Image(), image.Point{0, 0}, draw.Over)

	i.Set(bounds.Max.X-1, bounds.Max.Y-1, color.RGBA{0, 0, 0, 0})

	return gg.NewContextForImage(i)
}

func getLuminance(c color.Color) float64 {
	r, g, b, _ := c.RGBA()

	rgba := [3]float64{float64(r), float64(g), float64(b)}

	for i := 0; i < 3; i++ {
		rgb := rgba[i] / 65535

		if rgb < 0.03928 {
			rgb = rgb / 12.92
		} else {
			rgb = math.Pow((rgb+0.055)/1.055, 2.4)
		}

		rgba[i] = rgb
	}

	return 0.2126*rgba[0] + 0.7152*rgba[1] + 0.0722*rgba[2]
}

func getContrast(c1, c2 color.Color) float64 {
	l1 := getLuminance(c1) + 0.05
	l2 := getLuminance(c2) + 0.05

	ratio := l1 / l2

	if l2 > l1 {
		ratio = 1 / ratio
	}

	return ratio
}

func main() {
	models, err := filepath.Glob("./converted/*.obj")

	if err != nil {
		panic(err)
	}

	numberMeshes := r.Intn(2) + 2

	// eye := ln.Vector{4, 3, 2}    // camera position
	// center := ln.Vector{0, 0, 0} // camera looks at
	// up := ln.Vector{0, 0, 1}     // up direction

	width := 1024.0
	height := 1024.0

	dc := gg.NewContext(int(width), int(height))

	dc.InvertY()

	backgroundColor := randomColor()

	dc.SetColor(backgroundColor)
	dc.Clear()

	dc.SetLineWidth(1)

	eye := ln.Vector{0.65, 0, 1.2}
	center := ln.Vector{0, 0, 0}
	up := ln.Vector{0, 0, 1}

	fovy := 50.0
	znear := 0.01
	zfar := 100.0
	step := 0.01

	for i := 0; i < numberMeshes; i++ {
		scene := ln.Scene{}

		model := models[r.Intn(len(models))]

		fmt.Printf("loading %v\n", model)

		mesh, err := ln.LoadOBJ(model)

		if err != nil {
			panic(err)
		}

		// mesh.FitInside(
		// 	ln.Box{ln.Vector{-1, -1, -1}, ln.Vector{1, 1, 1}},
		// 	ln.Vector{0.5, 0.5, 0.5})

		mesh.UnitCube()

		dc.SetColor(randomContrastingColor(backgroundColor))

		scene.Add(ln.NewTransformedShape(&Shape{*mesh},
			ln.Rotate(ln.Vector{r.Float64(), r.Float64(), r.Float64()},
				ln.Radians(r.Float64()*360))))

		fmt.Println("rendering")

		paths := scene.Render(eye, center, up, width, height, fovy, znear, zfar, step)

		fmt.Println("drawing")

		for _, path := range paths {
			for _, v := range path {
				dc.LineTo(v.X, v.Y)
			}

			dc.NewSubPath()
		}

		dc.Stroke()
	}

	newContext := setTransparentPixel(dc)

	fmt.Println("saving")

	newContext.SavePNG("out.png")
}
