package ethersim

// Network

type NetworkMsg interface {
	Valid() bool
	Invalid()
	Copy() NetworkMsg
	From() int
	IsJam() bool
	Value() string
	Dest() int
	IsLast() bool
	SetLast()
}

type NetworkComponent interface {
	Tick()
	TickFalling() bool
}

type Network interface {
	OnMsg(msg NetworkMsg, sender Network)
	Id() int
	incomingMsg(dest Network) bool
	isResetting(from Network) bool
}

type BaseMsg struct {
	V      bool
	Msg    string
	Sender int
	Last   bool
	To     int
}

func (m *BaseMsg) Valid() bool   { return m.V }
func (m *BaseMsg) Invalid()      { m.V = false }
func (m *BaseMsg) From() int     { return m.Sender }
func (m *BaseMsg) IsJam() bool   { return false }
func (m *BaseMsg) IsLast() bool  { return m.Last }
func (m *BaseMsg) Value() string { return m.Msg }
func (m *BaseMsg) Dest() int     { return m.To }
func (m *BaseMsg) SetLast()      { m.Last = true }
func (m *BaseMsg) Copy() NetworkMsg {
	return &BaseMsg{
		V:      m.V,
		Sender: m.Sender,
		Msg:    m.Msg,
		To:     m.To,
		Last:   m.Last,
	}
}

type JamMsg struct{ Sender int }

func (m *JamMsg) Valid() bool      { return true }
func (m *JamMsg) Value() string    { return "" }
func (m *JamMsg) Invalid()         {}
func (m *JamMsg) From() int        { return m.Sender }
func (m *JamMsg) IsJam() bool      { return true }
func (m *JamMsg) IsLast() bool     { return true }
func (m *JamMsg) Dest() int        { return -1 }
func (m *JamMsg) SetLast()         {}
func (m *JamMsg) Copy() NetworkMsg { return &JamMsg{Sender: m.Sender} }
