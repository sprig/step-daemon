package physics

import (
	"fmt"

	"github.com/colinrgodsey/step-daemon/vec"
)

type Move struct {
	from, to   vec.Vec4
	delta, vel vec.Vec4
	fr, time   float64
}

func NewMove(from, to vec.Vec4, fr float64) (m Move) {
	m.from = from
	m.to = to
	m.fr = fr
	m.delta = to.Sub(from).Cache(true)
	m.time = m.delta.Dist() / fr
	m.vel = m.delta.Div(m.time).Cache(true)
	return
}

func (m *Move) From() vec.Vec4 {
	return m.from
}

func (m *Move) To() vec.Vec4 {
	return m.to
}

func (m *Move) Fr() float64 {
	return m.fr
}

func (m *Move) NonEmpty() bool {
	return m.fr > 0 && m.delta.Dist() > 0
}

func (m *Move) Delta() vec.Vec4 {
	return m.delta
}

func (m *Move) Vel() vec.Vec4 {
	return m.vel
}

func (m *Move) Time() float64 {
	return m.time
}

func (m *Move) Scale(scale float64) Move {
	return NewMove(m.from, m.to, m.fr*scale)
}

func (m *Move) String() string {
	return fmt.Sprintf("{%v -> %v (fr: %v)}", m.from, m.to, m.fr)
}

func (m *Move) IsEOrZOnly() bool {
	return m.delta.X() == 0 && m.delta.Y() == 0
}

func (m *Move) IsPrintMove() bool {
	return m.delta.E() > 0 && !m.IsEOrZOnly()
}
