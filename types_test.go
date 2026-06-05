package mezonlight

import (
	"reflect"
	"testing"

	"github.com/quangledang23/mezon-light-sdk-go/proto"
)

func TestNewChannelMessageFromProto(t *testing.T) {
	attachments := &proto.MessageAttachmentList{
		Attachments: []*proto.MessageAttachment{{Filename: "doc.pdf", URL: "https://cdn.test/doc.pdf"}},
	}
	pm := &proto.ChannelMessage{
		ClanID:            "100",
		ChannelID:         "200",
		MessageID:         "300",
		Code:              1,
		SenderID:          "400",
		Username:          "alice",
		Avatar:            "https://cdn.test/a.png",
		Content:           `{"t":"hello"}`,
		ChannelLabel:      "general",
		DisplayName:       "Alice",
		Attachments:       attachments.Marshal(),
		CreateTimeSeconds: 1700000000,
		Mode:              4,
		HideEditted:       true,
		IsPublic:          true,
		TopicID:           "500",
	}

	got := newChannelMessageFromProto(pm)

	if got.ID != "300" || got.MessageID != "300" {
		t.Errorf("ID = %q, MessageID = %q, want both %q", got.ID, got.MessageID, "300")
	}
	if got.ChannelID != "200" || got.ClanID != "100" || got.SenderID != "400" || got.TopicID != "500" {
		t.Error("identifier fields not copied")
	}
	wantContent := map[string]any{"t": "hello"}
	if !reflect.DeepEqual(got.Content, wantContent) {
		t.Errorf("Content = %#v, want %#v", got.Content, wantContent)
	}
	if len(got.Attachments) != 1 || got.Attachments[0].Filename != "doc.pdf" {
		t.Errorf("Attachments = %+v, want decoded protobuf list", got.Attachments)
	}
	if got.Username != "alice" || got.DisplayName != "Alice" || got.ChannelLabel != "general" {
		t.Error("descriptive fields not copied")
	}
	if got.Code != 1 || got.Mode != 4 || got.CreateTimeSeconds != 1700000000 {
		t.Error("numeric fields not copied")
	}
	if !got.HideEditted || !got.IsPublic {
		t.Error("boolean fields not copied")
	}
}

func TestNewChannelMessageFromProtoPlainTextContent(t *testing.T) {
	got := newChannelMessageFromProto(&proto.ChannelMessage{Content: "raw text"})
	wantContent := map[string]any{"t": "raw text"}
	if !reflect.DeepEqual(got.Content, wantContent) {
		t.Errorf("Content = %#v, want %#v", got.Content, wantContent)
	}
	if got.Attachments != nil {
		t.Errorf("Attachments = %v, want nil for empty payload", got.Attachments)
	}
}
