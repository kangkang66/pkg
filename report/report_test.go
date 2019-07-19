package report

import (
	"testing"
	"time"
)

func TestNewErrorLogReport(t *testing.T) {
	config := &ErrorLogReportConfig{
		SecondMaxNumToReport:20,
		MinuteMaxNumToReport:20,
		HourMaxNumToReport:20,
		MinTimeToReport:10,
	}

	e := NewErrorLogReport(config)
	go e.Reset()
	for {
		e.WriteMinute("kang")
		time.Sleep(50 * time.Millisecond)
	}
}
