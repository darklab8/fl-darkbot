package typelog

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type Logger struct {
	logger              *slog.Logger
	name                string
	enable_file_showing bool
	enable_json_format  bool
	log_level_slog      *slog.LevelVar
	level_log           LogLevel
}

type LoggerParam func(r *Logger)

func WithJsonFormat(state bool) LoggerParam {
	return func(logger *Logger) {
		logger.enable_json_format = state
	}
}

func WithFileShowing(state bool) LoggerParam {
	return func(logger *Logger) {
		logger.enable_file_showing = state
	}
}

func WithLogLevelStr(log_level_str string) LoggerParam {
	return WithLogLevel(LogLevel(log_level_str))
}

func WithLogLevel(log_level_str LogLevel) LoggerParam {
	return func(logger *Logger) {
		var log_level LogLevel = LogLevel(log_level_str)

		var programLevel = new(slog.LevelVar) // Info by default
		switch log_level {
		case LEVEL_DEBUG:
			programLevel.Set(slog.LevelDebug)
		case LEVEL_INFO:
			programLevel.Set(slog.LevelInfo)
		case LEVEL_WARN:
			programLevel.Set(slog.LevelWarn)
		case LEVEL_ERROR:
			programLevel.Set(slog.LevelError)
		case LEVEL_DEFAULT_WARN:
			programLevel.Set(slog.LevelWarn)
		default:
			panic(fmt.Sprintf("invalid log level=%s for logger=%s", log_level_str, logger.name))
		}
		logger.log_level_slog = programLevel
		logger.level_log = log_level_str
	}
}

/*
Leaving option for end applications to access all registered loggers
and choosing their own log levels to override for them
*/
var RegisteredLoggers []*Logger

func (l *Logger) GetName() string { return l.name }

func NewLogger(
	name string,
	options ...LoggerParam,
) *Logger {

	logger := &Logger{
		name: name,
	}
	RegisteredLoggers = append(RegisteredLoggers, logger)

	WithJsonFormat(bool(EnvTurnJSON))(logger)
	WithFileShowing(EnvTurnFileShowing)(logger)
	WithLogLevelStr(os.Getenv(strings.ToUpper(name) + "_LOG_LEVEL"))(logger)

	for _, opt := range options {
		opt(logger)
	}

	return logger.Initialized()
}

/*
For overrides by external libraries
*/
func (logger *Logger) OverrideOption(options ...LoggerParam) *Logger {
	for _, opt := range options {
		opt(logger)
	}

	return logger.Initialized()
}

func (logger *Logger) Initialized() *Logger {
	if logger.enable_json_format {
		logger.logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logger.log_level_slog}))
	} else {
		logger.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logger.log_level_slog}))
	}
	return logger
}

func (l *Logger) WithFields(opts ...LogType) *Logger {
	var new_logger Logger = *l
	new_logger.Initialized()
	new_logger.logger = new_logger.logger.With(newSlogArgs(opts...)...)
	return &new_logger
}
