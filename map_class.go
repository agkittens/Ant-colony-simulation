package main

import (
	"math"
	"math/rand"
)

type Map struct {
	structure      [SIZE][SIZE]float32
	nest_x, nest_y int
	food_x, food_y int
	food_count     int
}

// generates random location for nest and random location for food
func (m *Map) generate_map() {
	m.nest_x = rand.Intn(SIZE)
	m.nest_y = rand.Intn(SIZE)
	m.structure[m.nest_x][m.nest_y] = 1

	m.food_x = rand.Intn(SIZE)
	m.food_y = rand.Intn(SIZE)
	distance := math.Sqrt(math.Pow(float64(m.nest_x-m.food_x), 2) + math.Pow(float64(m.nest_y-m.food_y), 2))

	if m.nest_x != m.food_x && m.nest_y != m.food_y && distance > 30 {
		m.structure[m.food_x][m.food_y] = 1
	} else {
		m.generate_map()
	}
}
