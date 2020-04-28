package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/capability"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channelexported "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
	porttypes "github.com/cosmos/cosmos-sdk/x/ibc/05-port/types"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
	"github.com/desmos-labs/desmos/x/ibc/posts/internal/types"
	"github.com/desmos-labs/desmos/x/posts"
)

const (
	// DefaultPacketTimeout is the default packet timeout relative to the current block height
	DefaultPacketTimeout = 1000 // NOTE: in blocks

	// DefaultPacketTimeoutTimestamp is the default packet timeout timestamp relative
	// to the current block timestamp. The timeout is disabled when set to 0.
	DefaultPacketTimeoutTimestamp = 0 // NOTE: in nanoseconds
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec

	PostsKeeper   posts.Keeper
	channelKeeper types.ChannelKeeper
	portKeeper    types.PortKeeper
	scopedKeeper  capability.ScopedKeeper
}

func NewKeeper(
	cdc *codec.Codec, storeKey sdk.StoreKey, pk posts.Keeper,
	ck types.ChannelKeeper, portK types.PortKeeper, sk capability.ScopedKeeper,
) Keeper {
	return Keeper{
		storeKey: storeKey,
		cdc:      cdc,

		PostsKeeper:   pk,
		channelKeeper: ck,
		portKeeper:    portK,
		scopedKeeper:  sk,
	}
}

// PacketExecuted defines a wrapper function for the channel Keeper's function
// in order to expose it to the ICS20 transfer handler.
func (k Keeper) PacketExecuted(ctx sdk.Context, packet channelexported.PacketI, acknowledgement []byte) error {
	chanCap, ok := k.scopedKeeper.GetCapability(ctx, ibctypes.ChannelCapabilityPath(packet.GetDestPort(), packet.GetDestChannel()))
	if !ok {
		return sdkerrors.Wrap(channel.ErrChannelCapabilityNotFound, "channel capability could not be retrieved for packet")
	}
	return k.channelKeeper.PacketExecuted(ctx, chanCap, packet, acknowledgement)
}

// ChanCloseInit defines a wrapper function for the channel Keeper's function
// in order to expose it to the ICS20 trasfer handler.
func (k Keeper) ChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	capName := ibctypes.ChannelCapabilityPath(portID, channelID)
	chanCap, ok := k.scopedKeeper.GetCapability(ctx, capName)
	if !ok {
		return sdkerrors.Wrapf(channel.ErrChannelCapabilityNotFound, "could not retrieve channel capability at: %s", capName)
	}
	return k.channelKeeper.ChanCloseInit(ctx, portID, channelID, chanCap)
}

// BindPort defines a wrapper function for the ort Keeper's function in
// order to expose it to module's InitGenesis function
func (k Keeper) BindPort(ctx sdk.Context, portID string) error {
	// Set the portID into our store so we can retrieve it later
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.PortKey), []byte(portID))

	chanCap := k.portKeeper.BindPort(ctx, portID)
	return k.ClaimCapability(ctx, chanCap, porttypes.PortPath(portID))
}

// GetPort returns the portID for the IBC posts module.
func (k Keeper) GetPort(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get([]byte(types.PortKey)))
}

// ClaimCapability allows the transfer module that can claim a capability that IBC module
// passes to it
func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capability.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}

// SendPostCreation allows to create an IBC packet containing the data of the
// post to create into a Desmos-based chain through IBC
func (k Keeper) SendPostCreation(
	ctx sdk.Context,
	sourcePort,
	sourceChannel string,
	destHeight uint64,

	data types.PostCreationPacketData,
) error {
	sourceChannelEnd, found := k.channelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
	if !found {
		return sdkerrors.Wrap(channel.ErrChannelNotFound, sourceChannel)
	}

	destinationPort := sourceChannelEnd.Counterparty.PortID
	destinationChannel := sourceChannelEnd.Counterparty.ChannelID

	// get the next sequence
	sequence, found := k.channelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
	if !found {
		return channel.ErrSequenceSendNotFound
	}

	return k.createOutgoingPacket(
		ctx, sequence, sourcePort, sourceChannel, destinationPort, destinationChannel, destHeight,
		data,
	)
}

func (k Keeper) createOutgoingPacket(
	ctx sdk.Context,
	seq uint64,
	sourcePort, sourceChannel,
	destinationPort, destinationChannel string,
	destHeight uint64,

	data types.PostCreationPacketData,
) error {
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, ibctypes.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return sdkerrors.Wrap(channel.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	packet := channel.NewPacket(
		data.GetBytes(),
		seq,
		sourcePort,
		sourceChannel,
		destinationPort,
		destinationChannel,
		destHeight+DefaultPacketTimeout,
		DefaultPacketTimeoutTimestamp,
	)

	return k.channelKeeper.SendPacket(ctx, channelCap, packet)
}
