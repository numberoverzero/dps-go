package blob

// TypedReader provides read access to a partition, decoding raw bytes into T.
type TypedReader[T any] interface {
	PartitionKey() string
	Get(key string) (T, error)
	Has(key string) (bool, error)
	List(prefix string) ([]string, error)
	Read(func(TypedReader[T]) error) error
}

// TypedWriter provides read-write access to a partition, encoding T to raw bytes.
type TypedWriter[T any] interface {
	TypedReader[T]
	Set(key string, value T) error
	Remove(key string) (bool, error)
	Write(func(TypedWriter[T]) error) error
}

// NewTypedReader wraps a [PartitionReader] with a [Codec] to produce a [TypedReader].
func NewTypedReader[T any](r PartitionReader, c Codec[T]) TypedReader[T] {
	return &reader[T]{r, c}
}

// NewTypedWriter wraps a [PartitionWriter] with a [Codec] to produce a [TypedWriter].
func NewTypedWriter[T any](w PartitionWriter, c Codec[T]) TypedWriter[T] {
	return &writer[T]{reader[T]{w, c}, w}
}

// -- reader ------------------------------------------------------------------

type reader[T any] struct {
	r     PartitionReader
	codec Codec[T]
}

func (x *reader[T]) PartitionKey() string { return x.r.PartitionKey() }
func (x *reader[T]) Get(k string) (T, error) {
	b, err := x.r.Get(k)
	if err != nil {
		var zero T

		return zero, err
	}

	return x.codec.unmarshal(b)
}
func (x *reader[T]) Has(k string) (bool, error)      { return x.r.Has(k) }
func (x *reader[T]) List(p string) ([]string, error) { return x.r.List(p) }
func (x *reader[T]) Read(fn func(TypedReader[T]) error) error {
	return x.r.Read(func(inner PartitionReader) error {
		return fn(&reader[T]{inner, x.codec})
	})
}

// -- writer ------------------------------------------------------------------

type writer[T any] struct {
	reader[T]
	w PartitionWriter
}

func (x *writer[T]) Set(k string, v T) error {
	data, err := x.codec.marshal(v)
	if err != nil {
		return err
	}

	return x.w.Set(k, data)
}
func (x *writer[T]) Remove(k string) (bool, error) { return x.w.Remove(k) }
func (x *writer[T]) Write(fn func(TypedWriter[T]) error) error {
	return x.w.Write(func(inner PartitionWriter) error {
		return fn(&writer[T]{reader[T]{inner, x.codec}, inner})
	})
}
