package textwire

import (
	"fmt"
	"os"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

type WatcherLogger struct{}

func NewWatcherLogger() *WatcherLogger {
	return &WatcherLogger{}
}

func (l *WatcherLogger) Info(text string) {
	fmt.Printf("%s[Watcher]%s %s(info)%s %s\n", colorYellow, colorReset, colorCyan, colorReset, text)
}

func (l *WatcherLogger) Success(text string) {
	fmt.Printf("%s[Watcher]%s %s(good)%s %s\n", colorYellow, colorReset, colorGreen, colorReset, text)
}

func (l *WatcherLogger) Error(text string) {
	fmt.Printf("%s[Watcher]%s %s(error)%s %s\n", colorYellow, colorReset, colorRed, colorReset, text)
}

func (l *WatcherLogger) Fatal(text string) {
	l.Error(text)
	os.Exit(1)
}
