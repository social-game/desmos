package models_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/x/posts/internal/types/models"
	"github.com/stretchr/testify/require"
)

// -------------
// --- PostID
// -------------

func TestPostID_Next(t *testing.T) {
	tests := []struct {
		id       models.PostID
		expected models.PostID
	}{
		{
			id:       models.PostID(0),
			expected: models.PostID(1),
		},
		{
			id:       models.PostID(1123123),
			expected: models.PostID(1123124),
		},
	}

	for index, test := range tests {
		test := test
		t.Run(fmt.Sprintf("Test index: %dataType", index), func(t *testing.T) {
			require.Equal(t, test.expected, test.id.Next())
		})
	}
}

func TestPostID_MarshalJSON(t *testing.T) {
	bz, err := json.Marshal(models.PostID(10))
	require.NoError(t, err)
	require.Equal(t, `"10"`, string(bz))
}

func TestPostID_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expID    models.PostID
		expError string
	}{
		{
			name:     "Invalid ID returns error",
			value:    "id",
			expID:    models.PostID(0),
			expError: "invalid character 'i' looking for beginning of value",
		},
		{
			name:     "Valid id is read properly",
			value:    `"10"`,
			expID:    models.PostID(10),
			expError: "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			var id models.PostID
			err := json.Unmarshal([]byte(test.value), &id)

			if err == nil {
				require.Equal(t, test.expID, id)
			} else {
				require.Equal(t, test.expError, err.Error())
			}
		})
	}
}

func TestParsePostID(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expID    models.PostID
		expError string
	}{
		{
			name:     "Invalid id returns error",
			value:    "id",
			expID:    models.PostID(0),
			expError: "strconv.ParseUint: parsing \"id\": invalid syntax",
		},
		{
			name:     "Valid id returns proper value",
			value:    "10",
			expID:    models.PostID(10),
			expError: "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			id, err := models.ParsePostID(test.value)

			if err == nil {
				require.Equal(t, test.expID, id)
			} else {
				require.Equal(t, test.expError, err.Error())
			}
		})
	}
}

// -------------
// --- PostIDs
// -------------

func TestPostIDs_Equals(t *testing.T) {
	tests := []struct {
		name      string
		first     models.PostIDs
		second    models.PostIDs
		expEquals bool
	}{
		{
			name:      "Different length",
			first:     models.PostIDs{models.PostID(1), models.PostID(0)},
			second:    models.PostIDs{models.PostID(1)},
			expEquals: false,
		},
		{
			name:      "Different order",
			first:     models.PostIDs{models.PostID(0), models.PostID(1)},
			second:    models.PostIDs{models.PostID(1), models.PostID(0)},
			expEquals: false,
		},
		{
			name:      "Same length and order",
			first:     models.PostIDs{models.PostID(0), models.PostID(1)},
			second:    models.PostIDs{models.PostID(0), models.PostID(1)},
			expEquals: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expEquals, test.first.Equals(test.second))
		})
	}
}

func TestPostIDs_AppendIfMissing(t *testing.T) {
	tests := []struct {
		name      string
		IDs       models.PostIDs
		newID     models.PostID
		expIDs    models.PostIDs
		expEdited bool
	}{
		{
			name:      "AppendIfMissing dont append anything",
			IDs:       models.PostIDs{models.PostID(1)},
			newID:     models.PostID(1),
			expIDs:    models.PostIDs{models.PostID(1)},
			expEdited: false,
		},
		{
			name:      "AppendIfMissing append something",
			IDs:       models.PostIDs{models.PostID(1)},
			newID:     models.PostID(2),
			expIDs:    models.PostIDs{models.PostID(1), models.PostID(2)},
			expEdited: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			newIDs, edited := test.IDs.AppendIfMissing(test.newID)
			require.Equal(t, test.expIDs, newIDs)
			require.Equal(t, test.expEdited, edited)
		})
	}
}

// -----------
// --- Post
// -----------

