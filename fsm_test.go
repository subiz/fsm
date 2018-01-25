package fsm_test

import (
	"testing"
	"bitbucket.org/subiz/fsm"
)

var st fsm.FSM = fsm.New("van")

type S struct {}

func (s S) State1(e string, ps interface{}) (string, interface{}) {
	return "state2", nil
}

func (s S) State2(e string, ps interface{}) (string, interface{}) {
	st.Stop()
	return "", 4
}

func TestFSM(t *testing.T) {
	s := &S{}
	st.Map([]fsm.State{
		fsm.State{"init", "STATE1", s.State1},
		fsm.State{"state2", "STATE2", s.State2},
	})
	st.Run("init", nil)
}
