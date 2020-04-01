package types_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/desmos-labs/desmos/x/posts/internal/types"
	"github.com/stretchr/testify/require"
)

func TestCreatePostPacketData_MarshalJSON(t *testing.T) {
	data := types.NewPostCreationData(
		"My new post",
		types.PostID(53),
		false,
		"4e188d9c17150037d5199bbdb91ae1eb2a78a15aca04cb35530cccb81494b36e",
		map[string]string{},
		testOwner,
		date,
		types.NewPostMedias(types.NewPostMedia("https://uri.com", "text/plain")),
		&pollData,
	)
	packetData := types.NewCreatePostPacketData(data, 100)

	fmt.Println(string(packetData.GetBytes()))
}

func TestDesmosAddressMarshalUnmarshalJSON(t *testing.T) {
	address := types.DesmosAddress{AccAddress: testOwner}

	bz, err := address.MarshalJSON()
	require.NoError(t, err)

	var addr types.DesmosAddress
	err = json.Unmarshal(bz, &addr)
	require.NoError(t, err)

	require.True(t, address.Equals(addr))
}