func TestPost_String(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32("cosmos1cjf97gpzwmaf30pzvaargfgr884mpp5ak8f7ns")
	require.NoError(t, err)

	timeZone, err := time.LoadLocation("UTC")
	require.NoError(t, err)

	date := time.Date(2020, 1, 1, 12, 00, 00, 000, timeZone)
	post := models.Post{
		PostID:         models.PostID(19),
		ParentID:       models.PostID(1),
		Message:        "My post message",
		Created:        date,
		LastEdited:     date.AddDate(0, 0, 1),
		AllowsComments: true,
		Subspace:       "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
		OptionalData:   map[string]string{},
		Creator:        owner,
	}

	require.Equal(t,
		`{"id":"19","parent_id":"1","message":"My post message","created":"2020-01-01T12:00:00Z","last_edited":"2020-01-02T12:00:00Z","allows_comments":true,"subspace":"4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e","creator":"cosmos1cjf97gpzwmaf30pzvaargfgr884mpp5ak8f7ns"}`,
		post.String(),
	)
}

func TestPost_Validate(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32("cosmos1cjf97gpzwmaf30pzvaargfgr884mpp5ak8f7ns")
	require.NoError(t, err)

	timeZone, err := time.LoadLocation("UTC")
	require.NoError(t, err)

	date := time.Date(2020, 1, 1, 12, 00, 00, 000, timeZone)
	medias := models.PostMedias{
		models.PostMedia{
			URI:      "https://uri.com",
			MimeType: "text/plain",
		},
	}
	answer := models.PollAnswer{
		ID:   models.AnswerID(1),
		Text: "Yes",
	}

	answer2 := models.PollAnswer{
		ID:   models.AnswerID(2),
		Text: "No",
	}
	pollData := models.NewPollData("poll?", time.Now().UTC().Add(time.Hour), models.PollAnswers{answer, answer2}, true, false, true)

	tests := []struct {
		post     models.Post
		expError string
	}{
		{
			post:     models.NewPost(models.PostID(0), models.PostID(0), "Message", true, "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e", map[string]string{}, date, owner).WithMedias(medias).WithPollData(pollData),
			expError: "invalid post id: 0",
		},
		{
			post:     models.NewPost(models.PostID(1), models.PostID(0), "", true, "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e", map[string]string{}, date, nil).WithMedias(medias).WithPollData(pollData),
			expError: "invalid post owner: ",
		},
		{
			post:     models.NewPost(models.PostID(1), models.PostID(0), "", true, "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e", map[string]string{}, date, owner).WithPollData(pollData),
			expError: "post message or medias required, they cannot be both empty",
		},
		{
			post:     models.NewPost(models.PostID(1), models.PostID(0), " ", true, "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e", map[string]string{}, date, owner).WithPollData(pollData),
			expError: "post message or medias required, they cannot be both empty",
		},
		{
			post:     models.NewPost(models.PostID(1), models.PostID(0), "Message", true, "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e", map[string]string{}, time.Time{}, owner).WithMedias(medias).WithPollData(pollData),
			expError: "invalid post creation time: 0001-01-01 00:00:00 +0000 UTC",
		},
		{
			post:     models.Post{PostID: models.PostID(19), Creator: owner, Message: "Message", Subspace: "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e", Created: date, LastEdited: date.AddDate(0, 0, -1)},
			expError: "invalid post last edit time: 2019-12-31 12:00:00 +0000 UTC",
		},
		{
			post:     models.NewPost(models.PostID(1), models.PostID(0), "Message", true, "", map[string]string{}, date, owner).WithMedias(medias).WithPollData(pollData),
			expError: "post subspace must be a valid sha-256 hash",
		},
		{
			post:     models.NewPost(models.PostID(1), models.PostID(0), "Message", true, " ", map[string]string{}, date, owner).WithMedias(medias).WithPollData(pollData),
			expError: "post subspace must be a valid sha-256 hash",
		},
		{
			post: models.Post{
				PostID:         models.PostID(1),
				ParentID:       models.PostID(0),
				Message:        "Message",
				AllowsComments: true,
				Subspace:       "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				OptionalData:   map[string]string{},
				Created:        time.Now().UTC().Add(time.Hour),
				Creator:        owner,
				Medias:         medias,
			},
			expError: "post creation date cannot be in the future",
		},
		{
			post: models.Post{
				PostID:         models.PostID(1),
				ParentID:       models.PostID(0),
				Message:        "Message",
				AllowsComments: true,
				Subspace:       "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				OptionalData:   map[string]string{},
				Created:        time.Now().UTC(),
				LastEdited:     time.Now().UTC().Add(time.Hour),
				Creator:        owner,
				Medias:         medias,
			},
			expError: "post last edit date cannot be in the future",
		},
		{
			post: models.NewPost(
				models.PostID(1),
				models.PostID(0),
				`
				Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque massa felis, aliquam sed ipsum at, 
				mollis pharetra quam. Vestibulum nec nulla ante. Praesent sed dignissim turpis. Curabitur aliquam nunc 
				eu nisi porta, eu gravida purus faucibus. Duis commodo sagittis lacus, vitae luctus enim vulputate a. 
				Nulla tempor eget nunc vitae vulputate. Nulla facilities. Donec sollicitudin odio in arcu efficitur, 
				sit amet vestibulum diam ullamcorper. Ut ac dolor in velit gravida efficitur et et erat volutpat.
				`,
				true,
				"4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				map[string]string{},
				date,
				owner,
			).WithMedias(medias).WithPollData(pollData),
			expError: "post message cannot be longer than 500 characters",
		},
		{
			post: models.NewPost(
				models.PostID(1),
				models.PostID(0),
				"Message",
				true,
				"4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				map[string]string{
					"key1":  "value",
					"key2":  "value",
					"key3":  "value",
					"key4":  "value",
					"key5":  "value",
					"key6":  "value",
					"key7":  "value",
					"key8":  "value",
					"key9":  "value",
					"key10": "value",
					"key11": "value",
				},
				date,
				owner,
			).WithMedias(medias).WithPollData(pollData),
			expError: "post optional data cannot contain more than 10 key-value pairs",
		},
		{
			post: models.NewPost(
				models.PostID(1),
				models.PostID(0),
				"Message",
				true,
				"4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				map[string]string{
					"key1": `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque euismod, mi at commodo 
							efficitur, quam sapien congue enim, ut porttitor lacus tellus vitae turpis. Vivamus aliquam 
							sem eget neque metus.`,
				},
				date,
				owner,
			).WithMedias(medias).WithPollData(pollData),
			expError: "post optional data values cannot exceed 200 characters. key1 of post with id 1 is longer than this",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.expError, func(t *testing.T) {
			if len(test.expError) != 0 {
				require.Equal(t, test.expError, test.post.Validate().Error())
			} else {
				require.Nil(t, test.post.Validate())
			}
		})
	}
}

