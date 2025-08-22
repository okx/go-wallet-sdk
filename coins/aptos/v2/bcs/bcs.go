package bcs

// Marshaler is an interface for any type that can be serialized into BCS
//
// It's highly suggested to implement on a pointer to a type, and not the type directly.
// For example, you could implement a simple Marshaler for a given struct.
//
//	type MyStruct struct {
//	  num     uint8
//	  boolean bool
//	}
//
//	func (str *MyStruct) MarshalBCS(ser *Serializer) {
//	  ser.U8(str.num)
//	  ser.Bool(str.boolean)
//	}
//
// Additionally, if there is expected data, you can add errors to serialization.  It's suggested to stop serialization after any errors.
//
//	type MyStruct struct {
//	  num     uint8 // Only allowed to be 0-10
//	  boolean bool
//	}
//
//	func (str *MyStruct) MarshalBCS(ser *Serializer) {
//	  if str.num > 10 {
//	    ser.SetError(fmt.Error("Cannot serialize MyStruct, num is greater than 10: %d", str.num)
//	    return
//	  }
//	  ser.U8(str.num)
//	  ser.Bool(str.boolean)
//	}
type Marshaler interface {
	// MarshalBCS implements a way to serialize the type into BCS.  Note that the error will need to be directly set
	// using [Serializer.SetError] on the [Serializer].  If using this function, you will need to use [Serializer.Error]
	// to retrieve the error.
	MarshalBCS(ser *Serializer)
}

// Unmarshaler is an interface for any type that can be deserialized from BCS
//
// It's highly suggested to implement on a pointer to a type, and not the type directly.
// For example, you could implement a simple Unmarshaler for a given struct.  You will need to add any appropriate error handling.
//
//	type MyStruct struct {
//	  num     uint8
//	  boolean bool
//	}
//
//	func (str *MyStruct) UnmarshalBCS(des *Deserializer) {
//	  str.num = des.U8()
//	  str.boolean = des.Bool()
//	}
//
// Additionally, if there is expected formatting errors, you can add errors to deserialization.  It's suggested to stop serialization after any errors.
//
//	type MyStruct struct {
//	  num     uint8 // Only allowed to be 0-10
//	  boolean bool
//	}
//
//	func (str *MyStruct) UnmarshalBCS(des *Deserializer) {
//	  str.num = des.U8()
//	  if des.Error() {
//	    // End early, since deserialization failed
//	    return
//	  }
//	  if str.num > 10 {
//	    ser.SetError(fmt.Error("Cannot deserialize MyStruct, num is greater than 10: %d", str.num)
//	    return
//	  }
//	  str.boolean = des.Bool()
//	}
type Unmarshaler interface {
	// UnmarshalBCS implements a way to deserialize the type into BCS.  Note that the error will need to be directly set
	// using [Deserializer.SetError] on the [Deserializer].  If using this function, you will need to use [Deserializer.Error]
	// to retrieve the error.
	UnmarshalBCS(des *Deserializer)
}

// Struct is an interface for an on-chain type.  It must be able to be both Marshaler and Unmarshaler for BCS
type Struct interface {
	Marshaler
	Unmarshaler
}
