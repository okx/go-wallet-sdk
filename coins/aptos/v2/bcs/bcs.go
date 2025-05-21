package bcs

// Marshaler is an interface for any type that can be serialized into BCS
type Marshaler interface {
	MarshalBCS(*Serializer)
}

// Unmarshaler is an interface for any type that can be deserialized from BCS
type Unmarshaler interface {
	UnmarshalBCS(*Deserializer)
}

// Struct is an interface for an on-chain type.  It must be able to be both Marshaler and Unmarshaler for BCS
type Struct interface {
	Marshaler
	Unmarshaler
}

// reverse is a helper function for serialization / deserialization
func reverse(ub []byte) {
	lo := 0
	hi := len(ub) - 1
	for hi > lo {
		t := ub[lo]
		ub[lo] = ub[hi]
		ub[hi] = t
		lo++
		hi--
	}
}
