package ethersim

import "fmt"

var deviceid int = 0

// Devices
type Device struct {
	network        Network
	id             int
	queuedMessages []NetworkMsg
}

func (n *NetworkNode) createDevice(weight int) *Device {
	d := &Device{
		id:             deviceid,
		queuedMessages: make([]NetworkMsg, 0),
		network:        nil,
	}
	deviceid++
	edge := makeNetworkEdge(n.sim, n, d, weight)
	n.edges = append(n.edges, edge)
	d.network = edge
	n.sim.register(d)
	return d
}
func (d *Device) Id() int           { return d.id }
func (d *Device) TickFalling() bool { return true }
func (d *Device) Tick() {
	if len(d.queuedMessages) > 0 && !d.network.IncomingMsg(d) {
		msg := d.queuedMessages[0]
		d.queuedMessages = d.queuedMessages[1:]
		d.network.OnMsg(msg, d)
	}
}

// Expects to be called during rising edge of tick
func (d *Device) OnMsg(msg NetworkMsg, sender Network) {
	fmt.Printf("(%v) Received msg, valid %v\n", d.id, msg.Valid())
	d.QueueMessage(&BaseMsg{V: true})
}

func (d *Device) IncomingMsg(Network) bool { return false }
func (d *Device) QueueMessage(msg NetworkMsg) {
	if len(d.queuedMessages) < 10 {
		d.queuedMessages = append(d.queuedMessages, msg)
	}
}
