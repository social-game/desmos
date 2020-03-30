package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/codec"
	sim "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/desmos-labs/desmos/app/params"
	"github.com/desmos-labs/desmos/x/posts/internal/keeper"
)

// Simulation operation weights constants
const (
	OpWeightMsgCreatePost       = "op_weight_msg_create_post"
	OpWeightMsgEditPost         = "op_weight_msg_edit_post"
	OpWeightMsgAddReaction      = "op_weight_msg_add_reaction"
	OpWeightMsgRemoveReaction   = "op_weight_msg_remove_reaction"
	OpWeightMsgAnswerPoll       = "op_weight_msg_answer_poll"
	OpWeightMsgRegisterReaction = "op_weight_msg_register_reaction"

	DefaultGasValue = 5000000
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams sim.AppParams, cdc *codec.Codec, k keeper.Keeper, ak auth.AccountKeeper, bk bank.Keeper,
) simulation.WeightedOperations {

	var weightMsgCreatePost int
	appParams.GetOrGenerate(cdc, OpWeightMsgCreatePost, &weightMsgCreatePost, nil,
		func(_ *rand.Rand) {
			weightMsgCreatePost = params.DefaultWeightMsgCreatePost
		},
	)

	var weightMsgEditPost int
	appParams.GetOrGenerate(cdc, OpWeightMsgEditPost, &weightMsgEditPost, nil,
		func(_ *rand.Rand) {
			weightMsgEditPost = params.DefaultWeightMsgEditPost
		},
	)

	var weightMsgAddReaction int
	appParams.GetOrGenerate(cdc, OpWeightMsgAddReaction, &weightMsgAddReaction, nil,
		func(_ *rand.Rand) {
			weightMsgAddReaction = params.DefaultWeightMsgAddReaction
		},
	)

	var weightMsgRemoveReaction int
	appParams.GetOrGenerate(cdc, OpWeightMsgRemoveReaction, &weightMsgRemoveReaction, nil,
		func(_ *rand.Rand) {
			weightMsgRemoveReaction = params.DefaultWeightMsgRemoveReaction
		},
	)

	var weightMsgAnswerPoll int
	appParams.GetOrGenerate(cdc, OpWeightMsgAnswerPoll, &weightMsgAnswerPoll, nil,
		func(_ *rand.Rand) {
			weightMsgAnswerPoll = params.DefaultWeightMsgAnswerPoll
		},
	)

	var weightMsgRegisterReaction int
	appParams.GetOrGenerate(cdc, OpWeightMsgRegisterReaction, &weightMsgRegisterReaction, nil,
		func(_ *rand.Rand) {
			weightMsgRegisterReaction = params.DefaultWeightMsgRegisterReaction
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgCreatePost,
			SimulateMsgCreatePost(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightMsgEditPost,
			SimulateMsgEditPost(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightMsgRegisterReaction,
			SimulateMsgRegisterReaction(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightMsgAddReaction,
			SimulateMsgAddPostReaction(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightMsgRemoveReaction,
			SimulateMsgRemovePostReaction(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightMsgAnswerPoll,
			SimulateMsgAnswerToPoll(k, ak, bk),
		),
	}
}
