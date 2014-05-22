package sigmon

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/tedsuo/ifrit"
)

type sigmon struct {
	signals          []os.Signal
	monitoredProcess ifrit.Process
}

func New(p ifrit.Process, signals ...os.Signal) ifrit.Runner {
	signals = append(signals, syscall.SIGINT, syscall.SIGTERM)
	return &sigmon{
		signals:          signals,
		monitoredProcess: p,
	}
}

func (s *sigmon) Run(signals <-chan os.Signal, ready chan<- struct{}) error {
	osSignals := make(chan os.Signal, len(s.signals))
	signal.Notify(osSignals, s.signals...)

	close(ready)

	for {
		select {
		case sig := <-signals:
			s.monitoredProcess.Signal(sig)
		case sig := <-osSignals:
			s.monitoredProcess.Signal(sig)
		case err := <-s.monitoredProcess.Wait():
			return err
		}
	}
}
