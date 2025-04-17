package ethersim

import (
	"fmt"
)

// Network
var (
	nodeid   int = 0
	deviceid int = 0
	edgeid   int = 0
)

type NetworkMsg interface {
	Valid() bool
	Invalid()
	Copy() NetworkMsg
}

type Network interface {
	Update() error
	PostUpdate() error
	OnMsg(NetworkMsg, Network)
	Reset()
	Id() int
	IncomingMsg(dest Network) bool
}

type msgdata struct {
	msg   NetworkMsg
	stage int // between 0 and weight of edge incl
	dir   int // 1 or -1
}

type NetworkEdge struct {
	id       int
	n1       Network
	n2       Network
	edge     bool
	weight   int
	messages []*msgdata
}

func makeNetworkEdge(n1 Network, n2 Network, w int) *NetworkEdge {
	id := edgeid
	edgeid++
	return &NetworkEdge{
		id:       id,
		n1:       n1,
		n2:       n2,
		edge:     false,
		weight:   w,
		messages: make([]*msgdata, 0),
	}
}
func (e *NetworkEdge) Id() int { return e.id }
func (e *NetworkEdge) OnMsg(msg NetworkMsg, from Network) {
	if from != e.n1 && from != e.n2 {
		return
	}
	var dir, start int
	if from == e.n1 {
		dir = 1
		start = 0
	} else {
		dir = -1
		start = e.weight
	}

	for _, m2 := range e.messages {
		if m2.stage == start {
			m2.msg.Invalid()
			msg.Invalid()
		}
	}

	e.messages = append(e.messages, &msgdata{
		msg:   msg,
		stage: start,
		dir:   dir,
	})
}
func (e *NetworkEdge) Update() error {
	dirs := make(map[int]int)
	for _, msg := range e.messages {
		if v, ok := dirs[msg.stage]; !ok {
			dirs[msg.stage] = msg.dir
		} else if v != msg.dir {
			dirs[msg.stage] = 0
		}
	}

	for _, msg := range e.messages {
		if !msg.msg.Valid() {
			continue
		}
		// Case 1: Two messages will swap stages
		if v, ok := dirs[msg.stage+msg.dir]; ok && v != msg.dir {
			msg.msg.Invalid()
		}
		// Case 2: Two messages will be at the same stage
		if v, ok := dirs[msg.stage+2*msg.dir]; ok && v != msg.dir {
			msg.msg.Invalid()
		}
	}

	for i := len(e.messages) - 1; i >= 0; i-- {
		e.messages[i].stage += e.messages[i].dir
		msg := e.messages[i]
		if msg.stage <= 0 {
			e.messages[i] = e.messages[len(e.messages)-1]
			e.messages = e.messages[:len(e.messages)-1]

			e.n1.OnMsg(msg.msg, e.n2)
		} else if msg.stage >= e.weight {
			e.messages[i] = e.messages[len(e.messages)-1]
			e.messages = e.messages[:len(e.messages)-1]

			e.n2.OnMsg(msg.msg, e.n1)
		}
	}

	return nil
}
func (e *NetworkEdge) PostUpdate() error { return nil }
func (e *NetworkEdge) Reset()            {}
func (e *NetworkEdge) IncomingMsg(dest Network) bool {
	dir := 0
	if dest == e.n1 {
		dir = -1
	} else if dest == e.n2 {
		dir = 1
	}
	for _, msg := range e.messages {
		if msg.dir == dir {
			return true
		}
	}
	if dest == e.n1 && e.n2.IncomingMsg(e) {
		return true
	} else if dest == e.n2 && e.n1.IncomingMsg(e) {
		return true
	}

	return false
}

type incMessage struct {
	m    NetworkMsg
	from Network
}

type NetworkNode struct {
	id          int
	edges       []*NetworkEdge
	incMessages []incMessage
}

func MakeNetworkNode() *NetworkNode {
	n := &NetworkNode{
		id:    nodeid,
		edges: make([]*NetworkEdge, 0),
	}
	nodeid++
	return n
}

func (n *NetworkNode) OnMsg(msg NetworkMsg, from Network) {
	n.incMessages = append(n.incMessages, incMessage{
		m:    msg,
		from: from,
	})
}

func (n *NetworkNode) Id() int {
	return n.id
}

func (n *NetworkNode) Update() error {
	for _, edge := range n.edges {
		if edge.n1 == n {
			edge.Update()
		}
	}
	return nil
}

func (n *NetworkNode) PostUpdate() error {
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

	return nil
}
func (n *NetworkNode) IncomingMsg(dest Network) bool {
	for _, edge := range n.edges {
		if edge != dest && edge.IncomingMsg(n) {
			return true
		}
	}
	return false
}
func (n *NetworkNode) Reset() {}
func (n *NetworkNode) CreateNode(weight int) *NetworkNode {
	nn := MakeNetworkNode()
	edge := makeNetworkEdge(n, nn, weight)
	n.edges = append(n.edges, edge)
	nn.edges = append(nn.edges, edge)
	return nn
}

// Devices
type Device struct {
	network Network
	id      int
}

func (n *NetworkNode) createDevice(weight int) *Device {
	d := &Device{id: deviceid}
	deviceid++
	edge := makeNetworkEdge(n, d, weight)
	n.edges = append(n.edges, edge)
	d.network = edge
	return d
}

func (d *Device) SendPacket(msg NetworkMsg) {
	d.network.OnMsg(msg, d)
}

type BaseMsg struct {
	V bool
}

func (m *BaseMsg) Valid() bool      { return m.V }
func (m *BaseMsg) Invalid()         { m.V = false }
func (m *BaseMsg) Copy() NetworkMsg { return &BaseMsg{V: m.V} }

func (d *Device) OnMsg(msg NetworkMsg, sender Network) {
	fmt.Printf("(%v) Received msg, valid %v\n", d.id, msg.Valid())
	if !d.network.IncomingMsg(d) {
		d.SendPacket(&BaseMsg{V: true})
	}
}

func (d *Device) Id() int                  { return d.id }
func (d *Device) Update() error            { return nil }
func (d *Device) PostUpdate() error        { return nil }
func (d *Device) Reset()                   {}
func (d *Device) IncomingMsg(Network) bool { return false }
