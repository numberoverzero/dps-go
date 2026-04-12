package blob

import "encoding/json"

// Codec converts between T and raw bytes for storage.
// Construct one with [NewCodec] or use a built-in like [JSONCodec] or [RawCodec].
type Codec[T any] struct {
	marshal   func(T) ([]byte, error)
	unmarshal func([]byte) (T, error)
}

// NewCodec builds a [Codec] from a marshal/unmarshal pair.
// Returns [ErrCodec] if either argument is nil.
func NewCodec[T any](marshal func(T) ([]byte, error), unmarshal func([]byte) (T, error)) (Codec[T], error) {
	if marshal == nil || unmarshal == nil {
		return Codec[T]{}, ErrCodec
	}

	return Codec[T]{marshal, unmarshal}, nil
}

// JSONCodec returns a Codec that marshals T using encoding/json.
func JSONCodec[T any]() Codec[T] {
	return Codec[T]{
		marshal: func(v T) ([]byte, error) {
			return json.Marshal(v)
		},
		unmarshal: func(data []byte) (T, error) {
			var v T
			if err := json.Unmarshal(data, &v); err != nil {
				return v, err
			}

			return v, nil
		},
	}
}

// RawCodec is a Codec for []byte values that passes data through unchanged.
var RawCodec = Codec[[]byte]{
	marshal:   func(data []byte) ([]byte, error) { return data, nil },
	unmarshal: func(data []byte) ([]byte, error) { return data, nil },
}
