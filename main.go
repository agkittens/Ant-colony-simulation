package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	m := Map{}
	m.generate_map()
	fmt.Println(m.food_x, m.food_y, m.nest_x, m.nest_y)
	a := [NUM_ANTS]Ants{}

	for i := 0; i < NUM_ANTS; i++ {
		new := new_ant(m.nest_x, m.nest_y)
		a[i] = new
	}

	p := init_phero()

	ebiten.SetWindowSize(WIDTH, HIGHT)
	ebiten.SetWindowTitle("Ants colony")
	ebiten.SetTPS(TPS)

	g := Game{}
	g.m = &m
	g.a = &a
	g.p = &p

	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}

}
