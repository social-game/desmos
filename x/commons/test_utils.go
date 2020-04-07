package commons

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func RequireContainsEvent(t *testing.T, events []abci.Event, event sdk.Event) {
	expected := sdk.Events{event}.ToABCIEvents()[0]

	found := false
	for _, e := range events {
		if assert.ObjectsAreEqual(e, expected) {
			found = true
			break
		}
	}

	require.True(t, found)
}
