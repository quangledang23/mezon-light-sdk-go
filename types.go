package mezonlight

import (
	"github.com/mezonai/mezon-light-sdk-go/proto"
)

// Aliases re-exporting the wire-level types used in the public API, mirroring
// the Api* names of the TypeScript SDK.
type (
	ApiSession                  = proto.Session
	ApiSessionRefreshRequest    = proto.SessionRefreshRequest
	ApiChannelDescription       = proto.ChannelDescription
	ApiCreateChannelDescRequest = proto.CreateChannelDescRequest
	ApiMessageAttachment        = proto.MessageAttachment
	ApiUploadAttachment         = proto.UploadAttachment
	ApiUploadAttachmentRequest  = proto.UploadAttachmentRequest
)

// ClientInitConfig configures a LightClient created from existing tokens.
type ClientInitConfig struct {
	// Token is the authentication token.
	Token string `json:"token"`
	// RefreshToken is used for session renewal.
	RefreshToken string `json:"refresh_token"`
	// APIURL is the API URL for the Mezon server.
	APIURL string `json:"api_url"`
	// WSURL is the WebSocket host for session connectivity.
	WSURL string `json:"ws_url"`
	// UserID is the user ID associated with the session.
	UserID string `json:"user_id"`
	// ServerKey is the server key for authentication (optional, uses
	// DefaultServerKey if empty).
	ServerKey string `json:"serverkey,omitempty"`
}

// AuthenticateConfig configures authentication of a new user.
type AuthenticateConfig struct {
	// IDToken is the ID token from an identity provider.
	IDToken string `json:"id_token"`
	// UserID is the user ID to associate with the account.
	UserID string `json:"user_id"`
	// Username is the username for the account.
	Username string `json:"username"`
	// ServerKey is the server key for authentication (optional).
	ServerKey string `json:"serverkey,omitempty"`
	// GatewayURL is a custom gateway URL (optional, uses MezonGWURL if empty).
	GatewayURL string `json:"gateway_url,omitempty"`
}

// SendMessagePayload describes a message to send.
type SendMessagePayload struct {
	// ChannelID is the channel to send the message to.
	ChannelID string
	// Content is the message content; it is JSON-encoded before sending.
	Content any
	// Attachments holds optional file/media attachments.
	Attachments []*ApiMessageAttachment
	// HideLink hides link previews when true.
	HideLink bool
	// Code is the optional message code.
	Code int32
}

// ApiAuthenticationIdToken is the request body for ID-token authentication.
type ApiAuthenticationIdToken struct {
	// IDToken is the ID token from an identity provider.
	IDToken string `json:"id_token"`
	// UserID is the user ID associated with the token.
	UserID string `json:"user_id"`
	// Username is the username associated with the token.
	Username string `json:"username"`
}

// AuthenticationIdTokenResponse is the response of ID-token authentication.
type AuthenticationIdTokenResponse struct {
	// Token is the authentication token.
	Token string `json:"token"`
	// RefreshToken is used for session renewal.
	RefreshToken string `json:"refresh_token"`
	// APIURL is the API URL for the authenticated user.
	APIURL string `json:"api_url"`
	// WSURL is the WS host for the authenticated user.
	WSURL string `json:"ws_url"`
	// UserID is the user ID of the authenticated user.
	UserID string `json:"user_id"`
}

// ChannelMessage is a message received on a channel, with content and
// attachments already decoded (the wire-level counterpart is
// proto.ChannelMessage).
type ChannelMessage struct {
	ID                string                  `json:"id"`
	Avatar            string                  `json:"avatar,omitempty"`
	ChannelID         string                  `json:"channel_id"`
	ChannelLabel      string                  `json:"channel_label"`
	ClanID            string                  `json:"clan_id,omitempty"`
	Code              int32                   `json:"code"`
	Content           any                     `json:"content"`
	Attachments       []*ApiMessageAttachment `json:"attachments,omitempty"`
	SenderID          string                  `json:"sender_id"`
	ClanLogo          string                  `json:"clan_logo,omitempty"`
	CategoryName      string                  `json:"category_name,omitempty"`
	Username          string                  `json:"username,omitempty"`
	ClanNick          string                  `json:"clan_nick,omitempty"`
	ClanAvatar        string                  `json:"clan_avatar,omitempty"`
	DisplayName       string                  `json:"display_name,omitempty"`
	CreateTimeSeconds uint32                  `json:"create_time_seconds,omitempty"`
	UpdateTimeSeconds uint32                  `json:"update_time_seconds,omitempty"`
	Mode              int32                   `json:"mode,omitempty"`
	MessageID         string                  `json:"message_id,omitempty"`
	HideEditted       bool                    `json:"hide_editted,omitempty"`
	IsPublic          bool                    `json:"is_public,omitempty"`
	TopicID           string                  `json:"topic_id,omitempty"`
}

// newChannelMessageFromProto mirrors createChannelMessageFromEvent in the
// TypeScript SDK: it decodes content (JSON) and attachments (JSON or
// protobuf MessageAttachmentList).
func newChannelMessageFromProto(pm *proto.ChannelMessage) *ChannelMessage {
	return &ChannelMessage{
		ID:                pm.MessageID,
		Avatar:            pm.Avatar,
		ChannelID:         pm.ChannelID,
		ChannelLabel:      pm.ChannelLabel,
		ClanID:            pm.ClanID,
		Code:              pm.Code,
		Content:           SafeJSONParse([]byte(pm.Content)),
		Attachments:       DecodeAttachments(pm.Attachments),
		SenderID:          pm.SenderID,
		ClanLogo:          pm.ClanLogo,
		CategoryName:      pm.CategoryName,
		Username:          pm.Username,
		ClanNick:          pm.ClanNick,
		ClanAvatar:        pm.ClanAvatar,
		DisplayName:       pm.DisplayName,
		CreateTimeSeconds: pm.CreateTimeSeconds,
		UpdateTimeSeconds: pm.UpdateTimeSeconds,
		Mode:              pm.Mode,
		MessageID:         pm.MessageID,
		HideEditted:       pm.HideEditted,
		IsPublic:          pm.IsPublic,
		TopicID:           pm.TopicID,
	}
}
