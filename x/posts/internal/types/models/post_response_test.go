package models_test

import (
	"encoding/json"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/x/posts/internal/types/models"
	"github.com/stretchr/testify/require"
)

func TestPostQueryResponse_MarshalJSON(t *testing.T) {
	postOwner, err := sdk.AccAddressFromBech32("cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47")
	require.NoError(t, err)

	liker, err := sdk.AccAddressFromBech32("cosmos1s3nh6tafl4amaxkke9kdejhp09lk93g9ev39r4")
	require.NoError(t, err)

	otherLiker, err := sdk.AccAddressFromBech32("cosmos15lt0mflt6j9a9auj7yl3p20xec4xvljge0zhae")
	require.NoError(t, err)

	timeZone, err := time.LoadLocation("UTC")
	require.NoError(t, err)

	date := time.Date(2020, 2, 2, 15, 0, 0, 0, timeZone)
	medias := models.PostMedias{
		models.PostMedia{
			URI:      "https://uri.com",
			MimeType: "text/plain",
		},
	}

	testPostEndPollDate := time.Date(2050, 1, 1, 15, 15, 00, 000, timeZone)

	pollData := models.NewPollData(
		"poll?",
		testPostEndPollDate,
		models.NewPollAnswers(
			models.NewPollAnswer(models.AnswerID(1), "Yes"),
			models.NewPollAnswer(models.AnswerID(2), "No"),
		), true, false, true)

	answers2 := []models.AnswerID{models.AnswerID(1)}

	post := models.NewPost(
		models.PostID(10),
		models.PostID(0),
		"Post",
		true,
		"4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
		map[string]string{},
		date,
		postOwner,
	).WithMedias(medias).WithPollData(pollData)

	postNoMedia := models.NewPost(
		models.PostID(10),
		models.PostID(0),
		"Post",
		true,
		"4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
		map[string]string{},
		date,
		postOwner,
	).WithPollData(pollData)

	postNoPoll := models.NewPost(
		models.PostID(10),
		models.PostID(0),
		"Post",
		true,
		"4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
		map[string]string{},
		date,
		postOwner,
	).WithMedias(medias)

	answersDetails := []models.UserAnswer{models.NewUserAnswer(answers2, liker)}

	likes := models.PostReactions{
		models.NewPostReaction("like", liker),
		models.NewPostReaction("like", otherLiker),
	}
	children := models.PostIDs{models.PostID(98), models.PostID(100)}

	PostResponse := models.NewPostResponse(post, answersDetails, likes, children)

	tests := []struct {
		name        string
		response    models.PostQueryResponse
		expResponse string
	}{
		{
			name:        "Post Query Response with Post that contains media and poll",
			response:    PostResponse,
			expResponse: `{"id":"10","parent_id":"0","message":"Post","created":"2020-02-02T15:00:00Z","last_edited":"0001-01-01T00:00:00Z","allows_comments":true,"subspace":"4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e","creator":"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47","medias":[{"uri":"https://uri.com","mime_type":"text/plain"}],"poll_data":{"question":"poll?","provided_answers":[{"id":"1","text":"Yes"},{"id":"2","text":"No"}],"end_date":"2050-01-01T15:15:00Z","is_open":true,"allows_multiple_answers":false,"allows_answer_edits":true},"poll_answers":[{"answers":["1"],"user":"cosmos1s3nh6tafl4amaxkke9kdejhp09lk93g9ev39r4"}],"reactions":[{"owner":"cosmos1s3nh6tafl4amaxkke9kdejhp09lk93g9ev39r4","value":"like"},{"owner":"cosmos15lt0mflt6j9a9auj7yl3p20xec4xvljge0zhae","value":"like"}],"children":["98","100"]}`,
		},
		{
			name:        "Post Query with Post that not contains poll",
			response:    models.NewPostResponse(postNoPoll, nil, likes, children),
			expResponse: `{"id":"10","parent_id":"0","message":"Post","created":"2020-02-02T15:00:00Z","last_edited":"0001-01-01T00:00:00Z","allows_comments":true,"subspace":"4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e","creator":"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47","medias":[{"uri":"https://uri.com","mime_type":"text/plain"}],"reactions":[{"owner":"cosmos1s3nh6tafl4amaxkke9kdejhp09lk93g9ev39r4","value":"like"},{"owner":"cosmos15lt0mflt6j9a9auj7yl3p20xec4xvljge0zhae","value":"like"}],"children":["98","100"]}`,
		},
		{
			name:        "Post Query Response with Post that not contains media",
			response:    models.NewPostResponse(postNoMedia, answersDetails, likes, children),
			expResponse: `{"id":"10","parent_id":"0","message":"Post","created":"2020-02-02T15:00:00Z","last_edited":"0001-01-01T00:00:00Z","allows_comments":true,"subspace":"4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e","creator":"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47","poll_data":{"question":"poll?","provided_answers":[{"id":"1","text":"Yes"},{"id":"2","text":"No"}],"end_date":"2050-01-01T15:15:00Z","is_open":true,"allows_multiple_answers":false,"allows_answer_edits":true},"poll_answers":[{"answers":["1"],"user":"cosmos1s3nh6tafl4amaxkke9kdejhp09lk93g9ev39r4"}],"reactions":[{"owner":"cosmos1s3nh6tafl4amaxkke9kdejhp09lk93g9ev39r4","value":"like"},{"owner":"cosmos15lt0mflt6j9a9auj7yl3p20xec4xvljge0zhae","value":"like"}],"children":["98","100"]}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			jsonData, err := json.Marshal(&test.response)
			require.NoError(t, err)
			require.Equal(t, test.expResponse, string(jsonData))
		})
	}
}

func TestPostQueryResponse_String(t *testing.T) {
	postOwner, err := sdk.AccAddressFromBech32("cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47")
	require.NoError(t, err)

	liker, err := sdk.AccAddressFromBech32("cosmos1s3nh6tafl4amaxkke9kdejhp09lk93g9ev39r4")
	require.NoError(t, err)

	otherLiker, err := sdk.AccAddressFromBech32("cosmos15lt0mflt6j9a9auj7yl3p20xec4xvljge0zhae")
	require.NoError(t, err)

	timeZone, err := time.LoadLocation("UTC")
	require.NoError(t, err)

	answers2 := []models.AnswerID{models.AnswerID(1)}
	answersDetails := []models.UserAnswer{models.NewUserAnswer(answers2, liker)}

	date := time.Date(2020, 2, 2, 15, 0, 0, 0, timeZone)
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

	post := models.NewPost(
		models.PostID(10),
		models.PostID(0),
		"Post",
		true,
		"4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
		map[string]string{},
		date,
		postOwner,
	).WithMedias(medias).WithPollData(pollData)

	likes := models.PostReactions{
		models.NewPostReaction("like", liker),
		models.NewPostReaction("like", otherLiker),
	}
	children := models.PostIDs{models.PostID(98), models.PostID(100)}

	PostResponse := models.NewPostResponse(post, answersDetails, likes, children)

	tests := []struct {
		name        string
		response    models.PostQueryResponse
		expResponse string
	}{
		{
			name:        "Post query response string",
			response:    PostResponse,
			expResponse: "ID - [PostReactions] [Children] \n10 - [[{\"owner\":\"cosmos1s3nh6tafl4amaxkke9kdejhp09lk93g9ev39r4\",\"value\":\"like\"} {\"owner\":\"cosmos15lt0mflt6j9a9auj7yl3p20xec4xvljge0zhae\",\"value\":\"like\"}]] [[98 100]]",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			stringResponse := test.response.String()
			require.Equal(t, test.expResponse, stringResponse)
		})
	}
}
