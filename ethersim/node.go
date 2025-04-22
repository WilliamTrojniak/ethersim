package ethersim

var nodeid int = 0

type incMessage struct {
	m    NetworkMsg
	from Network
}

type NetworkNode struct {
	sim         *Simulation
	id          int
	edges       []*NetworkEdge
	incMessages []incMessage
	resetting   int
}

func MakeNetworkNode(s *Simulation) *NetworkNode {
	n := &NetworkNode{
		sim:       s,
		id:        nodeid,
		edges:     make([]*NetworkEdge, 0),
		resetting: 0,
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
	if len(n.incMessages) > 1 {
		for _, msg := range n.incMessages {
			msg.m.Invalid()
		}
	}

	for _, edge := range n.edges {
		for _, msg := range n.incMessages {
			if edge.n1 != msg.from && edge.n2 != msg.from {
				edge.OnMsg(msg.m.Copy(), n)
			}
		}
	}
	n.incMessages = make([]incMessage, 0)
}

// OnMsg expects to be called during the rising tick
func (n *NetworkNode) OnMsg(msg NetworkMsg, from Network) {
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

	// for _, edge := range n.edges {
	// 	if edge != dest && edge.incomingMsg(n) {
	// 		return true
	// 	}
	// }

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
	return n.isResetting(nil)
}
