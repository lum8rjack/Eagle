package modules

import (
	"fmt"
	"os"
	"time"
)

// Struct that contains the details for logging
type Logger struct {
	StartTime      time.Time
	EndTime        time.Time
	SaveOutput     bool
	OutputFilename string
}

func FileExists(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		return false
	}

	return true
}

// Write data to a file
func (l Logger) WriteToFile(data string) {
	fmt.Print(data)
	if l.SaveOutput {
		// Open file to append to it
		file, err := os.OpenFile(l.OutputFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			// Write to file
			file.WriteString(data)
		}
		defer file.Close()
	}
}

func timeToString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// Specify when the scan is starting
func (l Logger) Start() {
	l.StartTime = time.Now()
	stime := timeToString(l.StartTime)
	result := fmt.Sprintf("Eagle started scan at %s\n", stime)
	if l.SaveOutput {
		fd, _ := os.Create(l.OutputFilename)
		defer fd.Close()
	}
	l.WriteToFile(result)
}

// Specify when the scan is completed
func (l Logger) Stop() {
	l.EndTime = time.Now()
	dif := l.EndTime.Sub(l.StartTime).Seconds()
	data := fmt.Sprintf("Completed scan in %.2f seconds\n", dif)
	l.WriteToFile(data)
	if l.SaveOutput {
		fmt.Printf("Output file saved: %s\n", l.OutputFilename)
	}
}
