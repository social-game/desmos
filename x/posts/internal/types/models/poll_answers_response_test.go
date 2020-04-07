package models_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/x/posts/internal/types/models"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestPollAnswersQueryResponse_String(t *testing.T) {
	testOwner, err := sdk.AccAddressFromBech32("cosmos1cjf97gpzwmaf30pzvaargfgr884mpp5ak8f7ns")
	require.NoError(t, err)

	answers := models.UserAnswer{
		Answers: []models.AnswerID{1, 2},
		User:    testOwner,
	}

	pollResponse := models.PollAnswersQueryResponse{
		PostID:         models.PostID(0),
		AnswersDetails: models.UserAnswers{answers},
	}

	require.Equal(t, "Post ID [0] - Answers Details:\nUser: cosmos1cjf97gpzwmaf30pzvaargfgr884mpp5ak8f7ns \nAnswers IDs: 1 2", pollResponse.String())
}
