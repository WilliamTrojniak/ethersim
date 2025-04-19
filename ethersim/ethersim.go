package ethersim

// Network

type NetworkMsg interface {
	Valid() bool
	Invalid()
	Copy() NetworkMsg
}

type NetworkComponent interface {
	Tick()
	TickFalling() bool
}

type Network interface {
	OnMsg(msg NetworkMsg, sender Network)
	Id() int
	IncomingMsg(dest Network) bool
}

type BaseMsg struct {
	V bool
}

func (m *BaseMsg) Valid() bool      { return m.V }
func (m *BaseMsg) Invalid()         { m.V = false }
func (m *BaseMsg) Copy() NetworkMsg { return &BaseMsg{V: m.V} }
