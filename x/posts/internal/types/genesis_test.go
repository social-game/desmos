package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/x/posts/internal/types/models"
	"github.com/desmos-labs/desmos/x/posts/internal/types/models/polls"
	"github.com/desmos-labs/desmos/x/posts/internal/types/models/reactions"

	"github.com/desmos-labs/desmos/x/posts/internal/types"
	"github.com/stretchr/testify/require"
)

func TestValidateGenesis(t *testing.T) {
	user, err := sdk.AccAddressFromBech32("cosmos1s3nh6tafl4amaxkke9kdejhp09lk93g9ev39r4")
	require.NoError(t, err)

	tests := []struct {
		name        string
		genesis     types.GenesisState
		shouldError bool
	}{
		{
			name:        "DefaultGenesis does not error",
			genesis:     types.DefaultGenesisState(),
			shouldError: false,
		},
		{
			name: "Genesis with invalid post errors",
			genesis: types.GenesisState{
				Posts:         models.Posts{models.Post{PostID: models.PostID(0)}},
				PostReactions: map[string]reactions.PostReactions{},
			},
			shouldError: true,
		},
		{
			name: "Genesis with invalid post reaction errors",
			genesis: types.GenesisState{
				Posts: models.Posts{},
				PostReactions: map[string]reactions.PostReactions{
					"1": {reactions.PostReaction{Owner: nil}},
				},
			},
			shouldError: true,
		},
		{
			name: "Genesis with invalid poll answers errors",
			genesis: types.GenesisState{
				Posts: models.Posts{},
				PollAnswers: map[string]polls.UserAnswers{
					"1": {
						polls.NewUserAnswer([]polls.AnswerID{}, user),
					},
				},
				PostReactions: map[string]reactions.PostReactions{},
			},
			shouldError: true,
		},
		{
			name: "Genesis with invalid registered reaction errors",
			genesis: types.GenesisState{
				Posts: models.Posts{},
				RegisteredReactions: reactions.Reactions{reactions.NewReaction(user, ":smile", "smile.jpg",
					"4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e")},
			},
			shouldError: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			if test.shouldError {
				require.Error(t, types.ValidateGenesis(test.genesis))
			} else {
				require.NoError(t, types.ValidateGenesis(test.genesis))
			}
		})
	}
}
