package types

const (
	// AcceptSliceCapacity contains the initial capacity for slices to store MIME types and weights.
	AcceptSliceCapacity = 10

	// AcceptQualityWeight corresponds the default quality weight for MIME types.
	AcceptQualityWeight = 1.0

	// AcceptQualityParameter is used to specify quality weight parameter in the header.
	AcceptQualityParameter = "q"

	// AcceptSeparator is the used Separator for multiple MIME types in the header.
	AcceptSeparator = ","

	// AcceptHeader is the HTTP header key for the Accept header.
	AcceptHeader = "Accept"

	// SecWebsocketProtocol is the HTTP header key for the Websocket sub-protocol.
	SecWebsocketProtocol = "Sec-Websocket-Protocol"
)
