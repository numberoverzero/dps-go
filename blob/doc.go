// Package blob is typed key-value storage backed by partitions.
//
// A [Store] holds named partitions. Each partition is a flat key-value
// namespace. A [Codec] tells the partition how to convert between your
// type and raw bytes.
//
// # Basic Usage
//
// Open a store, then open a partition for writing and wrap it with a codec:
//
//	import "github.com/numberoverzero/dps-go/blob"
//
//	type Profile struct {
//	    Name  string `json:"name"`
//	    Score int    `json:"score"`
//	}
//
//	store, err := blob.Open("app.db")
//	defer store.Close()
//
//	raw, err := store.OpenWriter("players")
//	players := blob.NewTypedWriter(raw, blob.JSONCodec[Profile]())
//
//	err = players.Set("player:42", Profile{Name: "Ava", Score: 100})
//
//	profile, err := players.Get("player:42")
//
// For read-only access, use [StoreReader.OpenReader] and [NewTypedReader].
//
// # Codecs
//
// A [Codec] converts T to []byte and back. Built-ins: [JSONCodec] for
// structs, [RawCodec] for []byte passthrough. For anything else, build
// one with [NewCodec]:
//
//	func protoCodec[T proto.Message]() (blob.Codec[T], error) {
//	    return blob.NewCodec(
//	        func(v T) ([]byte, error) { return proto.Marshal(v) },
//	        func(b []byte) (T, error) {
//	            var v T
//	            return v, proto.Unmarshal(b, v)
//	        },
//	    )
//	}
//
// [NewCodec] returns [ErrCodec] if either function is nil.
//
// # Transactions
//
// Each call on a partition handle auto-wraps itself in a transaction.
// For multi-step atomicity, use [TypedReader.Read] (read-only) or
// [TypedWriter.Write] (read-write):
//
//	err = players.Write(func(tx blob.TypedWriter[Profile]) error {
//	    src, err := tx.Get("player:1")
//	    if err != nil { return err }
//	    dst, err := tx.Get("player:2")
//	    if err != nil { return err }
//
//	    src.Score -= 10
//	    dst.Score += 10
//
//	    if err := tx.Set("player:1", src); err != nil { return err }
//	    return tx.Set("player:2", dst)
//	})
//
// [Store] itself exposes [StoreReader.Read] and [StoreWriter.Write] for
// transactions that span partition management ([StoreWriter.CreatePartition],
// [StoreWriter.CopyPartition], [StoreWriter.DeletePartition]) alongside
// data changes.
//
// # Separating Logic from Transactions
//
// Define operations as functions that accept [TypedReader] (read) or
// [TypedWriter] (read-write). A partition handle and a transaction
// handle both satisfy these interfaces, so the same function works in
// either context:
//
//	func adjust(w blob.TypedWriter[Profile], key string, delta int) error {
//	    p, err := w.Get(key)
//	    if err != nil { return err }
//	    p.Score += delta
//	    return w.Set(key, p)
//	}
//
//	// standalone — auto-wraps in its own transaction
//	adjust(players, "player:1", 10)
//
//	// composed into a larger transaction
//	func transferTeam(w blob.TypedWriter[Profile], from, to []string, pts int) error {
//	    for _, k := range from {
//	        if err := adjust(w, k, -pts); err != nil { return err }
//	    }
//	    for _, k := range to {
//	        if err := adjust(w, k, pts); err != nil { return err }
//	    }
//	    return nil
//	}
package blob
