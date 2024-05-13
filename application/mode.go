package main

type Mode struct {
	value int
}

func (m *Mode) next_mode (){
	m.value = m.value + 1
	if m.value > create {
		m.value = 0
	}
}
