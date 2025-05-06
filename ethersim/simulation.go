package ethersim

type EventCb func(id int)

type Simulation struct {
	components        []NetworkComponent
	fallingComponents []NetworkComponent

	onTransceiverBeginTransmit EventCb
	onTransceiverEndTransmit   EventCb
	onTransceiverJam           EventCb
}

func MakeSimulation() *Simulation {
	return &Simulation{
		components:        make([]NetworkComponent, 0),
		fallingComponents: make([]NetworkComponent, 0),
	}
}
func (s *Simulation) Tick() {
	// fmt.Printf("-------------tick-------------\n")
	for _, c := range s.components {
		c.Tick()
	}
	// fmt.Printf("-------------post-------------\n")
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

func (s *Simulation) SetTransceiverBeginTransmitCb(f EventCb) { s.onTransceiverBeginTransmit = f }
func (s *Simulation) SetTransceiverEndTransmitCb(f EventCb)   { s.onTransceiverEndTransmit = f }
func (s *Simulation) SetTransceiverJamCb(f EventCb)           { s.onTransceiverJam = f }