func TestPost_Equals(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32("cosmos1cjf97gpzwmaf30pzvaargfgr884mpp5ak8f7ns")
	require.NoError(t, err)

	otherOwner, err := sdk.AccAddressFromBech32("cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47")
	require.NoError(t, err)

	timeZone, err := time.LoadLocation("UTC")
	require.NoError(t, err)

	date := time.Date(2020, 1, 1, 12, 00, 00, 000, timeZone)
	medias := models.PostMedias{
		models.PostMedia{
			URI:      "https://uri.com",
			MimeType: "text/plain",
		},
	}

	pollData := models.NewPollData(
		"poll?",
		time.Date(2050, 1, 1, 15, 15, 00, 000, timeZone),
		models.NewPollAnswers(
			models.NewPollAnswer(models.AnswerID(1), "Yes"),
			models.NewPollAnswer(models.AnswerID(2), "No"),
		),
		true,
		false,
		true,
	)

	tests := []struct {
		name      string
		first     models.Post
		second    models.Post
		expEquals bool
	}{
		{
			name: "Different post ID",
			first: models.Post{
				PostID:         models.PostID(19),
				ParentID:       models.PostID(1),
				Message:        "My post message",
				Created:        date,
				LastEdited:     date.AddDate(0, 0, 1),
				AllowsComments: true,
				Subspace:       "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				OptionalData:   map[string]string{},
				Creator:        owner,
				Medias:         medias,
			},
			second: models.Post{
				PostID:         models.PostID(10),
				ParentID:       models.PostID(1),
				Message:        "My post message",
				Created:        date,
				LastEdited:     date.AddDate(0, 0, 1),
				AllowsComments: true,
				Subspace:       "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				OptionalData:   map[string]string{},
				Creator:        owner,
				Medias:         medias,
			},
			expEquals: false,
		},
		{
			name: "Different parent ID",
			first: models.Post{
				PostID:         models.PostID(19),
				ParentID:       models.PostID(1),
				Message:        "My post message",
				Created:        date,
				LastEdited:     date.AddDate(0, 0, 1),
				AllowsComments: true,
				Subspace:       "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				OptionalData:   map[string]string{},
				Creator:        owner,
				Medias:         medias,
			},
			second: models.Post{
				PostID:         models.PostID(19),
				ParentID:       models.PostID(10),
				Message:        "My post message",
				Created:        date,
				LastEdited:     date.AddDate(0, 0, 1),
				AllowsComments: true,
				Subspace:       "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				OptionalData:   map[string]string{},
				Creator:        owner,
				Medias:         medias,
			},
			expEquals: false,
		},
		{
			name: "Different message",
			first: models.Post{
				PostID:         models.PostID(19),
				ParentID:       models.PostID(1),
				Message:        "My post message",
				Created:        date,
				LastEdited:     date.AddDate(0, 0, 1),
				AllowsComments: true,
				Subspace:       "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				OptionalData:   map[string]string{},
				Creator:        owner,
				Medias:         medias,
			},
			second: models.Post{
				PostID:         models.PostID(19),
				ParentID:       models.PostID(1),
				Message:        "Another post message",
				Created:        date,
				LastEdited:     date.AddDate(0, 0, 1),
				AllowsComments: true,
				Subspace:       "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				OptionalData:   map[string]string{},
				Creator:        owner,
				Medias:         medias,
			},
			expEquals: false,
		},
		{
			name: "Different creation time",
			first: models.Post{
				PostID:         models.PostID(19),
				ParentID:       models.PostID(1),
				Message:        "My post message",
				Created:        date,
				LastEdited:     date.AddDate(0, 0, 1),
				AllowsComments: true,
				Subspace:       "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				OptionalData:   map[string]string{},
				Creator:        owner,
				Medias:         medias,
			},
			second: models.Post{
				PostID:         models.PostID(19),
				ParentID:       models.PostID(1),
				Message:        "My post message",
				Created:        date.AddDate(0, 0, 1),
				LastEdited:     date.AddDate(0, 0, 1),
				AllowsComments: true,
				Subspace:       "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				OptionalData:   map[string]string{},
				Creator:        owner,
				Medias:         medias,
			},
			expEquals: false,
		},
		{
			name: "Different last edited",
			first: models.Post{
				PostID:         models.PostID(19),
				ParentID:       models.PostID(1),
				Message:        "My post message",
				Created:        date,
				LastEdited:     date.AddDate(0, 0, 1),
				AllowsComments: true,
				Subspace:       "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				OptionalData:   map[string]string{},
				Creator:        owner,
				Medias:         medias,
			},
			second: models.Post{
				PostID:         models.PostID(19),
				ParentID:       models.PostID(1),
				Message:        "My post message",
				Created:        date,
				LastEdited:     date.AddDate(0, 0, 2),
				AllowsComments: true,
				Subspace:       "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				OptionalData:   map[string]string{},
				Creator:        owner,
				Medias:         medias,
			},
			expEquals: false,
		},
		{
			name: "Different allows comments",
			first: models.Post{
				PostID:         models.PostID(19),
				ParentID:       models.PostID(1),
				Message:        "My post message",
				Created:        date,
				LastEdited:     date.AddDate(0, 0, 1),
				AllowsComments: true,
				Subspace:       "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				OptionalData:   map[string]string{},
				Creator:        owner,
				Medias:         medias,
			},
			second: models.Post{
				PostID:         models.PostID(19),
				ParentID:       models.PostID(1),
				Message:        "My post message",
				Created:        date,
				LastEdited:     date.AddDate(0, 0, 1),
				AllowsComments: false,
				Subspace:       "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				OptionalData:   map[string]string{},
				Creator:        owner,
				Medias:         medias,
			},
			expEquals: false,
		},
		{
			name: "Different subspace",
			first: models.Post{
				PostID:         models.PostID(19),
				ParentID:       models.PostID(1),
				Message:        "My post message",
				Created:        date,
				LastEdited:     date.AddDate(0, 0, 1),
				AllowsComments: true,
				Subspace:       "desmos-1",
				OptionalData:   map[string]string{},
				Creator:        owner,
				Medias:         medias,
			},
			second: models.Post{
				PostID:         models.PostID(19),
				ParentID:       models.PostID(1),
				Message:        "My post message",
				Created:        date,
				LastEdited:     date.AddDate(0, 0, 1),
				AllowsComments: true,
				Subspace:       "desmos-2",
				OptionalData:   map[string]string{},
				Creator:        owner,
				Medias:         medias,
			},
			expEquals: false,
		},
		{
			name: "Different optional data",
			first: models.Post{
				PostID:         models.PostID(19),
				ParentID:       models.PostID(1),
				Message:        "My post message",
				Created:        date,
				LastEdited:     date.AddDate(0, 0, 1),
				AllowsComments: true,
				Subspace:       "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				OptionalData: map[string]string{
					"field": "value",
				},
				Creator: owner,
				Medias:  medias,
			},
			second: models.Post{
				PostID:         models.PostID(19),
				ParentID:       models.PostID(1),
				Message:        "My post message",
				Created:        date,
				LastEdited:     date.AddDate(0, 0, 1),
				AllowsComments: true,
				Subspace:       "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				OptionalData: map[string]string{
					"field": "other-value",
				},
				Creator: owner,
				Medias:  medias,
			},
			expEquals: false,
		},
		{
			name: "Different owner",
			first: models.Post{
				PostID:         models.PostID(19),
				ParentID:       models.PostID(1),
				Message:        "My post message",
				Created:        date,
				LastEdited:     date.AddDate(0, 0, 1),
				AllowsComments: true,
				Subspace:       "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				OptionalData:   map[string]string{},
				Creator:        owner,
				Medias:         medias,
			},
			second: models.Post{
				PostID:         models.PostID(19),
				ParentID:       models.PostID(1),
				Message:        "My post message",
				Created:        date,
				LastEdited:     date.AddDate(0, 0, 1),
				AllowsComments: true,
				Subspace:       "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				OptionalData:   map[string]string{},
				Creator:        otherOwner,
				Medias:         medias,
			},
			expEquals: false,
		},
		{
			name: "Different medias",
			first: models.Post{
				PostID:         models.PostID(19),
				ParentID:       models.PostID(1),
				Message:        "My post message",
				Created:        date,
				LastEdited:     date.AddDate(0, 0, 1),
				AllowsComments: true,
				Subspace:       "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				OptionalData:   map[string]string{},
				Creator:        owner,
				Medias:         medias,
			},
			second: models.Post{
				PostID:         models.PostID(19),
				ParentID:       models.PostID(1),
				Message:        "My post message",
				Created:        date,
				LastEdited:     date.AddDate(0, 0, 1),
				AllowsComments: true,
				Subspace:       "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				OptionalData:   map[string]string{},
				Creator:        otherOwner,
				Medias:         models.PostMedias{},
			},
			expEquals: false,
		},
		{
			name: "Different polls",
			first: models.Post{
				PostID:         models.PostID(19),
				ParentID:       models.PostID(1),
				Message:        "My post message",
				Created:        date,
				LastEdited:     date.AddDate(0, 0, 1),
				AllowsComments: true,
				Subspace:       "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				OptionalData:   map[string]string{},
				Creator:        owner,
				Medias:         medias,
				PollData:       nil,
			},
			second: models.Post{
				PostID:         models.PostID(19),
				ParentID:       models.PostID(1),
				Message:        "My post message",
				Created:        date,
				LastEdited:     date.AddDate(0, 0, 1),
				AllowsComments: true,
				Subspace:       "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				OptionalData:   map[string]string{},
				Creator:        otherOwner,
				Medias:         medias,
				PollData:       &models.PollData{},
			},
			expEquals: false,
		},
		{
			name: "Equals posts",
			first: models.Post{
				PostID:         models.PostID(19),
				ParentID:       models.PostID(1),
				Message:        "My post message",
				Created:        date,
				LastEdited:     date.AddDate(0, 0, 1),
				AllowsComments: true,
				Subspace:       "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				OptionalData:   map[string]string{},
				Creator:        owner,
			}.WithMedias(medias).WithPollData(pollData),
			second: models.Post{
				PostID:         models.PostID(19),
				ParentID:       models.PostID(1),
				Message:        "My post message",
				Created:        date,
				LastEdited:     date.AddDate(0, 0, 1),
				AllowsComments: true,
				Subspace:       "4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				OptionalData:   map[string]string{},
				Creator:        owner,
			}.WithMedias(medias).WithPollData(pollData),
			expEquals: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expEquals, test.first.Equals(test.second))
		})
	}
}

