package types

import (
	"encoding/json"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channelexported "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
	transferTypes "github.com/cosmos/cosmos-sdk/x/ibc/20-transfer/types"
	"github.com/desmos-labs/desmos/x/commons"
	"github.com/tendermint/tendermint/libs/bech32"
)

var _ channelexported.PacketDataI = CreatePostPacketData{}

// DesmosAddress is a wrapper around sdk.AccAddress to make sure that it is properly serialized
// using the commons.Bech32MainPrefix prefix while getting signature bytes
type DesmosAddress struct {
	sdk.AccAddress
}

// MarshalJSON marshals to JSON using Bech32.
func (aa DesmosAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(aa.String())
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (aa *DesmosAddress) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	var aa2 sdk.AccAddress
	if len(strings.TrimSpace(s)) == 0 {
		aa2 = nil
	}

	bz, err := sdk.GetFromBech32(s, commons.Bech32MainPrefix)
	if err != nil {
		return err
	}

	err = sdk.VerifyAddressFormat(bz)
	if err != nil {
		return err
	}
	aa2 = bz

	*aa = DesmosAddress{aa2}
	return nil
}

// String implements the Stringer interface.
func (aa DesmosAddress) String() string {
	if aa.Empty() {
		return ""
	}

	bech32Addr, err := bech32.ConvertAndEncode(commons.Bech32MainPrefix, aa.Bytes())
	if err != nil {
		panic(err)
	}

	return bech32Addr
}

// CreatePostPacketData represents the packet data that should be sent when
// wanting to create a new post
type CreatePostPacketData struct {
	PostCreationData               // Include all the standard data
	Creator          DesmosAddress `json:"creator"`                // Override the creator to make sure it has the proper prefix
	Timeout          uint64        `json:"timeout" yaml:"timeout"` // Timeout of the packet
}

// NewCreatePostPacketData is the builder function for a new CreatePostPacketData
func NewCreatePostPacketData(data PostCreationData, timeout uint64) CreatePostPacketData {
	return CreatePostPacketData{
		PostCreationData: data,
		Creator:          DesmosAddress{data.Creator},
		Timeout:          timeout,
	}
}

// String returns a string representation of FungibleTokenPacketData
func (cppd CreatePostPacketData) String() string {
	return fmt.Sprintf(`CreatePostPacketData:
	%s
	Timeout:            %d`,
		cppd.PostCreationData,
		cppd.Timeout,
	)
}

// ValidateBasic implements channelexported.PacketDataI
func (cppd CreatePostPacketData) ValidateBasic() error {
	if err := cppd.PostCreationData.ValidateBasic(); err != nil {
		return err
	}

	if cppd.Timeout == 0 {
		return sdkerrors.Wrap(transferTypes.ErrInvalidPacketTimeout, "timeout cannot be 0")
	}

	return nil
}

// GetBytes implements channelexported.PacketDataI
func (cppd CreatePostPacketData) GetBytes() []byte {
	json := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(cppd))
	println(string(json))
	return json
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
