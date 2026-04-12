# dps-go

Go library for dumb partition storage, backed by SQLite.

A store holds named partitions. Each partition is a flat key-value namespace,
typed by wrapping with a codec:

```go
import (
    "github.com/numberoverzero/dps-go/blob"
)

type Profile struct {
    Name  string `json:"name"`
    Score int    `json:"score"`
}

store, err := blob.Open("app.db")
defer store.Close()

codec := blob.JSONCodec[Profile]()
writer, err := store.OpenWriter("players")
players := blob.NewTypedWriter(writer, codec)

err = players.Set("player:42", Profile{Name: "Ava", Score: 100})

profile, err := players.Get("player:42")
```

[Package docs](https://pkg.go.dev/github.com/numberoverzero/dps-go/blob) cover
codecs, transactions, and writing reusable operations.
