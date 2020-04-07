package keeper

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/desmos-labs/desmos/x/posts/internal/types"
)

// handleMsgAnswerPollPost handles the answer to a poll post
func handleMsgAnswerPollPost(ctx sdk.Context, keeper Keeper, msg types.MsgAnswerPoll) (*sdk.Result, error) {

	post, err := checkPostPollValid(ctx, msg.PostID, keeper)
	if err != nil {
		return nil, err
	}

	// checks if the post's poll allows multiple answers
	if len(msg.UserAnswers) > 1 && !post.PollData.AllowsMultipleAnswers {
		return nil, sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("the poll associated with ID %s doesn't allow multiple answers",
				post.PostID),
		)
	}

	// check if the user answers are more than the answers provided by the poll
	if len(msg.UserAnswers) > len(post.PollData.ProvidedAnswers) {
		return nil, sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("user's answers are more than the available ones in Poll"),
		)
	}

	for _, answer := range msg.UserAnswers {
		if found := answerExist(post.PollData.ProvidedAnswers.ExtractAnswersIDs(), answer); !found {
			return nil, sdkerrors.Wrap(
				sdkerrors.ErrInvalidRequest,
				fmt.Sprintf(
					"answer with ID %s isn't one of the poll's provided answers",
					strconv.FormatUint(uint64(answer), 10)),
			)
		}
	}

	pollAnswers := keeper.GetPollAnswersByUser(ctx, post.PostID, msg.Answerer)

	// check if the poll allows to edit previous answers
	if len(pollAnswers) > 0 && !post.PollData.AllowsAnswerEdits {
		return nil, sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("post with ID %s doesn't allow answers' edits", post.PostID),
		)
	}

	userPollAnswers := types.NewUserAnswer(msg.UserAnswers, msg.Answerer)

	keeper.SavePollAnswers(ctx, post.PostID, userPollAnswers)

	answerEvent := sdk.NewEvent(
		types.EventTypeAnsweredPoll,
		sdk.NewAttribute(types.AttributeKeyPostID, msg.PostID.String()),
		sdk.NewAttribute(types.AttributeKeyPollAnswerer, msg.Answerer.String()),
	)

	ctx.EventManager().EmitEvent(answerEvent)

	result := sdk.Result{
		Data:   keeper.Cdc.MustMarshalBinaryLengthPrefixed("Answered to poll correctly"),
		Events: sdk.Events{answerEvent}.ToABCIEvents(),
	}
	return &result, nil
}

// checkPostPollValid performs all the checks to ensure the post with the given id exists, contains a poll and such poll has not been closed
func checkPostPollValid(ctx sdk.Context, id types.PostID, keeper Keeper) (*types.Post, error) {
	// checks if post exists
	post, found := keeper.GetPost(ctx, id)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("post with id %s doesn't exist", id))
	}

	// checks if post has a poll
	if post.PollData == nil {
		return &post, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("no poll associated with ID: %s", id))
	}

	// checks if the poll is already closed or not
	if !post.PollData.Open {
		return &post, sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("the poll associated with ID %s was closed at %s", post.PostID, post.PollData.EndDate),
		)
	}

	return &post, nil
}

// answerExist checks if the answer is contained in providedAnswers slice
func answerExist(providedAnswers []types.AnswerID, answer types.AnswerID) bool {
	for _, ans := range providedAnswers {
		if ans == answer {
			return true
		}
	}
	return false
}
