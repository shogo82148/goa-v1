package utils

import (
	"os"
	"os/signal"
)

// Catch signals and invoke then callback
func Catch(signals []os.Signal, then func()) {
	c := make(chan os.Signal, 1)
	if signals == nil {
		signals = defaultSignals
	}
	signal.Notify(c, signals...)
	<-c
	signal.Stop(c)
	if then != nil {
		then()
	}
}
