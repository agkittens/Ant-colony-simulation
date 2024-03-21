package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	m *Map
	a *[NUM_ANTS]Ants

	p      *Pheromone
	pixels []byte
}

// updates screen
func (g *Game) Update() error {
	for i := 0; i < (NUM_ANTS - 1); i++ {
		go g.a[i].wander(g.m.food_x, g.m.food_y, &g.m.food_count, (*[SIZE][SIZE]float32)(&g.m.structure), Pheromone{g.p.phero_to_food, g.p.phero_to_nest})

		// if yes then rest of ants can follow phero
		if g.a[i].round != 0 {
			for j := 0; j < (NUM_ANTS); j++ {
				g.a[j].round = 1
			}
		}
		// g.a[i+1].wander(g.m.food_x, g.m.food_y, &g.m.food_count, (*[SIZE][SIZE]float32)(&g.m.structure), Pheromone{g.p.phero_to_food, g.p.phero_to_nest})
	}
	g.p.fade(&g.m.structure)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	if g.pixels == nil {
		g.pixels = make([]byte, SIZE*SIZE*4)
	}

	g.m.Draw(g.pixels)

	screen.WritePixels(g.pixels)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Food in nest: %d", g.m.food_count))

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SIZE, SIZE
}

// draws pixels on points with nest, food and tour
func (m *Map) Draw(pix []byte) {
	for x := 0; x < SIZE; x++ {
		for y := 0; y < SIZE; y++ {

			pos := y*SIZE + x
			val := min(m.structure[x][y], 255)
			pix[pos*4] = byte(val)
		}
	}
	pos_f := m.food_y*SIZE + m.food_x
	pix[pos_f*4+1] = byte(255)
	pos_n := m.nest_y*SIZE + m.nest_x
	pix[pos_n*4+2] = byte(255)
}
