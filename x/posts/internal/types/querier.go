package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/x/posts/internal/types/models"
)

// QueryPostsParams Params for query 'custom/posts/posts'
type QueryPostsParams struct {
	Page  int
	Limit int

	SortBy    string // Field that should determine the sorting
	SortOrder string // Either ascending or descending

	ParentID       *models.PostID
	CreationTime   *time.Time
	AllowsComments *bool
	Subspace       string
	Creator        sdk.AccAddress
	Hashtags       []string
}

func DefaultQueryPostsParams(page, limit int) QueryPostsParams {
	return QueryPostsParams{
		Page:  page,
		Limit: limit,

		SortBy:    models.PostSortByID,
		SortOrder: models.PostSortOrderAscending,

		ParentID:       nil,
		CreationTime:   nil,
		AllowsComments: nil,
		Subspace:       "",
		Creator:        nil,
		Hashtags:       nil,
	}
}
