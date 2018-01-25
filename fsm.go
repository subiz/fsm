package fsm

import (
	"sync"
	"time"
	"fmt"
)

type FSM interface {
	Map([]State)
	Run(string, interface{})
	Stop()
	GetState() string
}

type State struct {
	Event string
	State string
	Function func(string, interface{}) (string, interface{})
}

type FSMI struct {
	current_state string
	statex map[string]State
	state string
	lock *sync.Mutex
	id string
}

const (
	RUNNING = "running"
	STOPPED = "stopped"
)

func New(id string) *FSMI {
	return &FSMI{
		id: id,
		lock: &sync.Mutex{},
		state: STOPPED,
		statex: make(map[string]State),
	}
}

func (f *FSMI) Stop() {
	f.lock.Lock()
	f.state = STOPPED
	f.lock.Unlock()
}

func (f *FSMI) GetState() string {
	return f.current_state
}

func (f *FSMI) Map(ss []State) {
	for _, s := range ss {
		f.statex[s.Event] = s
	}
}

func (f *FSMI) Run(e string, ps interface{}) {
	f.state = RUNNING
	defer f.Stop()
	for {
		f.lock.Lock()
		if f.state == STOPPED {
			f.lock.Unlock()
			break
		}
		f.lock.Unlock()

		if _, ok := f.statex[e]; !ok {
			panic("unknow state " + e + ".")
		}
		func() {
			defer func() {
				r := recover()
				if r != nil {
					fmt.Printf("state error, id = %s; err%v\n", f.id, r)
					time.Sleep(5 * time.Second)
				}
			}()
			f.current_state = f.statex[e].State
			e, ps = f.statex[e].Function(e, ps)
		}()
	}
}
