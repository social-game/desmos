package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/desmos-labs/desmos/x/posts/internal/types"
)

func handleMsgRegisterReaction(ctx sdk.Context, keeper Keeper, msg types.MsgRegisterReaction) (*sdk.Result, error) {
	if _, isAlreadyRegistered := keeper.DoesReactionForShortCodeExist(ctx, msg.ShortCode, msg.Subspace); isAlreadyRegistered {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf(
			"reaction with shortcode %s and subspace %s has already been registered", msg.ShortCode, msg.Subspace))
	}

	reaction := types.NewReaction(msg.Creator, msg.ShortCode, msg.Value, msg.Subspace)

	keeper.RegisterReaction(ctx, reaction)

	event := sdk.NewEvent(
		types.EventTypeRegisterReaction,
		sdk.NewAttribute(types.AttributeKeyReactionCreator, msg.Creator.String()),
		sdk.NewAttribute(types.AttributeKeyReactionShortCode, msg.ShortCode),
		sdk.NewAttribute(types.AttributeKeyReactionSubSpace, msg.Subspace),
	)
	ctx.EventManager().EmitEvent(event)

	result := sdk.Result{
		Data:   []byte("reaction registered properly"),
		Events: sdk.Events{event},
	}

	return &result, nil
}

// handleMsgAddPostReaction handles the adding of a reaction to a post
func handleMsgAddPostReaction(ctx sdk.Context, keeper Keeper, msg types.MsgAddPostReaction) (*sdk.Result, error) {

	// Get the post
	post, found := keeper.GetPost(ctx, msg.PostID)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("post with id %s not found", msg.PostID))
	}

	// Create and store the reaction
	reaction := types.NewPostReaction(msg.Value, msg.User)

	if err := keeper.SavePostReaction(ctx, post.PostID, reaction); err != nil {
		return nil, err
	}

	// Emit the event
	event := sdk.NewEvent(
		types.EventTypePostReactionAdded,
		sdk.NewAttribute(types.AttributeKeyPostID, msg.PostID.String()),
		sdk.NewAttribute(types.AttributeKeyPostReactionOwner, msg.User.String()),
		sdk.NewAttribute(types.AttributeKeyPostReactionValue, msg.Value),
	)
	ctx.EventManager().EmitEvent(event)

	result := sdk.Result{
		Data:   []byte("reaction added properly"),
		Events: sdk.Events{event},
	}
	return &result, nil
}

// handleMsgRemovePostReaction handles the removal of a reaction from a post
func handleMsgRemovePostReaction(ctx sdk.Context, keeper Keeper, msg types.MsgRemovePostReaction) (*sdk.Result, error) {

	// Get the post
	post, found := keeper.GetPost(ctx, msg.PostID)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("post with id %s not found", msg.PostID))
	}

	// Remove the reaction
	if err := keeper.RemovePostReaction(ctx, post.PostID, msg.User, msg.Reaction); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	// Emit the event
	event := sdk.NewEvent(
		types.EventTypePostReactionRemoved,
		sdk.NewAttribute(types.AttributeKeyPostID, msg.PostID.String()),
		sdk.NewAttribute(types.AttributeKeyPostReactionOwner, msg.User.String()),
		sdk.NewAttribute(types.AttributeKeyPostReactionValue, msg.Reaction),
	)
	ctx.EventManager().EmitEvent(event)

	result := sdk.Result{
		Data:   []byte("reaction removed properly"),
		Events: sdk.Events{event},
	}
	return &result, nil
}
