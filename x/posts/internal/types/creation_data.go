package types

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// PostCreationData contains the data that can be sent while creating a new post
type PostCreationData struct {
	ParentID       PostID         `json:"parent_id"`
	Message        string         `json:"message"`
	AllowsComments bool           `json:"allows_comments"`
	Subspace       string         `json:"subspace"`
	OptionalData   OptionalData   `json:"optional_data,omitempty"`
	Creator        sdk.AccAddress `json:"creator"`
	CreationDate   time.Time      `json:"creation_date"`
	Medias         PostMedias     `json:"medias,omitempty"`
	PollData       *PollData      `json:"poll_data,omitempty"`
}

// NewPostCreationData is a constructor function for PostCreationData
func NewPostCreationData(message string, parentID PostID, allowsComments bool, subspace string,
	optionalData map[string]string, owner sdk.AccAddress, creationDate time.Time,
	medias PostMedias, pollData *PollData) PostCreationData {

	return PostCreationData{
		Message:        message,
		ParentID:       parentID,
		AllowsComments: allowsComments,
		Subspace:       subspace,
		OptionalData:   optionalData,
		Creator:        owner,
		CreationDate:   creationDate,
		Medias:         medias,
		PollData:       pollData,
	}
}

// String returns a string representation of FungibleTokenPacketData
func (data PostCreationData) String() string {
	return fmt.Sprintf(`FungibleTokenPacketData:
	ParentId:           %s
	Message:            %s
	AllowsComments:     %t
	Subspace:           %s
    OptionalData:       %s 
    Creator:            %s
    CreationDate:       %s
    Medias:             %s
    PollData:           %s`,
		data.ParentID,
		data.Message,
		data.AllowsComments,
		data.Subspace,
		data.OptionalData,
		data.Creator,
		data.CreationDate,
		data.Medias,
		data.PollData,
	)
}

// ValidateBasic runs stateless checks on the creation data
func (data PostCreationData) ValidateBasic() error {
	if data.Creator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("Invalid creator address: %s", data.Creator))
	}

	if len(strings.TrimSpace(data.Message)) == 0 && len(data.Medias) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Post message or medias are required and cannot be both blank or empty")
	}

	if len(data.Message) > MaxPostMessageLength {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("Post message cannot exceed %d characters", MaxPostMessageLength))
	}

	if !SubspaceRegEx.MatchString(data.Subspace) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Post subspace must be a valid sha-256 hash")
	}

	if len(data.OptionalData) > MaxOptionalDataFieldsNumber {
		data := fmt.Sprintf("Post optional data cannot be longer than %d fields", MaxOptionalDataFieldsNumber)
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, data)
	}

	for key, value := range data.OptionalData {
		if len(value) > MaxOptionalDataFieldValueLength {
			data := fmt.Sprintf("Post optional data value lengths cannot be longer than %d. %s exceeds the limit",
				MaxOptionalDataFieldValueLength, key)
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, data)
		}
	}

	if data.CreationDate.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid post creation date")
	}

	if data.CreationDate.After(time.Now().UTC()) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Creation date cannot be in the future")
	}

	if data.Medias != nil {
		if err := data.Medias.Validate(); err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
		}
	}

	if data.PollData != nil {
		if !data.PollData.Open {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Poll Post cannot be created closed")
		}
		if err := data.PollData.Validate(); err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
		}
	}

	return nil
}
