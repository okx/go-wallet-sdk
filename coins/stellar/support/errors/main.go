/*
 * Copyright 2016 Stellar Development Foundation and contributors.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * 	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * This file includes portions of third-party code from [https://github.com/stellar/go].
 * The original code is licensed under the Apache License 2.0.
 */

package errors

import (
	"github.com/pkg/errors"
)

// StackTracer represents a type (usually an error) that can provide a stack
// trace.
type StackTracer interface {
	StackTrace() errors.StackTrace
}

// Cause returns the underlying cause of the error, if possible.  See
// https://godoc.org/github.com/pkg/errors#Cause for further details.
func Cause(err error) error {
	return errors.Cause(err)
}

// Errorf formats according to a format specifier and returns the string as a
// value that satisfies error. See
// https://godoc.org/github.com/pkg/errors#Errorf for further details
func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

// New returns an error with the supplied message. See
// https://godoc.org/github.com/pkg/errors#New for further details
func New(message string) error {
	return errors.New(message)
}

// Wrap returns an error annotating err with message. If err is nil, Wrap
// returns nil.  See https://godoc.org/github.com/pkg/errors#Wrap for more
// details.
func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

// Wrapf returns an error annotating err with the format specifier. If err is
// nil, Wrapf returns nil. See https://godoc.org/github.com/pkg/errors#Wrapf
// for more details.
func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}
