/*******************************************************************************
Stuff I copy from project to project to handle logging

Created by thomas.a.cherry@nasa.gov
Created: January 2026
*******************************************************************************/

package lib

import (
	"io"
	"log"
	"os"
)

/******************************************************************************/
// #MARK: Variables and Structs

type Logger struct {
	Report  *log.Logger
	Error   *log.Logger
	Warning *log.Logger
	Warn    *log.Logger
	Info    *log.Logger
	Debug   *log.Logger
}

var Log Logger

type LogLevel int

const (
	LogLevelReport LogLevel = iota
	LogLevelError
	LogLevelWarning
	LogLevelInfo
	LogLevelDebug
)

/******************************************************************************/
// #MARK: Logging functions

func init() {
	file := os.Stderr

	Log.Report = log.New(file, "REPORT: ", log.Ldate|log.Ltime|log.Lshortfile)
	Log.Error = log.New(file, "❌ ERROR: ", log.Ldate|log.Ltime|log.Llongfile)
	Log.Warn = log.New(file, "⚠️ WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Log.Info = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Log.Debug = log.New(file, "🚀 DEBUG : ", log.Ldate|log.Ltime|log.Llongfile)

	Log.Info.SetOutput(io.Discard)
	Log.Debug.SetOutput(io.Discard)
}

func EnableInfo() {
	Log.Info.SetOutput(os.Stderr)
}

func EnableDebug() {
	Log.Debug.SetOutput(os.Stderr)
}

func SetLogLevel(level LogLevel) {
	switch level {
	case LogLevelReport:
		Log.Report.SetOutput(os.Stderr)
		//off
		Log.Error.SetOutput(io.Discard)
		Log.Warn.SetOutput(io.Discard)
		Log.Info.SetOutput(io.Discard)
		Log.Debug.SetOutput(io.Discard)
	case LogLevelError:
		Log.Report.SetOutput(os.Stderr)
		Log.Error.SetOutput(os.Stderr)
		//off
		Log.Warn.SetOutput(io.Discard)
		Log.Info.SetOutput(io.Discard)
		Log.Debug.SetOutput(io.Discard)
	case LogLevelWarning:
		Log.Report.SetOutput(os.Stderr)
		Log.Error.SetOutput(os.Stderr)
		Log.Warn.SetOutput(os.Stderr)
		//off
		Log.Info.SetOutput(io.Discard)
		Log.Debug.SetOutput(io.Discard)
	case LogLevelInfo:
		Log.Report.SetOutput(os.Stderr)
		Log.Error.SetOutput(os.Stderr)
		Log.Warn.SetOutput(os.Stderr)
		Log.Info.SetOutput(os.Stderr)
		//off
		Log.Debug.SetOutput(io.Discard)
	case LogLevelDebug:
		Log.Report.SetOutput(os.Stderr)
		Log.Error.SetOutput(os.Stderr)
		Log.Warn.SetOutput(os.Stderr)
		Log.Info.SetOutput(os.Stderr)
		Log.Debug.SetOutput(os.Stderr)
	}
}
