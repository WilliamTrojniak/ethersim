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
	// if d.network.isResetting(d) {
	// 	if !d.seenReset && d.hasSent {
	// 		d.timeoutRange *= 2
	// 		d.seenReset = true
	// 	}
	// 	d.randomizeTimeout()
	// 	return
	// }
	//
	// d.seenReset = false
	//
	// if d.timeout > 0 {
	// 	d.timeout--
	// 	return
	// }

	if len(d.queuedMessages) > 0 && !d.network.incomingMsg(d) && !d.network.isResetting(d) {
		msg := d.queuedMessages[0]
		d.queuedMessages = d.queuedMessages[1:]
		d.network.OnMsg(msg, d)
		// d.hasSent = true
	}
}

// Expects to be called during rising edge of tick
func (d *NetworkDevice) OnMsg(msg NetworkMsg, sender Network) {
	// if !msg.Valid() {
	// 	if !d.seenReset && d.hasSent {
	// 		d.timeoutRange *= 2
	// 		d.seenReset = true
	// 	}
	// 	d.resetTicks = numResetTicks
	// } else {
	// 	if d.hasSent {
	// 		d.timeoutRange = max(1, int(float32(d.timeoutRange)*0.9))
	// 	}
	// }
	// d.randomizeTimeout()
	// d.hasSent = false

	d.QueueMessage(&BaseMsg{V: true})
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
