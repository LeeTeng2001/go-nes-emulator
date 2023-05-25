package mlog

import (
	"github.com/charmbracelet/log"
	"os"
	"time"
)

var L = log.NewWithOptions(os.Stdout, log.Options{
	//ReportCaller:    true,
	ReportTimestamp: true,
	TimeFormat:      time.DateTime,
})

func init() {
	L.Info("Starting application...")
}
