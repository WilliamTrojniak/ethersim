package ethersim

import (
	"math/rand/v2"
)

var deviceid int = 0
var numResetTicks int = 10

// Devices
type NetworkDevice struct {
	network        Network
	id             int
	queuedMessages []NetworkMsg
	resetTicks     int
	timeout        int
}

func (n *NetworkNode) CreateDevice(weight int) (*NetworkDevice, *NetworkEdge) {
	d := &NetworkDevice{
		id:             deviceid,
		queuedMessages: make([]NetworkMsg, 0),
		network:        nil,
		resetTicks:     0,
	}
	deviceid++
	edge := makeNetworkEdge(n.sim, n, d, weight)
	n.edges = append(n.edges, edge)
	d.network = edge
	n.sim.register(d)
	return d, edge
}
func (d *NetworkDevice) Id() int           { return d.id }
func (d *NetworkDevice) TickFalling() bool { return true }
func (d *NetworkDevice) Tick() {
	if d.resetTicks > 0 {
		d.resetTicks--
		return
	}

	if d.network.isResetting(d) {
		d.randomizeTimeout()
		return
	}

	if d.timeout > 0 {
		d.timeout--
		return
	}

	if len(d.queuedMessages) > 0 && !d.network.incomingMsg(d) && !d.network.isResetting(d) {
		msg := d.queuedMessages[0]
		d.queuedMessages = d.queuedMessages[1:]
		d.network.OnMsg(msg, d)
	}
}

// Expects to be called during rising edge of tick
func (d *NetworkDevice) OnMsg(msg NetworkMsg, sender Network) {
	if !msg.Valid() {
		d.resetTicks = numResetTicks
	}

	d.QueueMessage(&BaseMsg{V: true})
}

func (d *NetworkDevice) incomingMsg(Network) bool { return false }
func (d *NetworkDevice) IncomingMsg() bool        { return d.network.incomingMsg(d) }
func (d *NetworkDevice) QueueMessage(msg NetworkMsg) {
	if len(d.queuedMessages) < 10 {
		d.queuedMessages = append(d.queuedMessages, msg)
	}
}

func (d *NetworkDevice) isResetting(from Network) bool {
	return d.resetTicks > 0
}

func (d *NetworkDevice) randomizeTimeout() {
	d.timeout = rand.IntN(500)

}
