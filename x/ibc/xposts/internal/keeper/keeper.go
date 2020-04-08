package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/capability"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channelexported "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
	"github.com/desmos-labs/desmos/x/posts"
)

type Keeper struct {
	postsKeeper   posts.Keeper
	channelKeeper channel.Keeper
	scopedKeeper  capability.ScopedKeeper
}

func NewKeeper(pk posts.Keeper, ck channel.Keeper, sk capability.ScopedKeeper) Keeper {
	return Keeper{
		postsKeeper:   pk,
		channelKeeper: ck,
		scopedKeeper:  sk,
	}
}

// PacketExecuted defines a wrapper function for the channel Keeper's function
// in order to expose it to the ICS20 transfer handler.
func (k Keeper) PacketExecuted(ctx sdk.Context, packet channelexported.PacketI, acknowledgement []byte) error {
	return k.channelKeeper.PacketExecuted(ctx, packet, acknowledgement)
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

// ClaimCapability allows the transfer module that can claim a capability that IBC module
// passes to it
func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capability.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}

// TimeoutTransfer handles post creation timeout logic.
func (k Keeper) TimeoutTransfer(ctx sdk.Context, packet channel.Packet, data posts.PostCreationData) error {
	// TODO: Implement
	return nil
	//return k.onTimeoutPacket(ctx, packet, data)
}

// TimeoutExecuted defines a wrapper function for the channel Keeper's function
// in order to expose it to the ICS20 transfer handler.
func (k Keeper) TimeoutExecuted(ctx sdk.Context, packet channelexported.PacketI) error {
	return k.channelKeeper.TimeoutExecuted(ctx, packet)
}
