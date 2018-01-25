package fsm

import (
	"sync"
	"time"
	"fmt"
)

type FSM interface {
	Map(string, func(string, interface{}) (string, interface{}))
	Run(string, interface{})
	Stop()
}

type FSMI struct {
	statex map[string]func(string, interface{}) (string, interface{})
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
		statex: make(map[string]func(string, interface{}) (string, interface{})),
	}
}

func (f *FSMI) Stop() {
	f.lock.Lock()
	f.state = STOPPED
	f.lock.Unlock()
}

func (f *FSMI) Map(e string, to func(string, interface{}) (string, interface{})) {
	f.statex[e] = to
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
			e, ps = f.statex[e](e, ps)
		}()
	}
}
