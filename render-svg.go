package main

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/fogleman/ln/ln"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

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

func main() {
	models, err := filepath.Glob("./converted/*.obj")

	if err != nil {
		panic(err)
	}

	width := 1024.0
	height := 1024.0

	eye := ln.Vector{0.65, 0, 1.2}
	center := ln.Vector{0, 0, 0}
	up := ln.Vector{0, 0, 1}

	fovy := 50.0
	znear := 0.01
	zfar := 100.0
	step := 0.01

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

	scene.Add(ln.NewTransformedShape(&Shape{*mesh},
		ln.Rotate(ln.Vector{r.Float64(), r.Float64(), r.Float64()},
			ln.Radians(r.Float64()*360))))

	fmt.Println("rendering")

	paths := scene.Render(eye, center, up, width, height, fovy, znear, zfar, step)

	fmt.Println("writing")

	paths.WriteToSVG("./out.svg", width, height)
}
