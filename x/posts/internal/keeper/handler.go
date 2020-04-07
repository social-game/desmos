package keeper

import (
	"fmt"

	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	"github.com/desmos-labs/desmos/x/posts/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns a handler for "magpie" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgCreatePost:
			return handlePostCreationRequest(ctx, keeper, msg.PostCreationData)
		case types.MsgEditPost:
			return handleMsgEditPost(ctx, keeper, msg)
		case types.MsgAddPostReaction:
			return handleMsgAddPostReaction(ctx, keeper, msg)
		case types.MsgRemovePostReaction:
			return handleMsgRemovePostReaction(ctx, keeper, msg)
		case types.MsgAnswerPoll:
			return handleMsgAnswerPollPost(ctx, keeper, msg)
		case types.MsgRegisterReaction:
			return handleMsgRegisterReaction(ctx, keeper, msg)

		case channeltypes.MsgPacket:
			var data types.PostCreationData
			if err := types.ModuleCdc.UnmarshalJSON(msg.GetData(), &data); err != nil {
				return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal post creation transfer packet data: %s", err.Error())
			}
			return handleCreationPacketData(ctx, keeper, msg, data)

		default:
			errMsg := fmt.Sprintf("Unrecognized Posts message type: %v", msg.Type())
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// handlePostCreationRequest handles the creation of a new post
func handlePostCreationRequest(ctx sdk.Context, keeper Keeper, data types.PostCreationData) (*sdk.Result, error) {
	post := types.NewPost(
		keeper.GetLastPostID(ctx).Next(),
		data.ParentID,
		data.Message,
		data.AllowsComments,
		data.Subspace,
		data.OptionalData,
		data.CreationDate,
		data.Creator,
	).WithMedias(data.Medias)

	if data.PollData != nil {
		post = post.WithPollData(*data.PollData)
	}

	// Check for double posting
	if existing, found := keeper.IsPostConflicting(ctx, post); found {
		msg := `the provided post conflicts with the one having id %s. Please check that either their creation date, subspace or creator are different`
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf(msg, existing.PostID))
	}

	// If valid, check the parent post
	if post.ParentID.Valid() {
		parentPost, found := keeper.GetPost(ctx, post.ParentID)
		if !found {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("parent post with id %s not found", post.ParentID))
		}

		if !parentPost.AllowsComments {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("post with id %s does not allow comments", parentPost.PostID))
		}
	}

	keeper.SavePost(ctx, post)

	createEvent := sdk.NewEvent(
		types.EventTypePostCreated,
		sdk.NewAttribute(types.AttributeKeyPostID, post.PostID.String()),
		sdk.NewAttribute(types.AttributeKeyPostParentID, post.ParentID.String()),
		sdk.NewAttribute(types.AttributeKeyCreationTime, post.Created.String()),
		sdk.NewAttribute(types.AttributeKeyPostOwner, post.Creator.String()),
	)
	ctx.EventManager().EmitEvent(createEvent)

	result := sdk.Result{
		Data:   keeper.Cdc.MustMarshalBinaryLengthPrefixed(post.PostID),
		Events: sdk.Events{createEvent}.ToABCIEvents(),
	}
	return &result, nil
}

// handleCreationPacketData handles a MsgPacket containing a CreatePostPacketData
func handleCreationPacketData(
	ctx sdk.Context, k Keeper, msg channeltypes.MsgPacket, data types.PostCreationData,
) (*sdk.Result, error) {
	result, err := handlePostCreationRequest(ctx, k, data)
	if err != nil {

		if err := k.ChanCloseInit(ctx, msg.Packet.DestinationPort, msg.Packet.DestinationChannel); err != nil {
			return nil, err
		}
		return nil, err
	}

	acknowledgement := types.AckDataCreation{}.GetBytes()
	if err := k.PacketExecuted(ctx, msg.Packet, acknowledgement); err != nil {
		return nil, err
	}

	return result, err
}

// handleMsgEditPost handles the edit of posts
func handleMsgEditPost(ctx sdk.Context, keeper Keeper, msg types.MsgEditPost) (*sdk.Result, error) {

	// Get the existing post
	existing, found := keeper.GetPost(ctx, msg.PostID)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("post with id %s not found", msg.PostID))
	}

	// Checks if the the msg sender is the same as the current owner
	if !msg.Editor.Equals(existing.Creator) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	// Check the validity of the current block height respect to the creation date of the post
	if existing.Created.After(msg.EditDate) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "edit date cannot be before creation date")
	}

	// Edit the post
	existing.Message = msg.Message
	existing.LastEdited = msg.EditDate
	keeper.SavePost(ctx, existing)

	editEvent := sdk.NewEvent(
		types.EventTypePostEdited,
		sdk.NewAttribute(types.AttributeKeyPostID, existing.PostID.String()),
		sdk.NewAttribute(types.AttributeKeyPostEditTime, existing.LastEdited.String()),
	)
	ctx.EventManager().EmitEvent(editEvent)

	result := sdk.Result{
		Data:   keeper.Cdc.MustMarshalBinaryLengthPrefixed(existing.PostID),
		Events: sdk.Events{editEvent}.ToABCIEvents(),
	}
	return &result, nil
}
