package msgs

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/desmos-labs/desmos/x/posts/internal/types/models"
)

// ----------------------
// --- MsgCreatePost
// ----------------------

// MsgCreatePost defines a CreatePost message
type MsgCreatePost struct {
	models.PostCreationData
}

// NewMsgCreatePost is a constructor function for MsgCreatePost
func NewMsgCreatePost(message string, parentID models.PostID, allowsComments bool, subspace string,
	optionalData map[string]string, owner sdk.AccAddress, creationDate time.Time,
	medias models.PostMedias, pollData *models.PollData) MsgCreatePost {
	return MsgCreatePost{
		PostCreationData: models.NewPostCreationData(
			message, parentID, allowsComments, subspace, optionalData, owner, creationDate, medias, pollData,
		),
	}
}

// Route should return the name of the module
func (msg MsgCreatePost) Route() string { return models.RouterKey }

// Type should return the action
func (msg MsgCreatePost) Type() string { return models.ActionCreatePost }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreatePost) ValidateBasic() error {
	return msg.PostCreationData.ValidateBasic()
}

// GetSignBytes encodes the message for signing
func (msg MsgCreatePost) GetSignBytes() []byte {
	return sdk.MustSortJSON(MsgsCodec.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreatePost) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

// MarshalJSON implements the json.Marshaler interface.
// This is done due to the fact that Amino does not respect omitempty clauses
func (msg MsgCreatePost) MarshalJSON() ([]byte, error) {
	type temp MsgCreatePost
	return json.Marshal(temp(msg))
}

// ----------------------
// --- MsgEditPost
// ----------------------

// MsgEditPost defines the EditPostMessage message
type MsgEditPost struct {
	PostID   models.PostID  `json:"post_id"`
	Message  string         `json:"message"`
	Editor   sdk.AccAddress `json:"editor"`
	EditDate time.Time      `json:"edit_date"`
}

// NewMsgEditPost is the constructor function for MsgEditPost
func NewMsgEditPost(id models.PostID, message string, owner sdk.AccAddress, editDate time.Time) MsgEditPost {
	return MsgEditPost{
		PostID:   id,
		Message:  message,
		Editor:   owner,
		EditDate: editDate,
	}
}

// Route should return the name of the module
func (msg MsgEditPost) Route() string { return models.RouterKey }

// Type should return the action
func (msg MsgEditPost) Type() string { return models.ActionEditPost }

// ValidateBasic runs stateless checks on the message
func (msg MsgEditPost) ValidateBasic() error {
	if !msg.PostID.Valid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid post id")
	}

	if msg.Editor.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("Invalid editor address: %s", msg.Editor))
	}

	if len(strings.TrimSpace(msg.Message)) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Post message cannot be empty nor blank")
	}

	if msg.EditDate.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid edit date")
	}

	if msg.EditDate.After(time.Now().UTC()) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Edit date cannot be in the future")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgEditPost) GetSignBytes() []byte {
	return sdk.MustSortJSON(MsgsCodec.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgEditPost) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Editor}
}
