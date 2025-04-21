package ethersim

var edgeid int = 0

type msgdata struct {
	msg   NetworkMsg
	stage int // between 0 and weight of edge incl
	dir   int // 1 or -1
}

func (m *msgdata) Msg() NetworkMsg { return m.msg }
func (m *msgdata) Stage() int      { return m.stage }
func (m *msgdata) Dir() int        { return m.dir }

type NetworkEdge struct {
	id       int
	n1       Network
	n2       Network
	edge     bool
	weight   int
	messages []*msgdata
	incn1    bool
	incn2    bool
}

func makeNetworkEdge(s *Simulation, n1 Network, n2 Network, w int) *NetworkEdge {
	id := edgeid
	edgeid++
	edge := &NetworkEdge{
		id:       id,
		n1:       n1,
		n2:       n2,
		edge:     false,
		weight:   w,
		messages: make([]*msgdata, 0),
		incn1:    false,
		incn2:    false,
	}
	s.register(edge)

	return edge
}

func (e *NetworkEdge) Id() int { return e.id }

func (e *NetworkEdge) TickFalling() bool { return false }
func (e *NetworkEdge) Tick() {
	dirs := make(map[int]int) // Maps stages to directions
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

	e.incn1 = false
	e.incn2 = false

	for i := len(e.messages) - 1; i >= 0; i-- {
		e.messages[i].stage += e.messages[i].dir
		msg := e.messages[i]

		if msg.dir == -1 {
			e.incn1 = true
		} else if msg.dir == 1 {
			e.incn2 = true
		}

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
}

// Expects to be called during falling edge of tick
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

// Valid during rising and falling of tick
func (e *NetworkEdge) incomingMsg(dest Network) bool {
	if dest == e.n1 && (e.incn1 || e.n2.incomingMsg(e)) {
		return true
	} else if dest == e.n2 && (e.incn2 || e.n1.incomingMsg(e)) {
		return true
	}

	return false
}

func (e *NetworkEdge) Weight() int          { return e.weight }
func (e *NetworkEdge) Messages() []*msgdata { return e.messages }
func (e *NetworkEdge) isResetting(from Network) bool {
	if from == e.n1 {
		return e.n2.isResetting(e)
	}

	if from == e.n2 {
		return e.n1.isResetting(e)
	}

	return false
}

func (e *NetworkEdge) IsResetting() bool {
	return e.n1.isResetting(e) || e.n2.isResetting(e)
}
