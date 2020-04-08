package xposts

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	"github.com/desmos-labs/desmos/x/ibc/xposts/internal/types"
	"github.com/desmos-labs/desmos/x/posts"
)

// NewHandler returns sdk.Handler for cross-chain posts creation module messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		// TODO: Add here
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized cross chain posts creation message type: %T", msg)
		}
	}
}

// handleCreationPacketData handles a MsgPacket containing a PostCreationData
func handleCreationPacketData(
	ctx sdk.Context, k Keeper, pk posts.Keeper, packet channeltypes.Packet, data posts.PostCreationData,
) (*sdk.Result, error) {
	result, err := posts.HandlePostCreationRequest(ctx, pk, data)
	if err != nil {

		if err := k.ChanCloseInit(ctx, packet.DestinationPort, packet.DestinationChannel); err != nil {
			return nil, err
		}
		return nil, err
	}

	acknowledgement := types.AckDataCreation{}.GetBytes()
	if err := k.PacketExecuted(ctx, packet, acknowledgement); err != nil {
		return nil, err
	}

	return result, err
}

// See onTimeoutPacket in spec: https://github.com/cosmos/ics/tree/master/spec/ics-020-fungible-token-transfer#packet-relay
func handleTimeoutDataTransfer(
	ctx sdk.Context, k Keeper, packet channeltypes.Packet, data posts.PostCreationData,
) (*sdk.Result, error) {
	if err := k.TimeoutTransfer(ctx, packet, data); err != nil {
		// This shouldn't happen, since we've already validated that we've sent the packet.
		panic(err)
	}

	if err := k.TimeoutExecuted(ctx, packet); err != nil {
		// This shouldn't happen, since we've already validated that we've sent the packet.
		// TODO: Figure out what happens if the capability authorisation changes.
		panic(err)
	}

	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}
