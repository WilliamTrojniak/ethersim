package ethersim

import "math/rand/v2"

var nodeid int = 0

type incMessage struct {
	m    NetworkMsg
	from Network
}

type NetworkNode struct {
	sim          *Simulation
	id           int
	edges        []*NetworkEdge
	deviceEdge   *NetworkEdge
	incMessages  []incMessage
	outMessages  []NetworkMsg
	resetting    int
	transmitting bool
	resetTicks   int
	timeout      int
	timeoutRange int
	timeoutFrom  int
	seenReset    bool
	hasSent      bool
	transmitRem  int
}

func MakeNetworkNode(s *Simulation) *NetworkNode {
	n := &NetworkNode{
		sim:          s,
		id:           nodeid,
		edges:        make([]*NetworkEdge, 0),
		deviceEdge:   nil,
		resetting:    0,
		transmitting: false,
		resetTicks:   0,
		timeoutRange: 20,
		seenReset:    false,
		hasSent:      false,
	}
	s.register(n)
	nodeid++
	return n
}

func (n *NetworkNode) CreateNode(weight int) (*NetworkNode, *NetworkEdge) {
	nn := MakeNetworkNode(n.sim)
	edge := makeNetworkEdge(n.sim, n, nn, weight)
	n.edges = append(n.edges, edge)
	nn.edges = append(nn.edges, edge)
	return nn, edge
}

func (n *NetworkNode) Id() int { return n.id }

// Network Component Interface
// Distribute messages to edges after edges have ticked
func (n *NetworkNode) TickFalling() bool { return true }
func (n *NetworkNode) Tick() {
	if len(n.incMessages) > 1 || (len(n.incMessages) == 1 && n.transmitting) {
		if !n.seenReset {
			n.seenReset = true
			if n.transmitting {
				n.timeoutRange *= 2
			}
			n.resetTicks = numResetTicks
		}
		n.transmitting = false
	}

	if len(n.incMessages) > 0 {
		n.randomizeTimeout()
	}

	if n.resetTicks > 0 {
		n.resetTicks--
	} else if n.timeout > 0 && len(n.outMessages) > 0 && !n.seenReset {
		n.timeout--
	} else if n.timeout == 0 && len(n.outMessages) > 0 && !n.transmitting {
		n.transmitting = true
		n.transmitRem = 50
	}

	if n.resetTicks == 0 {
		n.seenReset = false
	}

	for _, edge := range n.edges {
		if n.resetTicks > 0 {
			edge.OnMsg(&JamMsg{}, n)
		} else if n.transmitting {
			edge.OnMsg(n.outMessages[0].Copy(), n)
			n.transmitRem--
			if n.transmitRem <= 0 {
				n.transmitting = false
				n.outMessages = n.outMessages[1:]
				n.timeoutRange = int(float32(n.timeoutRange)*0.9) + 2

				n.randomizeTimeout()
			}
		} else if len(n.incMessages) > 0 {
			msg := n.incMessages[0]
			if edge.n1 != msg.from && edge.n2 != msg.from {
				edge.OnMsg(n.incMessages[0].m.Copy(), n)
			}
		}
	}
	n.incMessages = n.incMessages[:0]
}

// OnMsg expects to be called during the rising tick
func (n *NetworkNode) OnMsg(msg NetworkMsg, from Network) {
	if n.deviceEdge != nil && from == n.deviceEdge.n2 {
		n.outMessages = append(n.outMessages, msg)
		return
	}

	n.incMessages = append(n.incMessages, incMessage{
		m:    msg,
		from: from,
	})
}

// Valid during both rising and falling ticks
func (n *NetworkNode) incomingMsg(dest Network) bool {
	for _, msg := range n.incMessages {
		if msg.from != dest {
			return true
		}
	}

	return false
}

func (n *NetworkNode) isResetting(from Network) bool {
	for _, edge := range n.edges {
		if edge != from && edge.isResetting(n) {
			return true
		}

	}
	return false
}

func (n *NetworkNode) IsResetting() bool {
	return n.seenReset
}
func (n *NetworkNode) randomizeTimeout() {
	n.timeout = rand.IntN(n.timeoutRange) + 1
	n.timeoutFrom = n.timeout
}

func (n *NetworkNode) TimeoutRange() int    { return n.timeoutRange }
func (n *NetworkNode) TimeoutFrom() int     { return n.timeoutFrom }
func (n *NetworkNode) Timeout() int         { return n.timeout }
func (n *NetworkNode) NQueued() int         { return len(n.outMessages) }
func (n *NetworkNode) IsTransmitting() bool { return n.transmitting }
func (n *NetworkNode) SendingValue() string {
	if n.transmitting {
		return n.outMessages[0].Value()
	}
	return ""
}
