package handler

import (
	"os"

	"github.com/InsideGallery/core/oslistener"
)

// SignalListener contains signal and callbacks
type SignalListener struct {
	callbacks map[os.Signal]func()
}

// NewSignalListener return new signal listener
func NewSignalListener() *SignalListener {
	return &SignalListener{
		callbacks: map[os.Signal]func(){},
	}
}

// Add add signal to listen
func (l *SignalListener) Add(signal os.Signal, fn func()) {
	l.callbacks[signal] = fn
}

// SignalsToSubscribe return list of signals
func (l *SignalListener) SignalsToSubscribe() oslistener.OsSignalsList {
	signals := make(oslistener.OsSignalsList, len(l.callbacks))
	var i int

	for s := range l.callbacks {
		signals[i] = s
		i++
	}

	return signals
}

// ReceiveSignal call when signal received
func (l *SignalListener) ReceiveSignal(s os.Signal) {
	fn, ok := l.callbacks[s]
	if ok {
		fn()
	}
}
