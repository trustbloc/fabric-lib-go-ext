/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
/*
Notice: This file has been modified for TrustBloc Fabric Lib Go EXT usage.
Please review third_party pinning scripts and patches for more details.
*/

package logbridge

import (
	"github.com/trustbloc/fabric-lib-go-ext/pkg/common/logging"
)

// Log levels (from pkg/logging/level.go).
const (
	CRITICAL logging.Level = iota
	ERROR
	WARNING
	INFO
	DEBUG
)

// Logger bridges the lib's logger struct
type Logger struct {
	*logging.Logger
	module string
}

// MustGetLogger bridges calls the lib's NewLogger
func MustGetLogger(module string) *Logger {
	fabModule := "fablibgoext"
	logger := logging.NewLogger(fabModule)
	return &Logger{
		Logger: logger,
		module: fabModule,
	}
}

// Warningf bridges calls to the lib logger's Warnf.
func (l *Logger) Warningf(format string, args ...interface{}) {
	l.Warnf(format, args...)
}

// Warning bridges calls to the lib logger's Warn.
func (l *Logger) Warning(args ...interface{}) {
	l.Warn(args...)
}

// IsEnabledFor bridges calls to the lib logger's IsEnabledFor.
func (l *Logger) IsEnabledFor(level logging.Level) bool {
	return logging.IsEnabledFor(l.module, level)
}
