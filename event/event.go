package event

import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

const (
	Event_EXIT = "exit"
	Event_WAIT = "wait"
)

var (
	Events = make(map[string][]func(), 2)
)

// ------------  事件管理 ---------------

func OnEvent(name string, fs ...func()) error {
	evs, ok := Events[name]
	if !ok {
		evs = make([]func(), 0, len(fs))
	}

	for _, f := range fs {
		fp := reflect.ValueOf(f).Pointer()
		for i := 0; i < len(evs); i++ {
			if reflect.ValueOf(evs[i]).Pointer() == fp {
				return fmt.Errorf("func[%v] already exists in event[%s]", fp, name)
			}
		}
		evs = append(evs, f)
	}
	Events[name] = evs
	return nil
}

func EmitEvent(name string) {
	evs, ok := Events[name]
	if !ok {
		return
	}

	for _, f := range evs {
		f()
	}
}

func EmitAllEvents() {
	for _, fs := range Events {
		for _, f := range fs {
			f()
		}
	}
	return
}

func OffEvent(name string, f func(interface{})) error {
	evs, ok := Events[name]
	if !ok || len(evs) == 0 {
		return fmt.Errorf("envet[%s] doesn't have any funcs", name)
	}

	fp := reflect.ValueOf(f).Pointer()
	for i := 0; i < len(evs); i++ {
		if reflect.ValueOf(evs[i]).Pointer() == fp {
			evs = append(evs[:i], evs[i+1:]...)
			Events[name] = evs
			return nil
		}
	}

	return fmt.Errorf("%v func dones't exist in event[%s]", fp, name)
}

func OffAllEvents(name string) error {
	Events[name] = nil
	return nil
}

func WaitEvent(sig ...os.Signal) os.Signal {
	c := make(chan os.Signal, 1)
	if len(sig) == 0 {
		signal.Notify(c, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT)
	} else {
		signal.Notify(c, sig...)
	}
	return <-c
}
