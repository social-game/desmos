package types

import (
	"fmt"

	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
)

// IBC transfer events
const (
	EventTypeTimeout      = "timeout"
	EventTypePacket       = "create_post_packet"
	EventTypeChannelClose = "channel_closed"

	AttributeKeyCreator = "creator"
	AttributeKeyPostID  = "post_id"

	AttributeKeyAckSuccess = "success"
	AttributeKeyAckError   = "error"
)

// IBC transfer events vars
var (
	AttributeValueCategory = fmt.Sprintf("%s_%s", ibctypes.ModuleName, ModuleName)
)
