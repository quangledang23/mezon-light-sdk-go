package mezonlight

import (
	"encoding/json"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/quangledang23/mezon-light-sdk-go/proto"
)

// SafeJSONParse decodes raw JSON content. On failure (or for empty input) it
// returns map[string]any{"t": <raw string>}, mirroring safeJSONParse in the
// TypeScript SDK.
func SafeJSONParse(raw []byte) any {
	s := string(raw)
	if s == "" || s == "[]" {
		return map[string]any{"t": s}
	}

	var out any
	if err := json.Unmarshal(raw, &out); err == nil {
		return out
	}

	// Retry with bare newlines escaped, as some payloads contain raw control
	// characters inside string literals.
	fixed := strings.NewReplacer("\n", "\\n", "\r", "\\r").Replace(s)
	if err := json.Unmarshal([]byte(fixed), &out); err == nil {
		return out
	}

	return map[string]any{"t": s}
}

// DecodeAttachments decodes a channel message attachments payload, which may
// be either JSON or a protobuf-encoded MessageAttachmentList.
func DecodeAttachments(data []byte) []*proto.MessageAttachment {
	if len(data) == 0 {
		return nil
	}

	// '[' (JSON array) or '{' (JSON object) marks a JSON payload.
	if data[0] == '[' || data[0] == '{' {
		var list []*proto.MessageAttachment
		if err := json.Unmarshal(data, &list); err == nil {
			return list
		}
		var wrapper struct {
			Attachments []*proto.MessageAttachment `json:"attachments"`
		}
		if err := json.Unmarshal(data, &wrapper); err == nil {
			return wrapper.Attachments
		}
		return nil
	}

	list := &proto.MessageAttachmentList{}
	if err := list.Unmarshal(data); err == nil {
		return list.Attachments
	}
	return nil
}

// DecodeMentions decodes a channel message mentions payload, which may be
// either JSON or a protobuf-encoded MessageMentionList.
func DecodeMentions(data []byte) []*proto.MessageMention {
	if len(data) == 0 {
		return nil
	}

	// '[' (JSON array) or '{' (JSON object) marks a JSON payload.
	if data[0] == '[' || data[0] == '{' {
		var list []*proto.MessageMention
		if err := json.Unmarshal(data, &list); err == nil {
			return list
		}
		var wrapper struct {
			Mentions []*proto.MessageMention `json:"mentions"`
		}
		if err := json.Unmarshal(data, &wrapper); err == nil {
			return wrapper.Mentions
		}
		return nil
	}

	list := &proto.MessageMentionList{}
	if err := list.Unmarshal(data); err == nil {
		return list.Mentions
	}
	return nil
}

var linkRegexp = regexp.MustCompile(`https?://\S+`)

// ExtractLinks finds http(s) URLs in text and returns "lk" markup tokens for
// them, which clients render as clickable links. Offsets are characters, not
// bytes, matching how Mezon clients index content.
func ExtractLinks(text string) []*MessageMarkup {
	var marks []*MessageMarkup
	for _, loc := range linkRegexp.FindAllStringIndex(text, -1) {
		url := strings.TrimRight(text[loc[0]:loc[1]], `.,;:!?'")]}`)
		if url == "" {
			continue
		}
		s := utf8.RuneCountInString(text[:loc[0]])
		marks = append(marks, &MessageMarkup{
			Type: MarkupTypeLink,
			S:    int32(s),
			E:    int32(s + utf8.RuneCountInString(url)),
		})
	}
	return marks
}

// NewTextContent builds message content from plain text, marking any URLs in
// it as clickable links.
func NewTextContent(text string) *MessageContent {
	return &MessageContent{T: text, Mk: ExtractLinks(text)}
}
