// ezmq: An easy golang amqp client.
// Copyright (C) 2022  super9du
//
// This library is free software; you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 2.1 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library; If not, see <https://www.gnu.org/licenses/>.

package ezmq

var (
	_default Logger = &printLogger{}
)

func SetLogger(lg Logger) {
	_default = lg
}

// erro will print only when printLogger's Level is erro, warn, info or debug.
func erro(v ...interface{}) {
	_default.Error(v...)
}

// errorf will print only when printLogger's Level is erro, warn, info or debug.
func errorf(format string, args ...interface{}) {
	_default.Errorf(format, args...)
}

// warn will print when printLogger's Level is warn, info or debug.
func warn(v ...interface{}) {
	_default.Warn(v...)
}

// warnf will print when printLogger's Level is warn, info or debug.
func warnf(format string, args ...interface{}) {
	_default.Warnf(format, args...)
}

// info will print when printLogger's Level is info or debug.
func info(v ...interface{}) {
	_default.Info(v...)
}

// infof will print when printLogger's Level is info or debug.
func infof(format string, args ...interface{}) {
	_default.Infof(format, args...)
}

// debug will print when printLogger's Level is debug.
func debug(v ...interface{}) {
	_default.Debug(v...)
}

// debugf will print when printLogger's Level is debug.
func debugf(format string, args ...interface{}) {
	_default.Debugf(format, args...)
}
