// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"log"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/spatial/barneshut"
	"gonum.org/v1/gonum/spatial/r2"
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

func galaxy(numStars int, rnd *rand.Rand) ([]*mass, barneshut.Plane) {
	// Make 50 stars in random locations and velocities.
	stars := make([]*mass, numStars)
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
	// Make a plane to calculate approximate forces
	plane := barneshut.Plane{Particles: p}

	return stars, plane
}

func simulate(stars []*mass, plane barneshut.Plane, dist *distribution) {
	vectors := make([]r2.Vec, len(stars))
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
	}

	// Recompute the distribution of stars
	dist.Update(stars)
	dist.EnsureSquare()
}
