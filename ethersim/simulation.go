package ethersim

import "fmt"

type Simulation struct {
	components        []NetworkComponent
	fallingComponents []NetworkComponent
}

func MakeSimulation() *Simulation {
	return &Simulation{
		components:        make([]NetworkComponent, 0),
		fallingComponents: make([]NetworkComponent, 0),
	}
}
func (s *Simulation) Tick() {
	fmt.Printf("-------------tick-------------\n")
	for _, c := range s.components {
		c.Tick()
	}
	fmt.Printf("-------------post-------------\n")
	for _, c := range s.fallingComponents {
		c.Tick()
	}
}
func (s *Simulation) register(c NetworkComponent) {
	if c.TickFalling() {
		s.fallingComponents = append(s.fallingComponents, c)
	} else {
		s.components = append(s.components, c)
	}
}
