package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	host "github.com/cosmos/cosmos-sdk/x/ibc/24-host"
)

// MsgCrossPost allows to create a new post inside another chain using IBC
type MsgCrossPost struct {
	SourcePort    string `json:"source_port" yaml:"source_port"`       // the port on which the packet will be sent
	SourceChannel string `json:"source_channel" yaml:"source_channel"` // the channel by which the packet will be sent
	DestHeight    uint64 `json:"dest_height" yaml:"dest_height"`       // the current height of the destination chain

	PostData PostCreationPacketData `json:"post_data"`
	Sender   sdk.AccAddress         `json:"sender"`
}

// NewMsgCrossPost allows to create a new MsgCrossPost created from the given
// sender and containing the given post data.
func NewMsgCrossPost(
	sourcePort, sourceChannel string, destHeight uint64,
	postData PostCreationPacketData, sender sdk.AccAddress,
) MsgCrossPost {
	return MsgCrossPost{
		SourcePort:    sourcePort,
		SourceChannel: sourceChannel,
		DestHeight:    destHeight,
		PostData:      postData,
		Sender:        sender,
	}
}

// Route implements sdk.Msg
func (MsgCrossPost) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (MsgCrossPost) Type() string {
	return "cross-post"
}

// ValidateBasic implements sdk.Msg
func (msg MsgCrossPost) ValidateBasic() error {
	if err := host.DefaultPortIdentifierValidator(msg.SourcePort); err != nil {
		return sdkerrors.Wrap(err, "invalid source port ID")
	}

	if err := host.DefaultChannelIdentifierValidator(msg.SourceChannel); err != nil {
		return sdkerrors.Wrap(err, "invalid source channel ID")
	}

	if len(msg.Sender) == 0 {
		return fmt.Errorf("invalid msg sender: %s", msg.Sender)
	}

	return msg.PostData.ValidateBasic()
}

// GetSignBytes implements sdk.Msg
func (msg MsgCrossPost) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgCrossPost) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