func TestPost_GetPostHashtags(t *testing.T) {
	owner, err := sdk.AccAddressFromBech32("cosmos1cjf97gpzwmaf30pzvaargfgr884mpp5ak8f7ns")
	require.NoError(t, err)

	timeZone, err := time.LoadLocation("UTC")
	require.NoError(t, err)

	date := time.Date(2020, 1, 1, 12, 00, 00, 000, timeZone)

	tests := []struct {
		name        string
		post        models.Post
		expHashtags []string
	}{
		{
			name: "Hashtags in message extracted correctly (spaced hashtags)",
			post: models.NewPost(models.PostID(1),
				models.PostID(0),
				"Post with #test #desmos",
				false,
				"4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				map[string]string{},
				date,
				owner,
			),
			expHashtags: []string{"test", "desmos"},
		},
		{
			name: "Hashtags in message extracted correctly (non-spaced hashtags)",
			post: models.NewPost(models.PostID(1),
				models.PostID(0),
				"Post with #test#desmos",
				false,
				"4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				map[string]string{},
				date,
				owner,
			),
			expHashtags: []string{},
		},
		{
			name: "Hashtags in message extracted correctly (underscore separated hashtags)",
			post: models.NewPost(models.PostID(1),
				models.PostID(0),
				"Post with #test_#desmos",
				false,
				"4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				map[string]string{},
				date,
				owner,
			),
			expHashtags: []string{},
		},
		{
			name: "Hashtags in message extracted correctly (only number hashtag)",
			post: models.NewPost(models.PostID(1),
				models.PostID(0),
				"Post with #101112",
				false,
				"4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				map[string]string{},
				date,
				owner,
			),
			expHashtags: []string{},
		},
		{
			name: "No hashtags in message",
			post: models.NewPost(models.PostID(1),
				models.PostID(0),
				"Post with no hashtag",
				false,
				"4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				map[string]string{},
				date,
				owner,
			),
			expHashtags: []string{},
		},
		{
			name: "No same hashtags inside string array",
			post: models.NewPost(models.PostID(1),
				models.PostID(0),
				"Post with double #hashtag #hashtag",
				false,
				"4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
				map[string]string{},
				date,
				owner,
			),
			expHashtags: []string{"hashtag"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			hashtags := test.post.GetPostHashtags()
			require.Equal(t, test.expHashtags, hashtags)
		})
	}
}

