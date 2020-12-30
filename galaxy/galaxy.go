// SPDX-License-Identifier: Unlicense OR MIT

// The galaxy command is a seed for producing a demo
// for multiple objects moving in a Gio window.
package main

import (
	"log"
	"time"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/spatial/barneshut"
	"gonum.org/v1/gonum/spatial/r2"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

type mass struct {
	d r2.Vec  // position
	v r2.Vec  // velocity
	m float64 // mass
}

func (m *mass) Coord2() r2.Vec { return m.d }
func (m *mass) Mass() float64  { return m.m }
func (m *mass) move(f r2.Vec) {
	// F = ma
	f.X /= m.m
	f.Y /= m.m
	m.v = m.v.Add(f)

	// Update position.
	m.d = m.d.Add(m.v)
}

func main() {
	rnd := rand.New(rand.NewSource(uint64(time.Now().Unix())))

	// Make 50 stars in random locations and velocities.
	stars := make([]*mass, 50)
	p := make([]barneshut.Particle2, len(stars))
	for i := range stars {
		s := &mass{
			d: r2.Vec{
				X: 100*rnd.Float64() - 50,
				Y: 100*rnd.Float64() - 50,
			},
			m: rnd.Float64(),
		}
		// Aim at the ground and miss.
		s.d = s.d.Scale(-1).Add(r2.Vec{
			X: 10 * rnd.NormFloat64(),
			Y: 10 * rnd.NormFloat64(),
		})

		stars[i] = s
		p[i] = s
	}
	vectors := make([]r2.Vec, len(stars))

	tracks := make([]plotter.XYs, len(stars))

	// Make a plane to calculate approximate forces
	plane := barneshut.Plane{Particles: p}

	// Run a simulation for 10000 updates.
	for i := 0; i < 10000; i++ {
		// Build the data structure. For small systems
		// this step may be omitted and ForceOn will
		// perform the naive quadratic calculation
		// without building the data structure.
		err := plane.Reset()
		if err != nil {
			log.Fatal(err)
		}

		// Calculate the force vectors using the theta
		// parameter.
		const theta = 0.1
		// and an imaginary gravitational constant.
		const G = 10
		for j, s := range stars {
			vectors[j] = plane.ForceOn(s, theta, barneshut.Gravity2).Scale(G)
		}

		// Update positions.
		for j, s := range stars {
			s.move(vectors[j])
			tracks[j] = append(tracks[j], plotter.XY{X: s.d.X, Y: s.d.Y})
		}
	}

	plt, err := plot.New()
	if err != nil {
		log.Fatalf("failed create plot: %v", err)
	}
	for i, t := range tracks {
		l, err := plotter.NewLine(t)
		if err != nil {
			log.Fatalf("failed create track: %v", err)
		}
		l.Color = plotutil.Color(i)
		l.Dashes = plotutil.Dashes(i)
		plt.Add(l)
	}
	plt.X.Min = -1000
	plt.X.Max = 1000
	plt.Y.Min = -1000
	plt.Y.Max = 1000
	err = plt.Save(20*vg.Centimeter, 20*vg.Centimeter, "galaxy.svg")
	if err != nil {
		log.Fatalf("failed to save file: %v", err)
	}
}
