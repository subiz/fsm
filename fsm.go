package fsm

import "sync"

type FSM interface {
	Map(string, func(string, []interface{}) (string, []interface{}))
	Run(string, []interface{})
	Stop()
}

type FSMI struct {
	statex map[string]func(string, []interface{}) (string, []interface{})
	state string
	lock *sync.Mutex
}

const (
	RUNNING = "running"
	STOPPED = "stopped"
)

func New() *FSMI {
	return &FSMI{
		lock: &sync.Mutex{},
		state: STOPPED,
		statex: make(map[string]func(string, []interface{}) (string, []interface{})),
	}
}

func (f *FSMI) Stop() {
	f.lock.Lock()
	f.state = STOPPED
	f.lock.Unlock()
}

func (f *FSMI) Map(e string, to func(string, []interface{}) (string, []interface{})) {
	f.statex[e] = to
}

func (f *FSMI) Run(e string, ps []interface{}) {
	f.state = RUNNING
	for {
		f.lock.Lock()
		if f.state == STOPPED {
			f.lock.Unlock()
			break
		}
		f.lock.Unlock()

		if _, ok := f.statex[e]; !ok {
			panic("unknow state " + e)
		}
		e, ps = f.statex[e](e, ps)
	}
}
