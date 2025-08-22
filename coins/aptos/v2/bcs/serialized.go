package bcs

// Serialized represents a serialized transaction argument
type Serialized struct {
	Value []byte
}

// NewSerialized creates a new Serialized instance
func NewSerialized(value []byte) *Serialized {
	return &Serialized{
		Value: value,
	}
}

// Serialize serializes the Serialized instance
func (s *Serialized) Serialized(serializer *Serializer) {
	serializer.WriteBytes(s.Value)
}

// SerializeForEntryFunction serializes the Serialized instance for entry function
func (s *Serialized) SerializedForEntryFunction(serializer *Serializer) {
	s.Serialized(serializer)
}

// SerializeForScriptFunction serializes the Serialized instance for script function
func (s *Serialized) SerializedForScriptFunction(serializer *Serializer) {
	serializer.Uleb128(uint32(9))
	s.Serialized(serializer)
}
