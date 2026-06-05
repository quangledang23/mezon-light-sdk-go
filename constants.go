// Package mezonlight is a lightweight Go SDK for Mezon chat, ported from the
// TypeScript package mezon-light-sdk.
package mezonlight

import "time"

const (
	// MezonGWURL is the default Mezon Gateway URL.
	MezonGWURL = "https://gw.mezon.ai"

	// SocketReadyMaxRetry is the maximum number of retries when waiting for
	// the socket to be ready.
	SocketReadyMaxRetry = 20

	// SocketReadyRetryDelay is the initial delay between socket ready
	// retries (uses exponential backoff).
	SocketReadyRetryDelay = 100 * time.Millisecond

	// ClanDM is the clan ID used for Direct Messages.
	ClanDM = "0"

	// ChannelTypeDM is the channel type for Direct Messages.
	ChannelTypeDM = 3

	// ChannelTypeGroup is the channel type for group DMs.
	ChannelTypeGroup = 2

	// StreamModeDM is the stream mode for Direct Messages.
	StreamModeDM = 4

	// StreamModeGroup is the stream mode for group DMs.
	StreamModeGroup = 3

	// DefaultServerKey is the default server key if none is provided.
	DefaultServerKey = "DefaultServerKey"
)
