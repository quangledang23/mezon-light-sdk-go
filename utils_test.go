package mezonlight

import (
	"reflect"
	"testing"

	"github.com/quangledang23/mezon-light-sdk-go/proto"
)

func TestSafeJSONParse(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want any
	}{
		{
			name: "empty input",
			raw:  "",
			want: map[string]any{"t": ""},
		},
		{
			name: "empty array literal",
			raw:  "[]",
			want: map[string]any{"t": "[]"},
		},
		{
			name: "valid object",
			raw:  `{"t":"hello"}`,
			want: map[string]any{"t": "hello"},
		},
		{
			name: "valid array",
			raw:  `[1,2]`,
			want: []any{float64(1), float64(2)},
		},
		{
			name: "valid scalar",
			raw:  `42`,
			want: float64(42),
		},
		{
			name: "raw newline inside string literal",
			raw:  "{\"t\":\"line1\nline2\"}",
			want: map[string]any{"t": "line1\nline2"},
		},
		{
			name: "raw carriage return inside string literal",
			raw:  "{\"t\":\"a\rb\"}",
			want: map[string]any{"t": "a\rb"},
		},
		{
			name: "invalid JSON falls back to t wrapper",
			raw:  "plain text message",
			want: map[string]any{"t": "plain text message"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SafeJSONParse([]byte(tt.raw))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SafeJSONParse(%q) = %#v, want %#v", tt.raw, got, tt.want)
			}
		})
	}
}

func TestDecodeAttachments(t *testing.T) {
	att := &proto.MessageAttachment{
		Filename: "photo.png",
		Size:     1024,
		URL:      "https://cdn.example.com/photo.png",
		Filetype: "image/png",
		Width:    800,
		Height:   600,
	}

	t.Run("empty input", func(t *testing.T) {
		if got := DecodeAttachments(nil); got != nil {
			t.Errorf("DecodeAttachments(nil) = %v, want nil", got)
		}
		if got := DecodeAttachments([]byte{}); got != nil {
			t.Errorf("DecodeAttachments(empty) = %v, want nil", got)
		}
	})

	t.Run("JSON array", func(t *testing.T) {
		data := []byte(`[{"filename":"photo.png","size":1024,"url":"https://cdn.example.com/photo.png","filetype":"image/png","width":800,"height":600}]`)
		got := DecodeAttachments(data)
		if len(got) != 1 || !reflect.DeepEqual(got[0], att) {
			t.Errorf("DecodeAttachments(JSON array) = %+v, want [%+v]", got, att)
		}
	})

	t.Run("JSON wrapper object", func(t *testing.T) {
		data := []byte(`{"attachments":[{"filename":"photo.png","size":1024,"url":"https://cdn.example.com/photo.png","filetype":"image/png","width":800,"height":600}]}`)
		got := DecodeAttachments(data)
		if len(got) != 1 || !reflect.DeepEqual(got[0], att) {
			t.Errorf("DecodeAttachments(JSON wrapper) = %+v, want [%+v]", got, att)
		}
	})

	t.Run("protobuf MessageAttachmentList", func(t *testing.T) {
		list := &proto.MessageAttachmentList{Attachments: []*proto.MessageAttachment{att}}
		got := DecodeAttachments(list.Marshal())
		if len(got) != 1 || !reflect.DeepEqual(got[0], att) {
			t.Errorf("DecodeAttachments(protobuf) = %+v, want [%+v]", got, att)
		}
	})

	t.Run("invalid JSON array", func(t *testing.T) {
		if got := DecodeAttachments([]byte(`[{"bad`)); got != nil {
			t.Errorf("DecodeAttachments(invalid JSON) = %v, want nil", got)
		}
	})

	t.Run("invalid protobuf", func(t *testing.T) {
		if got := DecodeAttachments([]byte{0xff}); got != nil {
			t.Errorf("DecodeAttachments(invalid protobuf) = %v, want nil", got)
		}
	})
}
