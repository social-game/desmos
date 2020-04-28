package models

// nolint
// autogenerated code using github.com/haasted/alias-generator.
// based on functionality in github.com/rigelrozanski/multitool

import (
	"github.com/desmos-labs/desmos/x/posts/internal/types/models/common"
	"github.com/desmos-labs/desmos/x/posts/internal/types/models/polls"
	"github.com/desmos-labs/desmos/x/posts/internal/types/models/reactions"
)

const (
	ModuleName                      = common.ModuleName
	RouterKey                       = common.RouterKey
	StoreKey                        = common.StoreKey
	MaxPostMessageLength            = common.MaxPostMessageLength
	MaxOptionalDataFieldsNumber     = common.MaxOptionalDataFieldsNumber
	MaxOptionalDataFieldValueLength = common.MaxOptionalDataFieldValueLength
	ActionCreatePost                = common.ActionCreatePost
	ActionEditPost                  = common.ActionEditPost
	ActionAnswerPoll                = common.ActionAnswerPoll
	ActionAddPostReaction           = common.ActionAddPostReaction
	ActionRemovePostReaction        = common.ActionRemovePostReaction
	ActionRegisterReaction          = common.ActionRegisterReaction
	QuerierRoute                    = common.QuerierRoute
	QueryPost                       = common.QueryPost
	QueryPosts                      = common.QueryPosts
	QueryPollAnswers                = common.QueryPollAnswers
	QueryRegisteredReactions        = common.QueryRegisteredReactions
	PostSortByID                    = common.PostSortByID
	PostSortByCreationDate          = common.PostSortByCreationDate
	PostSortOrderAscending          = common.PostSortOrderAscending
	PostSortOrderDescending         = common.PostSortOrderDescending
)

var (
	// functions aliases
	NewPostMedia      = common.NewPostMedia
	ValidateURI       = common.ValidateURI
	NewPostMedias     = common.NewPostMedias
	NewPollData       = polls.NewPollData
	ArePollDataEquals = polls.ArePollDataEquals
	NewUserAnswer     = polls.NewUserAnswer
	ParseAnswerID     = polls.ParseAnswerID
	NewPollAnswer     = polls.NewPollAnswer
	NewPollAnswers    = polls.NewPollAnswers
	NewReaction       = reactions.NewReaction
	IsEmoji           = reactions.IsEmoji
	NewPostReaction   = reactions.NewPostReaction

	// variable aliases
	SubspaceRegEx            = common.SubspaceRegEx
	HashtagRegEx             = common.HashtagRegEx
	ShortCodeRegEx           = common.ShortCodeRegEx
	URIRegEx                 = common.URIRegEx
	LastPostIDStoreKey       = common.LastPostIDStoreKey
	PostStorePrefix          = common.PostStorePrefix
	PostCommentsStorePrefix  = common.PostCommentsStorePrefix
	PostReactionsStorePrefix = common.PostReactionsStorePrefix
	ReactionsStorePrefix     = common.ReactionsStorePrefix
	PollAnswersStorePrefix   = common.PollAnswersStorePrefix
)

type (
	PostMedia     = common.PostMedia
	PostMedias    = common.PostMedias
	OptionalData  = common.OptionalData
	KeyValue      = common.KeyValue
	PollData      = polls.PollData
	UserAnswer    = polls.UserAnswer
	UserAnswers   = polls.UserAnswers
	AnswerID      = polls.AnswerID
	PollAnswer    = polls.PollAnswer
	PollAnswers   = polls.PollAnswers
	Reaction      = reactions.Reaction
	Reactions     = reactions.Reactions
	PostReaction  = reactions.PostReaction
	PostReactions = reactions.PostReactions
)
