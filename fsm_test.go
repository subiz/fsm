package fsm_test

import (
	"testing"
	"bitbucket.org/subiz/fsm"
)

var st fsm.FSM = fsm.New()

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
	st.Map("init", s.State1)
	st.Map("state2", s.State2)
	st.Run("init", nil)
}
