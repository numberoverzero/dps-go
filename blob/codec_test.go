package blob

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type jsonTestStruct struct {
	Name     string          `json:"name"`
	Score    int             `json:"score"`
	Tags     []string        `json:"tags,omitempty"`
	Metadata map[string]int  `json:"metadata,omitempty"`
	Nested   *jsonTestStruct `json:"nested,omitempty"`
}

func TestJSON_Marshal(t *testing.T) {
	codec := JSONCodec[jsonTestStruct]()

	tests := []struct {
		name  string
		value jsonTestStruct
	}{
		{"zero value", jsonTestStruct{}},
		{"populated", jsonTestStruct{
			Name:     "Ava",
			Score:    42,
			Tags:     []string{"a", "b"},
			Metadata: map[string]int{"x": 1},
			Nested:   &jsonTestStruct{Name: "inner"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// codec.marshal output should be valid for encoding/json
			data, err := codec.marshal(tt.value)
			require.NoError(t, err)

			var got jsonTestStruct
			require.NoError(t, json.Unmarshal(data, &got))
			assert.Equal(t, tt.value, got)
		})
	}
}

func TestJSON_Unmarshal(t *testing.T) {
	codec := JSONCodec[jsonTestStruct]()

	tests := []struct {
		name    string
		data    []byte
		want    jsonTestStruct
		wantErr bool
	}{
		{
			name: "valid json",
			data: []byte(`{"name":"Ava","score":42}`),
			want: jsonTestStruct{Name: "Ava", Score: 42},
		},
		{
			name: "encoding/json output",
			data: func() []byte {
				b, _ := json.Marshal(jsonTestStruct{Name: "Bob", Tags: []string{"x"}})

				return b
			}(),
			want: jsonTestStruct{Name: "Bob", Tags: []string{"x"}},
		},
		{
			name: "unknown fields ignored",
			data: []byte(`{"name":"Ava","extra":"ignored"}`),
			want: jsonTestStruct{Name: "Ava"},
		},
		{
			name:    "not json",
			data:    []byte("not json"),
			wantErr: true,
		},
		{
			name:    "truncated",
			data:    []byte(`{"name":`),
			wantErr: true,
		},
		{
			name:    "empty",
			data:    []byte{},
			wantErr: true,
		},
		{
			name:    "wrong type",
			data:    []byte(`"just a string"`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := codec.unmarshal(tt.data)
			if tt.wantErr {
				assert.Error(t, err)

				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestJSON_Roundtrip(t *testing.T) {
	codec := JSONCodec[jsonTestStruct]()

	tests := []struct {
		name  string
		value jsonTestStruct
	}{
		{"zero value", jsonTestStruct{}},
		{"full", jsonTestStruct{
			Name:     "日本語",
			Score:    -1,
			Tags:     []string{"a", "b"},
			Metadata: map[string]int{"x": 1, "y": 2},
			Nested:   &jsonTestStruct{Name: "inner", Score: 99},
		}},
		{"nil collections", jsonTestStruct{Name: "bare"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := codec.marshal(tt.value)
			require.NoError(t, err)

			got, err := codec.unmarshal(data)
			require.NoError(t, err)
			assert.Equal(t, tt.value, got)
		})
	}
}

// empty collections with omitempty round-trip to nil (standard encoding/json behavior)
func TestJSON_OmitEmptyLosesEmpty(t *testing.T) {
	codec := JSONCodec[jsonTestStruct]()
	original := jsonTestStruct{Tags: []string{}, Metadata: map[string]int{}}

	data, err := codec.marshal(original)
	require.NoError(t, err)

	got, err := codec.unmarshal(data)
	require.NoError(t, err)

	assert.Nil(t, got.Tags, "empty slice becomes nil after omitempty roundtrip")
	assert.Nil(t, got.Metadata, "empty map becomes nil after omitempty roundtrip")
}

func TestJSON_MarshalError(t *testing.T) {
	codec := JSONCodec[any]()

	tests := []struct {
		name  string
		value any
	}{
		{"channel", make(chan int)},
		{"func", func() {}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := codec.marshal(tt.value)
			assert.Error(t, err)
		})
	}
}

func TestRaw(t *testing.T) {
	tests := []struct {
		name  string
		value []byte
	}{
		{"content", []byte("hello world")},
		{"binary", []byte{0x00, 0xFF, 0x80, 0x01}},
		{"empty", []byte{}},
		{"nil", nil},
		{"large", make([]byte, 1<<16)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := RawCodec.marshal(tt.value)
			require.NoError(t, err)
			assert.Equal(t, tt.value, data)

			got, err := RawCodec.unmarshal(data)
			require.NoError(t, err)
			assert.Equal(t, tt.value, got)
		})
	}
}
