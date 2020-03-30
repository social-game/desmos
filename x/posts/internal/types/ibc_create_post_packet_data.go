package types

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channelexported "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
	transferTypes "github.com/cosmos/cosmos-sdk/x/ibc/20-transfer/types"
)

var _ channelexported.PacketDataI = CreatePostPacketData{}

type CreatePostPacketData struct {
	CreationData
	Timeout uint64 `json:"timeout" yaml:"timeout"`
}

// String returns a string representation of FungibleTokenPacketData
func (cppd CreatePostPacketData) String() string {
	return fmt.Sprintf(`CreatePostPacketData:
	%s
	Timeout:            %d`,
		cppd.CreationData,
		cppd.Timeout,
	)
}

// ValidateBasic implements channelexported.PacketDataI
func (cppd CreatePostPacketData) ValidateBasic() error {
	if err := cppd.CreationData.ValidateBasic(); err != nil {
		return err
	}

	if cppd.Timeout == 0 {
		return sdkerrors.Wrap(transferTypes.ErrInvalidPacketTimeout, "timeout cannot be 0")
	}

	return nil
}

// GetBytes implements channelexported.PacketDataI
func (cppd CreatePostPacketData) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(cppd))
}

// GetTimeoutHeight implements channelexported.PacketDataI
func (cppd CreatePostPacketData) GetTimeoutHeight() uint64 {
	return cppd.Timeout
}

// Type implements channelexported.PacketDataI
func (cppd CreatePostPacketData) Type() string {
	return "posts/create"
}

// MarshalJSON implements the json.Marshaler interface.
// This is done due to the fact that Amino does not respect omitempty clauses
func (cppd CreatePostPacketData) MarshalJSON() ([]byte, error) {
	type temp CreatePostPacketData
	return json.Marshal(temp(cppd))
}

// AckDataCreation is a no-op packet
// See spec for onAcknowledgePacket: https://github.com/cosmos/ics/tree/master/spec/ics-020-fungible-token-transfer#packet-relay
type AckDataCreation struct{}

// GetBytes implements channelexported.PacketAcknowledgementI
func (ack AckDataCreation) GetBytes() []byte {
	return []byte("post creation ack")
}
