package textwire

import (
	"fmt"
	"io"
	"os"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

type WatcherLogger struct {
	out io.Writer
}

func NewWatcherLogger() *WatcherLogger {
	return &WatcherLogger{}
}

func (l *WatcherLogger) Info(text string) {
	l.print(colorCyan, text)
}

func (l *WatcherLogger) Error(text string) {
	l.print(colorRed, text)
}

func (l *WatcherLogger) print(color, text string) {
	out := l.out
	if out == nil {
		out = os.Stdout
	}

	fmt.Fprintf(
		out,
		"%s[Watcher]%s %s%s%s\n",
		colorYellow,
		colorReset,
		color,
		text,
		colorReset,
	)
}

func (l *WatcherLogger) Fatal(text string) {
	l.Error(text)
	os.Exit(1)
}
