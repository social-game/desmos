package types

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CreatePostPacketData represents the packet data that should be sent when
// wanting to create a new post
type CreatePostPacketData struct {
	PostCreationData // Include all the standard data
}

// NewCreatePostPacketData is the builder function for a new CreatePostPacketData
func NewCreatePostPacketData(data PostCreationData) CreatePostPacketData {
	return CreatePostPacketData{
		PostCreationData: data,
	}
}

// String returns a string representation of FungibleTokenPacketData
func (cppd CreatePostPacketData) String() string {
	return fmt.Sprintf(`CreatePostPacketData:
	%s`,
		cppd.PostCreationData,
	)
}

// ValidateBasic implements channelexported.PacketDataI
func (cppd CreatePostPacketData) ValidateBasic() error {
	return cppd.PostCreationData.ValidateBasic()
}

// GetBytes implements channelexported.PacketDataI
func (cppd CreatePostPacketData) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(cppd))
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
