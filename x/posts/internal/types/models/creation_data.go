package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PostCreationData contains the data that can be sent while creating a new post
type PostCreationData struct {
	ParentID       PostID       `json:"parent_id"`
	Message        string       `json:"message"`
	AllowsComments bool         `json:"allows_comments"`
	Subspace       string       `json:"subspace"`
	OptionalData   OptionalData `json:"optional_data,omitempty"`
	Creator        string       `json:"creator"`
	CreationDate   time.Time    `json:"creation_date"`
	Medias         PostMedias   `json:"medias,omitempty"`
	PollData       *PollData    `json:"poll_data,omitempty"`
}

// NewPostCreationData is a constructor function for PostCreationData
func NewPostCreationData(message string, parentID PostID, allowsComments bool, subspace string,
	optionalData map[string]string, owner string, creationDate time.Time,
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
	if len(strings.TrimSpace(data.Message)) == 0 && len(data.Medias) == 0 {
		return fmt.Errorf("post message or medias are required and cannot be both blank or empty")
	}

	if len(data.Message) > MaxPostMessageLength {
		return fmt.Errorf("post message cannot exceed %dataType characters", MaxPostMessageLength)
	}

	if !SubspaceRegEx.MatchString(data.Subspace) {
		return fmt.Errorf("post subspace must be a valid sha-256 hash")
	}

	if len(data.OptionalData) > MaxOptionalDataFieldsNumber {
		return fmt.Errorf("post optional data cannot be longer than %dataType fields", MaxOptionalDataFieldsNumber)
	}

	for key, value := range data.OptionalData {
		if len(value) > MaxOptionalDataFieldValueLength {
			return fmt.Errorf("post optional data value lengths cannot be longer than %dataType. %s exceeds the limit",
				MaxOptionalDataFieldValueLength, key)
		}
	}

	if len(strings.TrimSpace(data.Creator)) == 0 {
		return fmt.Errorf("invalid creator address: %s", data.Creator)
	}

	if data.CreationDate.IsZero() {
		return fmt.Errorf("invalid post creation date")
	}

	if data.CreationDate.After(time.Now().UTC()) {
		return fmt.Errorf("creation date cannot be in the future")
	}

	if data.Medias != nil {
		if err := data.Medias.Validate(); err != nil {
			return fmt.Errorf(err.Error())
		}
	}

	if data.PollData != nil {
		if !data.PollData.Open {
			return fmt.Errorf("poll Post cannot be created closed")
		}
		if err := data.PollData.Validate(); err != nil {
			return fmt.Errorf(err.Error())
		}
	}

	return nil
}

// MarshalJSON implements json.Marshaler
// This is done due to the fact that Amino does not handle properly omitempty clauses
func (data PostCreationData) MarshalJSON() ([]byte, error) {
	type tmp PostCreationData
	return json.Marshal(tmp(data))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// This is done due to the fact that Amino does not handle properly omitempty clauses
func (data *PostCreationData) UnmarshalJSON(bytes []byte) error {
	type dataType PostCreationData
	var temp dataType
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return err
	}

	*data = PostCreationData(temp)
	return nil
}

// GetCreatorAddress returns the address of the creator as an sdk.AccAddress
func (data PostCreationData) GetCreatorAddress() (sdk.AccAddress, error) {
	address, err := sdk.AccAddressFromBech32(data.Creator)
	if err != nil {
		return nil, err
	}
	return address, nil
}
