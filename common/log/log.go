package log

import (
	"flag"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

const logFormat = `date=%s, method=%s, url=%s,  response_time=%s`

var logLevel = flag.String("log_level", "info", "set log level")
var errorLogPath = flag.String("error_log", "logs/error.log", "log path")

var version = flag.Bool("v", false, "for version")

func init() {
	if *version {
		os.Exit(0)
		return
	}
	SetLevel(*logLevel)
	if f := reopen(*errorLogPath); f != nil {
		logrus.SetOutput(f)
	}
}

func reopen(filename string) *os.File {
	if filename == "" {
		return nil
	}

	logFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		// logrus.Error("Error in opening ", filename, err)
		return nil
	}
	return logFile
}

type Fields logrus.Fields

func SetLevel(level string) {
	switch level {
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "warning":
		logrus.SetLevel(logrus.WarnLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func GetLevel() string {
	return strings.ToUpper(logrus.GetLevel().String())
}

func Request(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		t := time.Now()
		next(w, r, ps)
		Infof(logFormat, t, r.Method, r.RequestURI, time.Since(t))
	}
}

func Info(args ...interface{}) {
	var function string
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
		function = runtime.FuncForPC(pc).Name()[strings.LastIndex(runtime.FuncForPC(pc).Name(), "/")+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d %s", file, line, function)).Info(args...)
}

func Infoln(args ...interface{}) {
	var function string
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
		function = runtime.FuncForPC(pc).Name()[strings.LastIndex(runtime.FuncForPC(pc).Name(), "/")+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d %s", file, line, function)).Infoln(args...)
}

func Infof(format string, args ...interface{}) {
	var function string
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
		function = runtime.FuncForPC(pc).Name()[strings.LastIndex(runtime.FuncForPC(pc).Name(), "/")+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d %s", file, line, function)).Infof(format, args...)
}

func Print(args ...interface{}) {
	var function string
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
		function = runtime.FuncForPC(pc).Name()[strings.LastIndex(runtime.FuncForPC(pc).Name(), "/")+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d %s", file, line, function)).Info(args...)
}

func Println(args ...interface{}) {
	var function string
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
		function = runtime.FuncForPC(pc).Name()[strings.LastIndex(runtime.FuncForPC(pc).Name(), "/")+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d %s", file, line, function)).Infoln(args...)
}

func Printf(format string, args ...interface{}) {
	var function string
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
		function = runtime.FuncForPC(pc).Name()[strings.LastIndex(runtime.FuncForPC(pc).Name(), "/")+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d %s", file, line, function)).Infof(format, args...)
}

func Debug(args ...interface{}) {
	var function string
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
		function = runtime.FuncForPC(pc).Name()[strings.LastIndex(runtime.FuncForPC(pc).Name(), "/")+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d %s", file, line, function)).Debug(args...)
}

func Debugln(args ...interface{}) {
	var function string
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
		function = runtime.FuncForPC(pc).Name()[strings.LastIndex(runtime.FuncForPC(pc).Name(), "/")+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d %s", file, line, function)).Debugln(args...)
}

func Debugf(format string, args ...interface{}) {
	var function string
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
		function = runtime.FuncForPC(pc).Name()[strings.LastIndex(runtime.FuncForPC(pc).Name(), "/")+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d %s", file, line, function)).Debugf(format, args...)
}

func Warn(args ...interface{}) {
	var function string
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
		function = runtime.FuncForPC(pc).Name()[strings.LastIndex(runtime.FuncForPC(pc).Name(), "/")+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d %s", file, line, function)).Warn(args...)
}

func Warnln(args ...interface{}) {
	var function string
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
		function = runtime.FuncForPC(pc).Name()[strings.LastIndex(runtime.FuncForPC(pc).Name(), "/")+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d %s", file, line, function)).Warnln(args...)
}

func Warnf(format string, args ...interface{}) {
	var function string
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
		function = runtime.FuncForPC(pc).Name()[strings.LastIndex(runtime.FuncForPC(pc).Name(), "/")+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d %s", file, line, function)).Warnf(format, args...)
}

func Error(args ...interface{}) {
	var function string
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
		function = runtime.FuncForPC(pc).Name()[strings.LastIndex(runtime.FuncForPC(pc).Name(), "/")+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d %s", file, line, function)).Error(args...)
}

func Errorln(args ...interface{}) {
	var function string
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
		function = runtime.FuncForPC(pc).Name()[strings.LastIndex(runtime.FuncForPC(pc).Name(), "/")+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d %s", file, line, function)).Errorln(args...)
}

func Errorf(format string, args ...interface{}) {
	var function string
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
		function = runtime.FuncForPC(pc).Name()[strings.LastIndex(runtime.FuncForPC(pc).Name(), "/")+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d %s", file, line, function)).Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	var function string
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
		function = runtime.FuncForPC(pc).Name()[strings.LastIndex(runtime.FuncForPC(pc).Name(), "/")+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d %s", file, line, function)).Fatal(args...)
}

func Fatalln(args ...interface{}) {
	var function string
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
		function = runtime.FuncForPC(pc).Name()[strings.LastIndex(runtime.FuncForPC(pc).Name(), "/")+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d %s", file, line, function)).Fatalln(args...)
}

func Fatalf(format string, args ...interface{}) {
	var function string
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
		function = runtime.FuncForPC(pc).Name()[strings.LastIndex(runtime.FuncForPC(pc).Name(), "/")+1:]
	}
	logrus.WithField("source", fmt.Sprintf("%s:%d %s", file, line, function)).Fatalf(format, args...)
}

func WithFields(fields Fields) *logrus.Entry {
	var function string
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		file = file[slash+1:]
		function = runtime.FuncForPC(pc).Name()[strings.LastIndex(runtime.FuncForPC(pc).Name(), "/")+1:]
	}

	fields["source"] = fmt.Sprintf("%s:%d %s", file, line, function)

	logrusFields := logrus.Fields{}

	for key, value := range fields {
		logrusFields[key] = value
	}

	return logrus.WithFields(logrusFields)
}
