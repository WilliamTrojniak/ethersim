package ethersim

type EventCb func(id int)
type MsgEventCb func(id int, msg NetworkMsg)

type Simulation struct {
	components        []NetworkComponent
	fallingComponents []NetworkComponent

	onTransceiverBeginTransmit MsgEventCb
	onTransceiverEndTransmit   MsgEventCb
	onTransceiverJam           EventCb
	onDeviceQueueMsg           MsgEventCb
	onDeviceReceiveMsg         MsgEventCb
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

func (s *Simulation) SetTransceiverBeginTransmitCb(f MsgEventCb) { s.onTransceiverBeginTransmit = f }
func (s *Simulation) SetTransceiverEndTransmitCb(f MsgEventCb)   { s.onTransceiverEndTransmit = f }
func (s *Simulation) SetTransceiverJamCb(f EventCb)              { s.onTransceiverJam = f }
func (s *Simulation) SetDeviceQueueMsgCb(f MsgEventCb)           { s.onDeviceQueueMsg = f }
func (s *Simulation) SetDeviceReceiveMsgCb(f MsgEventCb)         { s.onDeviceReceiveMsg = f }
