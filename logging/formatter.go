package logging

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

type (
	// Formatter is a custom logrus formatter to provide bunyan-style logs
	Formatter struct {
		Name string
	}
)

var (
	pid          int
	hostname     string
	bunyanLevels = map[string]int{
		"debug":   20,
		"info":    30,
		"warning": 40,
		"error":   50,
		"fatal":   60,
	}
)

func init() {
	// Initialize vars
	pid = os.Getpid()

	var err error
	hostname, err = os.Hostname()
	if err != nil {
		hostname = "<hostname n/a>"
	}
}

// NewFormatter delivers a formatter
func NewFormatter(name string) *Formatter {
	return &Formatter{name}
}

// Format implements the formatter interface
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Build data list
	data := logrus.Fields{}
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}

	data["name"] = f.Name
	data["hostname"] = hostname
	data["pid"] = pid
	level, ok := bunyanLevels[entry.Level.String()]
	if ok {
		data["level"] = level
	} else {
		data["level"] = entry.Level.String()
	}

	data["msg"] = entry.Message
	data["time"] = entry.Time.Format("2006-01-02T15:04:05.999Z07:00")
	data["v"] = 0

	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal log fields: %v", err)
	}

	return append(serialized, '\n'), nil
}
