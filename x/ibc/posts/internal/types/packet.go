package types

import (
	"encoding/json"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/x/posts"
)

// PostCreationPacketData represents the IBC packet data that needs to be sent when
// creating a Desmos post
type PostCreationPacketData struct {
	posts.PostCreationData
}

// NewPostCreationData is a constructor function for PostCreationPacketData
func NewPostCreationPacketData(
	message string, parentID uint64, allowsComments bool, subspace string,
	optionalData map[string]string, owner string, creationDate time.Time,
	medias posts.PostMedias, pollData *posts.PollData,
) PostCreationPacketData {
	return PostCreationPacketData{
		PostCreationData: posts.NewPostCreationData(
			message,
			posts.PostID(parentID),
			allowsComments,
			subspace,
			optionalData,
			owner,
			creationDate,
			medias,
			pollData,
		),
	}
}

// MarshalJSON implements json.Marshaler
func (data PostCreationPacketData) MarshalJSON() ([]byte, error) {
	return data.PostCreationData.MarshalJSON()
}

func (data *PostCreationPacketData) UnmarshalJSON(bytes []byte) error {
	var postData posts.PostCreationData
	err := json.Unmarshal(bytes, &postData)
	if err != nil {
		return err
	}
	*data = PostCreationPacketData{postData}
	return nil
}

// GetBytes allows to use this inside an IBC packet
func (data PostCreationPacketData) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(data))
}

// PostCreationPacketAcknowledgement contains a boolean success flag and an optional error msg
// error msg is empty string on success
// See spec for onAcknowledgePacket: https://github.com/cosmos/ics/tree/master/spec/ics-020-fungible-token-transfer#packet-relay
type PostCreationPacketAcknowledgement struct {
	Success bool   `json:"success" yaml:"success"`
	Error   string `json:"error" yaml:"error"`
}

// GetBytes is a helper for serialising
func (ack PostCreationPacketAcknowledgement) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(ack))
}
