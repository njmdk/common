package logger

import (
	"path/filepath"
	"sync/atomic"
	"time"
	"unsafe"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/njmdk/common/network/ulimit"
	"github.com/njmdk/common/timer"
	"github.com/njmdk/common/utils"
)

var emptyStr = ""

var log *Logger

func SetDefaultLog(l *Logger) {
	log = l
}

func GetDefaultLog() *Logger {
	return log
}

func InitDefaultLogger(name, path string, lvl zapcore.Level, isDebug bool) (*Logger, error) {
	var err error
	log, err = New(name, path, lvl, isDebug)

	return log, err
}

func New(name, path string, lvl zapcore.Level, isDebug bool) (*Logger, error) {
	err := ulimit.SetRLimit()
	if err != nil {
		return nil, err
	}
	err = utils.CheckAndCreate(path)
	if err != nil {
		return nil, err
	}

	l := &Logger{
		tomorrow: unsafe.Pointer(&emptyStr),
		name:     name,
		path:     path,
		lvl:      lvl,
		isDebug:  isDebug,
	}
	l.checkTomorrow()
	l.deleteBeforeLog()
	return l, nil
}

func newLogger(path string, lvl zapcore.Level, isDebug bool) (*zap.Logger, error) {
	ec := zap.NewDevelopmentEncoderConfig()

	outPath := []string{"stderr", path}
	if !isDebug {
		outPath = []string{path}
	}

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(lvl),
		Development:      false,
		Encoding:         "console",
		EncoderConfig:    ec,
		OutputPaths:      outPath,
		ErrorOutputPaths: []string{"stderr"},
	}
	
	return cfg.Build(zap.AddStacktrace(zap.ErrorLevel), zap.AddCallerSkip(1))
}

type Logger struct {
	log      unsafe.Pointer
	sugar    unsafe.Pointer
	tomorrow unsafe.Pointer
	name     string
	path     string
	t        time.Time
	lvl      zapcore.Level
	isDebug  bool
}

func (this_ *Logger) SetCheckTomorrowTime(t time.Time) {
	this_.t = t
}

func (this_ *Logger)IsLogDebug()bool  {
	return this_.isDebug
}

func (this_ *Logger) checkTomorrow() {
	t := timer.Now()

	tomorrowStr := t.Format("2006_01_02")
	if tomorrowStr != *(*string)(atomic.LoadPointer(&this_.tomorrow)) {
		if atomic.CompareAndSwapPointer(&this_.tomorrow, this_.tomorrow, unsafe.Pointer(&tomorrowStr)) {
			pathFile := filepath.Join(this_.path, tomorrowStr+"_"+this_.name+".log")

			log, err := newLogger(pathFile, this_.lvl, this_.isDebug)
			if err != nil {
				panic(err)
				return
			}

			l := (*zap.Logger)(atomic.LoadPointer(&this_.log))
			if l != nil {
				_ = l.Sync()
			}
			
			atomic.StorePointer(&this_.log, unsafe.Pointer(log))
			atomic.StorePointer(&this_.sugar, unsafe.Pointer(log.Sugar()))
		}
	}
}

func (this_ *Logger) getLog() *zap.Logger {
	this_.checkTomorrow()
	return (*zap.Logger)(atomic.LoadPointer(&this_.log))
}

func (this_ *Logger) getSugar() *zap.SugaredLogger {
	this_.checkTomorrow()
	return (*zap.SugaredLogger)(atomic.LoadPointer(&this_.sugar))
}

func (this_ *Logger) Debug(msg string, fields ...zap.Field) {
	this_.getLog().Debug(msg, fields...)
}

func (this_ *Logger) Info(msg string, fields ...zap.Field) {
	this_.getLog().Info(msg, fields...)
}

func (this_ *Logger) Warn(msg string, fields ...zap.Field) {
	this_.getLog().Warn(msg, fields...)
}

func (this_ *Logger) Error(msg string, fields ...zap.Field) {
	this_.getLog().Error(msg, fields...)
}

func (this_ *Logger) DPanic(msg string, fields ...zap.Field) {
	this_.getLog().DPanic(msg, fields...)
}

func (this_ *Logger) Panic(msg string, fields ...zap.Field) {
	this_.getLog().Panic(msg, fields...)
}

func (this_ *Logger) Fatal(msg string, fields ...zap.Field) {
	this_.getLog().Fatal(msg, fields...)
}

func (this_ *Logger) DebugFormat(format string, args ...interface{}) {
	this_.getSugar().Debugf(format, args...)
}

func (this_ *Logger) InfoFormat(format string, args ...interface{}) {
	this_.getSugar().Infof(format, args...)
}

func (this_ *Logger) WarnFormat(format string, args ...interface{}) {
	this_.getSugar().Warnf(format, args...)
}

func (this_ *Logger) ErrorFormat(format string, args ...interface{}) {
	this_.getSugar().Errorf(format, args...)
}

func (this_ *Logger) DPanicFormat(format string, args ...interface{}) {
	this_.getSugar().DPanicf(format, args...)
}

func (this_ *Logger) PanicFormat(format string, args ...interface{}) {
	this_.getSugar().Panicf(format, args...)
}

func (this_ *Logger) FatalFormat(format string, args ...interface{}) {
	this_.getSugar().Fatalf(format, args...)
}

func Debug(msg string, fields ...zap.Field) {
	log.getLog().Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	log.getLog().Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	log.getLog().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	log.getLog().Error(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	log.getLog().DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	log.getLog().Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	log.getLog().Fatal(msg, fields...)
}

func DebugFormat(format string, args ...interface{}) {
	log.getSugar().Debugf(format, args...)
}

func InfoFormat(format string, args ...interface{}) {
	log.getSugar().Infof(format, args...)
}

func WarnFormat(format string, args ...interface{}) {
	log.getSugar().Warnf(format, args...)
}

func ErrorFormat(format string, args ...interface{}) {
	log.getSugar().Errorf(format, args...)
}

func DPanicFormat(format string, args ...interface{}) {
	log.getSugar().DPanicf(format, args...)
}

func PanicFormat(format string, args ...interface{}) {
	log.getSugar().Panicf(format, args...)
}

func FatalFormat(format string, args ...interface{}) {
	log.getSugar().Fatalf(format, args...)
}
