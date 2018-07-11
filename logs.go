// MIT License
//
// Copyright (c) 2017 Lalamove.com
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
//
// Reference: https://lalamove.atlassian.net/wiki/spaces/TECH/pages/82149406/Kubernetes
// Lalamove kubernetes logging format
//{
//    "message": "", // string describing what happened
//    "src_file": "", // file path
//    "src_line": "", // line number
//    "context": {}, // custom field here
//    "level": "", // debug/info/warning/error/fatal
//    "time": "", // ISO8601.nanoseconds+TZ (in node only support precision up to milliseconds)
//    "backtrace": "" // err stack
//}

package logs

import (
	"runtime"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	ISO8601 = "2006-01-02T15:04:05.000000000Z0700"

	TimeKey        = "time"
	LevelKey       = "level"
	CallerKey      = "src_file"
	MessageKey     = "message"
	StacktraceKey  = "backtrace"
	CustomFieldKey = "context"

	SourceLineKey = "src_line"

	Warning = "warning"
)

var (
	Log    *zap.Logger
	Config *zap.Config
	once   sync.Once
)

// InitLogger will build the logger with default config
// It is thread safe as it creates instance once only
// In order to build a customized logger,
// developer must run InitLogger before Logger() function
func InitLogger() {
	once.Do(func() {
		cfg := *NewConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		cfg.EncoderConfig = *NewLalamoveEncoderConfig()
		Log, _ = cfg.Build()
	})
}

// NewConfig will return a zap config
func NewConfig() *zap.Config {
	once.Do(func() {
		c := zap.NewProductionConfig()
		Config = &c
	})
	return Config
}

// NewLalamoveEncoderConfig will create an EncoderConfig
func NewLalamoveEncoderConfig() *zapcore.EncoderConfig {
	return &zapcore.EncoderConfig{
		TimeKey:        TimeKey,
		LevelKey:       LevelKey,
		CallerKey:      CallerKey,
		MessageKey:     MessageKey,
		StacktraceKey:  StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    LalamoveLevelEncoder,
		EncodeTime:     LalamoveISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   LalamoveCallerEncoder,
	}
}

// LalamoveLevelEncoder will convert the warn display string to warning
func LalamoveLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	if l == zapcore.WarnLevel {
		// Set warn label to warning
		enc.AppendString(Warning)
	} else {
		enc.AppendString(l.String())
	}
}

func LalamoveCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(caller.File)
}

// LalamoveISO8601TimeEncoder will convert the time to ISO8601 based on Lalamove k8s logging format
func LalamoveISO8601TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.UTC().Format(ISO8601))
}

// Logger will create a zap based logger
// Extra field will inside fields namespace
// return a *zap.Logger for logging
func Logger() *zap.Logger {
	InitLogger()
	// Skip this function by one
	// ln int is line number of source file
	_, _, ln, _ := runtime.Caller(1)

	return Log.With(zap.String(SourceLineKey, strconv.FormatInt(int64(ln), 10)), zap.Namespace(CustomFieldKey))
}
