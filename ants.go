package main

import (
	"math"
	"math/rand"
)

type Ants struct {
	home_x, home_y int
	pos_x, pos_y   int
	cargo          int
	direction      int
	slot           int
	round          int
	step           int
}

func new_ant(nest_x int, nest_y int) Ants {
	return Ants{
		home_x: nest_x,
		home_y: nest_y,
		pos_x:  nest_x,
		pos_y:  nest_y,
		cargo:  0,
		round:  0,
		step:   0,
	}

}

func (a *Ants) wander(food_x int, food_y int, count *int, structure *[SIZE][SIZE]float32, p Pheromone) {

	var des_x, des_y int
	prob := rand.Intn(10)

	//chosing direction to go
	if a.cargo == 1 && prob <= 6 {
		des_x = a.home_x
		des_y = a.home_y
		a.direction = a.to_des(des_x, des_y)
	} else if a.cargo == 0 && a.sense_food2(food_x, food_y) {
		a.direction = a.to_des(food_x, food_y)
	} else {
		a.direction = a.calc_dir(a.direction+a.slot-1, 4)
	}

	if a.round != 0 && a.step%10 == 0 && a.pos_x > (MIN+2) && a.pos_y > (MIN+2) && a.pos_x < (MAX-2) && a.pos_y < (MAX-2) {
		a.slot = a.follow_phero(structure)

	} else {
		a.slot = rand.Intn(3)
	}

	a.sense_food(food_x, food_y)

	a.move()

	a.control_movement()

	//leaving pheromones on slot
	if a.cargo == 0 {
		structure[a.pos_x][a.pos_y] = structure[a.pos_x][a.pos_y] + p.phero_to_food
	} else {
		structure[a.pos_x][a.pos_y] = structure[a.pos_x][a.pos_y] + p.phero_to_nest
	}

	if a.check_location(food_x, food_y) {
		*count = *count + 1
	}
	a.step += 1

}

// checks if ant is in nest or on food, when ant visits nest with cargo it drops it there
// returns bool value, when true, food counter in nest rises
func (a *Ants) check_location(food_x int, food_y int) bool {
	if a.pos_x == food_x && a.pos_y == food_y {
		a.cargo = 1
		a.round = a.round + 1
		return false
	}

	if a.pos_x == a.home_x && a.pos_y == a.home_y {
		if a.cargo == 1 {
			a.cargo = 0
			return true
		}

	}
	return false

}

// calcs movement for ant based on chosen direction and slot
// ant can peak 1 from 4 directions and 1 from 3 slots in front of them
//
//				d0
//			s0	s1	s2
//		s2				s0
//	d3	s1		A		s1	d1
//		s0				s2
//			s2	s1	s0
//				d2
func (a *Ants) move() {
	if a.direction == 0 {
		a.pos_x = a.pos_x - 1 + a.slot
		a.pos_y = a.pos_y - 1

	} else if a.direction == 1 {
		a.pos_x = a.pos_x + 1
		a.pos_y = a.pos_y - 1 + a.slot

	} else if a.direction == 2 {
		a.pos_x = a.pos_x + 1 - a.slot
		a.pos_y = a.pos_y + 1

	} else if a.direction == 3 {
		a.pos_x = a.pos_x - 1
		a.pos_y = a.pos_y + 1 - a.slot
	}
}

// controls ants for not leaving boundries
func (a *Ants) control_movement() {
	pos_x_temp := a.pos_x
	pos_y_temp := a.pos_y

	a.pos_x = min(max(a.pos_x, MIN), MAX)
	a.pos_y = min(max(a.pos_y, MIN), MAX)

	if pos_x_temp != a.pos_x || pos_y_temp != a.pos_y {
		a.direction = (a.direction + 2) % 4
	}

}

// calcs next ant direction
func (a *Ants) calc_dir(d, m int) int {
	var res int = d % m
	if (res < 0 && m > 0) || (res > 0 && m < 0) {
		return res + m
	}
	return res
}

// checks if ant has food on 1 of 3 slots in front of them, if yes returns true
func (a *Ants) sense_food(food_x int, food_y int) bool {
	orientation := 1
	var x_temp, y_temp, x_temp2, y_temp2 int

	if a.direction == 2 || a.direction == 3 {
		orientation = -1
	}

	if a.direction == 0 || a.direction == 2 {
		y_temp = a.pos_y + orientation*(-1)
		for i := 0; i < 3; i++ {
			x_temp = a.pos_x + orientation*(-1+i)
			if x_temp == food_x && y_temp == food_y {
				a.pos_x, a.pos_y = food_x, food_y
				return true
			}
		}
	} else {
		x_temp2 = a.pos_x + orientation*(-1)
		for i := 0; i < 3; i++ {
			y_temp2 = a.pos_y + orientation*(-1+i)
			if x_temp2 == food_x && y_temp2 == food_y {
				a.pos_x, a.pos_y = food_x, food_y
				return true
			}
		}
	}
	return false

}

// makes ants to follow phero based on values in front of ant
func (a *Ants) follow_phero(structure *[SIZE][SIZE]float32) int {
	orientation := 1
	pheromones := []int{}
	idx_max := 0
	idx_min := 0
	idx_prob := 0
	prob := rand.Intn(10)

	var x_temp, y_temp int

	if a.direction == 2 || a.direction == 3 {
		orientation = -1
	}

	if a.direction == 0 || a.direction == 2 {
		y_temp = a.pos_y + orientation*(-1)
		for i := 0; i < 3; i++ {
			x_temp = a.pos_x + orientation*(-1+i)
			pheromones = append(pheromones, int(structure[x_temp][y_temp]))

		}
	} else {
		x_temp = a.pos_x + orientation*(-1)
		for i := 0; i < 3; i++ {
			y_temp = a.pos_y + orientation*(-1+i)
			pheromones = append(pheromones, int(structure[x_temp][y_temp]))
		}
	}

	phero_max := pheromones[0]
	phero_min := pheromones[0]

	for i := 1; i < 3; i++ {
		if pheromones[i] > phero_max {
			idx_max = i
		}
		if pheromones[i] < phero_min {
			idx_min = i
		}
	}

	//adds probability for chosing a slot
	if prob <= 10 {
		idx_prob = idx_max
	} else if prob <= 4 {
		idx_prob = (2 * (idx_max + idx_min)) % 3
	} else if prob <= 2 {
		idx_prob = idx_min
	}
	return idx_prob
}

// makes ants to change direction based on their destination
func (a *Ants) to_des(des_x int, des_y int) int {
	dir := 0

	if a.pos_x > des_x && a.pos_y > des_y {
		dir = 3 * rand.Intn(2)

	} else if a.pos_x > des_x && a.pos_y < des_y {
		dir = rand.Intn(2) + 2

	} else if a.pos_x < des_x && a.pos_y > des_y {
		dir = rand.Intn(2)

	} else if a.pos_x < des_x && a.pos_y < des_y {
		dir = rand.Intn(2) + 1

	} else if a.pos_x == des_x && a.pos_y > des_y {
		dir = 0
	} else if a.pos_x == des_x && a.pos_y < des_y {
		dir = 2
	} else if a.pos_y == des_y && a.pos_x > des_x {
		dir = 3
	} else if a.pos_y == des_y && a.pos_x < des_x {
		dir = 1
	}
	return dir
}

// checks if ant see food
func (a *Ants) sense_food2(food_x int, food_y int) bool {
	distance := math.Sqrt(math.Pow(float64(a.pos_x-food_x), 2) + math.Pow(float64(a.pos_y-food_y), 2))
	if distance <= 15 {
		return true
	}
	return false
}
