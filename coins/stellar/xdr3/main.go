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

package xdr

// Enum indicates this implementing type should be serialized/deserialized
// as an XDR Enum.  Implement ValidEnum to specify what values are valid for
// this enum.
type Enum interface {
	ValidEnum(int32) bool
}

// Sized types are types that have an explicit maximum size.  By default, the
// variable length XDR types (VarArray, VarOpaque and String) have a maximum
// byte size of a 2^32-1, but an implementor of this type may reduce that
// maximum to an appropriate value for the XDR schema in use.
type Sized interface {
	XDRMaxSize() int
}

// Union indicates the implementing type should be serialized/deserialized as
// an XDR Union.  The implementer must provide public fields, one for the
// union's disciminant, whose name must be returned by ArmForSwitch(), and
// one per potential value of the union, which must be a pointer.  For example:
//
//	type Result struct {
//	  Type ResultType  // this is the union's disciminant, may be 0 to indicate success, 1 to indicate error
//	  Msg  *string // this field will be populated when Type == 1
//	}
type Union interface {
	ArmForSwitch(int32) (string, bool)
	SwitchFieldName() string
}
