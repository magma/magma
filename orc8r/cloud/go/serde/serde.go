package serde

// Serde (SERializer-DEserializer) implements logic to serialize/deserialize
// a specific piece of data.
type Serde interface {
	// GetDomain returns a globally unique key which represents the domain of
	// this Serde. Serde types are unique within each domain but not across
	// domains.
	GetDomain() string

	// GetType returns a unique key within the domain for the specific Serde
	// implementation. This represents the type of data that the Serde will be
	// responsible for serializing and deserialing.
	GetType() string

	// Serialize a piece of data
	Serialize(in interface{}) ([]byte, error)

	// Deserialize a piece of data
	Deserialize(in []byte) (interface{}, error)
}
