package textwire

import (
	"fmt"
	"os"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

type WatcherLogger struct{}

func NewWatcherLogger() *WatcherLogger {
	return &WatcherLogger{}
}

func (l *WatcherLogger) Info(text string) {
	fmt.Printf(
		"%s[Watcher]%s %s%s%s\n",
		colorYellow,
		colorReset,
		colorCyan,
		text,
		colorReset,
	)
}

func (l *WatcherLogger) Error(text string) {
	fmt.Printf(
		"%s[Watcher]%s %s%s%s\n",
		colorYellow,
		colorReset,
		colorRed,
		text,
		colorReset,
	)
}

func (l *WatcherLogger) Fatal(text string) {
	l.Error(text)
	os.Exit(1)
}
