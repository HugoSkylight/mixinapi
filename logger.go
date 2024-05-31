package mixinapi

import (
	"context"
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, err error, args ...any)
	Error(ctx context.Context, msg string, err error, args ...any)
	Fatal(ctx context.Context, msg string, err error, args ...any)
}

var (
	logger      Logger
	sp          = string(filepath.Separator)
	logLevelMap = map[int]zapcore.Level{
		6: zapcore.DebugLevel,
		5: zapcore.DebugLevel,
		4: zapcore.InfoLevel,
		3: zapcore.WarnLevel,
		2: zapcore.ErrorLevel,
		1: zapcore.FatalLevel,
		0: zapcore.PanicLevel,
	}
)

type ZapLogger struct {
	logger       *zap.SugaredLogger
	prefix       string
	level        zapcore.Level
	rotationTime time.Duration
}

func InitLogger(
	prefix string,
	logLevel int,
) error {
	l, err := NewZapLogger(prefix, true, logLevel, false, ".logs", 7, 24)
	if err != nil {
		return err
	}
	logger = l
	return nil
}

func Debug(msg string, args ...any) {
	logger.Debug(context.Background(), msg, args)
}

func Info(msg string, args ...any) {
	logger.Info(context.Background(), msg, args)
}

func Warn(msg string, err error, args ...any) {
	logger.Warn(context.Background(), msg, err, args)
}

func Error(msg string, err error, args ...any) {
	logger.Error(context.Background(), msg, err, args)
}

func Fatal(msg string, err error, args ...any) {
	logger.Fatal(context.Background(), msg, err, args)
}

func NewZapLogger(
	prefix string,
	isStdout bool,
	logLevel int,
	isJson bool,
	logLocation string,
	rotateCount uint,
	rotationTime uint,
) (*ZapLogger, error) {
	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(logLevelMap[logLevel]),
		DisableStacktrace: true,
	}
	if isJson {
		zapConfig.Encoding = "json"
	} else {
		zapConfig.Encoding = "console"
	}
	zl := &ZapLogger{prefix: prefix, rotationTime: time.Duration(rotationTime) * time.Hour}
	opts, err := zl.cores(isStdout, logLocation, rotateCount)
	if err != nil {
		return nil, err
	}
	l, err := zapConfig.Build(opts)
	if err != nil {
		return nil, err
	}
	zl.logger = l.Sugar()
	return zl, nil
}

func (l *ZapLogger) getWriter(logLocation string, rorateCount uint) (zapcore.WriteSyncer, error) {
	logf, err := rotatelogs.New(logLocation+sp+l.prefix+".%Y-%m-%d",
		rotatelogs.WithRotationCount(rorateCount),
		rotatelogs.WithRotationTime(l.rotationTime),
	)
	if err != nil {
		return nil, err
	}
	return zapcore.AddSync(logf), nil
}

func (l *ZapLogger) timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	layout := "2006-01-02 15:04:05.000"
	type appendTimeEncoder interface {
		AppendTimeLayout(time.Time, string)
	}
	if enc, ok := enc.(appendTimeEncoder); ok {
		enc.AppendTimeLayout(t, layout)
		return
	}
	enc.AppendString(t.Format(layout))
}

func (l *ZapLogger) cores(isStdout bool, logLocation string, rotateCount uint) (zap.Option, error) {
	c := zap.NewProductionEncoderConfig()
	c.EncodeDuration = zapcore.SecondsDurationEncoder
	c.EncodeTime = l.timeEncoder
	c.MessageKey = "msg"
	c.LevelKey = "level"
	c.TimeKey = "time"
	c.CallerKey = "caller"
	c.NameKey = "logger"

	fileEncoder := zapcore.NewConsoleEncoder(c)
	writer, err := l.getWriter(logLocation, rotateCount)
	if err != nil {
		return nil, err
	}

	var cores []zapcore.Core
	if logLocation != "" {
		cores = []zapcore.Core{
			zapcore.NewCore(fileEncoder, writer, zap.NewAtomicLevelAt(l.level)),
		}
	}
	if isStdout {
		cores = append(cores, zapcore.NewCore(fileEncoder, zapcore.Lock(os.Stdout), zap.NewAtomicLevelAt(l.level)))
	}
	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	}), nil
}

func (l *ZapLogger) Debug(ctx context.Context, msg string, args ...any) {
	args = l.toArgs(ctx, args)
	l.logger.Info(msg, zap.Any("args", args))
}

func (l *ZapLogger) Info(ctx context.Context, msg string, args ...any) {
	args = l.toArgs(ctx, args)
	l.logger.Info(msg, zap.Any("args", args))
}

func (l *ZapLogger) Warn(ctx context.Context, msg string, err error, args ...any) {
	if err != nil {
		args = append(args, "error", err.Error())
	}
	args = l.toArgs(ctx, args)
	l.logger.Warn(msg, zap.Any("args", args))
}

func (l *ZapLogger) Error(ctx context.Context, msg string, err error, args ...any) {
	if err != nil {
		args = append(args, "error", err.Error())
	}
	args = l.toArgs(ctx, args)
	l.logger.Error(msg, zap.Any("args", args))
}

func (l *ZapLogger) Fatal(ctx context.Context, msg string, err error, args ...any) {
	if err != nil {
		args = append(args, "error", err.Error())
	}
	args = l.toArgs(ctx, args)
	l.logger.Fatal(msg, zap.Any("args", args))
}

func (l *ZapLogger) toArgs(ctx context.Context, args []any) []any {
	if ctx == nil {
		return args
	}
	return args
}