// -----------
// --- Posts
// -----------
func TestPosts_String(t *testing.T) {
	owner1, err := sdk.AccAddressFromBech32("cosmos1cjf97gpzwmaf30pzvaargfgr884mpp5ak8f7ns")
	require.NoError(t, err)

	owner2, err := sdk.AccAddressFromBech32("cosmos1r2plnngkwnahajl3d2a7fvzcsxf6djlt380f3l")
	require.NoError(t, err)

	medias := models.PostMedias{
		models.PostMedia{
			URI:      "https://uri.com",
			MimeType: "text/plain",
		},
	}
	answer := models.PollAnswer{
		ID:   models.AnswerID(1),
		Text: "Yes",
	}

	answer2 := models.PollAnswer{
		ID:   models.AnswerID(2),
		Text: "No",
	}
	pollData := models.NewPollData("poll?", time.Now().UTC().Add(time.Hour), models.PollAnswers{answer, answer2}, true, false, true)

	timeZone, err := time.LoadLocation("UTC")
	require.NoError(t, err)

	date := time.Date(2020, 1, 1, 12, 0, 00, 000, timeZone)

	posts := models.Posts{
		models.NewPost(models.PostID(1), models.PostID(10), "Post 1", false, "external-ref-1", map[string]string{}, date, owner1).WithMedias(medias).WithPollData(pollData),
		models.NewPost(models.PostID(2), models.PostID(10), "Post 2", false, "external-ref-1", map[string]string{}, date, owner2).WithMedias(medias).WithPollData(pollData),
	}

	expected := `ID - [Creator] Message
1 - [cosmos1cjf97gpzwmaf30pzvaargfgr884mpp5ak8f7ns] Post 1
2 - [cosmos1r2plnngkwnahajl3d2a7fvzcsxf6djlt380f3l] Post 2`
	require.Equal(t, expected, posts.String())
}
