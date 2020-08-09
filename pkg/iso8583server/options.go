package iso8583server

import (
	"net"
	"time"
)

type Option func(*configuration)

func OptionListener(l net.Listener) Option {
	return func(config *configuration) {
		config.Listener = l
	}
}

func OptionTimeout(t time.Duration) Option {
	return func(config *configuration) {
		config.Timeout = t
	}
}

func OptionReader(reader NetReadFunc) Option {
	return func(config *configuration) {
		config.NetRead = reader
	}
}

func OptionInfoLogger(logger LogFunc) Option {
	return func(config *configuration) {
		if logger == nil {
			config.LogInfo = nopLogger
		}
		config.LogInfo = logger
	}
}

func OptionErrLogger(logger LogFunc) Option {
	return func(config *configuration) {
		if logger == nil {
			config.LogErr = nopLogger
		}
		config.LogErr = logger
	}
}

func OptionMTIRead(reader ReadMTIFunc) Option {
	return func(config *configuration) {
		config.ReadMTI = reader
	}
}

func OptionUnknownHandler(handler HandlerFunc) Option {
	return func(config *configuration) {
		config.UnknownHandler = handler
	}
}

var nopLogger = func(v ...interface{}) {}
