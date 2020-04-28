package posts

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	"github.com/desmos-labs/desmos/x/ibc/posts/internal/types"
	"github.com/desmos-labs/desmos/x/posts"
)

// NewHandler returns sdk.Handler for cross-chain posts creation module messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgCrossPost:
			return handleMsgCrossPost(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized cross chain posts creation message type: %T", msg)
		}
	}
}

// handleCreationPacketData handles a MsgCrossPost
func handleMsgCrossPost(ctx sdk.Context, k Keeper, msg types.MsgCrossPost) (*sdk.Result, error) {
	if err := k.SendPostCreation(ctx, msg.SourcePort, msg.SourceChannel, msg.DestHeight, msg.PostData); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(AttributeKeyCreator, msg.PostData.Creator),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
	)

	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

// handleCreationPacketData handles a MsgPacket containing a PostCreationData
func handleCreationPacketData(
	ctx sdk.Context, k Keeper, packet channeltypes.Packet, data PostCreationPacketData,
) (*sdk.Result, error) {
	acknowledgement := types.PostCreationPacketAcknowledgement{
		Success: true,
		Error:   "",
	}

	_, err := posts.HandlePostCreationRequest(ctx, k.PostsKeeper, data.PostCreationData)
	if err != nil {
		acknowledgement = types.PostCreationPacketAcknowledgement{
			Success: false,
			Error:   err.Error(),
		}
	}

	if err := k.PacketExecuted(ctx, packet, acknowledgement.GetBytes()); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypePacket,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(AttributeKeyCreator, data.Creator),
		),
	)

	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}
