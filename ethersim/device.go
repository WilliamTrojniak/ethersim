package ethersim

import (
// "math/rand/v2"
)

var deviceid int = 0
var numResetTicks int = 40

// Devices
type NetworkDevice struct {
	network        Network
	id             int
	queuedMessages []NetworkMsg
}

func (n *NetworkNode) CreateDevice(weight int) (*NetworkDevice, *NetworkEdge) {
	if n.deviceEdge != nil {
		return nil, nil
	}
	d := &NetworkDevice{
		id:             deviceid,
		queuedMessages: make([]NetworkMsg, 0),
		network:        nil,
	}
	deviceid++
	edge := makeNetworkEdge(n.sim, n, d, weight)
	n.deviceEdge = edge
	d.network = edge
	n.sim.register(d)
	// d.randomizeTimeout()
	return d, edge
}
func (d *NetworkDevice) Id() int           { return d.id }
func (d *NetworkDevice) TickFalling() bool { return true }
func (d *NetworkDevice) Tick() {

	if len(d.queuedMessages) > 0 && !d.network.incomingMsg(d) && !d.network.isResetting(d) {
		msg := d.queuedMessages[0]
		d.queuedMessages = d.queuedMessages[1:]
		d.network.OnMsg(msg, d)
	}
}

// Expects to be called during rising edge of tick
func (d *NetworkDevice) OnMsg(msg NetworkMsg, sender Network) {
}

func (d *NetworkDevice) incomingMsg(Network) bool { return false }
func (d *NetworkDevice) IncomingMsg() bool        { return d.network.incomingMsg(d) }
func (d *NetworkDevice) QueueMessage(msg NetworkMsg) {

	if len(d.queuedMessages) < 100 {
		d.queuedMessages = append(d.queuedMessages, msg)
	}
}

func (d *NetworkDevice) isResetting(from Network) bool {
	return false
}

func (d *NetworkDevice) QueuedMessages() []NetworkMsg { return d.queuedMessages }
func (d *NetworkDevice) Timeout() int                 { return 0 }
func (d *NetworkDevice) TimeoutFrom() int             { return 1 }
func (d *NetworkDevice) TimeoutRange() int            { return 1 }
func (d *NetworkDevice) SeenT() bool                  { return false }
