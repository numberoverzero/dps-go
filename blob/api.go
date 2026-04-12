package blob

import (
	"errors"
)

var (
	// ErrCodec is returned when a Codec is invalid.
	ErrCodec = errors.New("blob: invalid codec")

	// ErrNotFound is returned when a key does not exist in the partition.
	ErrNotFound = errors.New("blob: not found")
)

// Store manages raw partitions.
type Store interface {
	StoreWriter
	Close() error
}

// StoreWriter provides raw read-write access to the store's partitions.
type StoreWriter interface {
	StoreReader
	OpenWriter(key string) (PartitionWriter, error)
	CreatePartition(key string) error
	CopyPartition(srcKey, dstKey string) error
	DeletePartition(key string) error
	Write(func(StoreWriter) error) error
}

// StoreReader provides raw read access to the store's partitions.
type StoreReader interface {
	OpenReader(key string) (PartitionReader, error)
	HasPartition(key string) (bool, error)
	ListPartitions(prefix string) ([]string, error)
	Read(func(StoreReader) error) error
}

// PartitionWriter provides raw read-write access to a partition's key-value data.
type PartitionWriter interface {
	PartitionReader
	Set(key string, value []byte) error
	Remove(key string) (bool, error)
	Write(func(PartitionWriter) error) error
}

// PartitionReader provides raw read access to a partition's key-value data.
type PartitionReader interface {
	PartitionKey() string
	Get(key string) ([]byte, error)
	Has(key string) (bool, error)
	List(prefix string) ([]string, error)
	Read(func(PartitionReader) error) error
}
