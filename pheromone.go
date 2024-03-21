package main

type Pheromone struct {
	phero_to_food float32
	phero_to_nest float32
}

func init_phero() Pheromone {
	return Pheromone{
		phero_to_food: 10,
		phero_to_nest: 70,
	}
}

// fading pheromones on whole map
func (p *Pheromone) fade(structure *[SIZE][SIZE]float32) {
	for w := 0; w < len(structure[0]); w++ {
		for h := 0; h < len(structure[1]); h++ {
			structure[w][h] = structure[w][h] - FADE_VALUE
			structure[w][h] = max(structure[w][h], 0)
		}
	}
}
