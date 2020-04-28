package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	commitmenttypes "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment/types"
	"github.com/desmos-labs/desmos/x/posts"
)

var ModuleCdc = codec.New()

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCrossPost{}, "desmos/MsgCrossPost", nil)
	cdc.RegisterConcrete(PostCreationPacketData{}, "ibc/desmos/PostCreationPacketData", nil)
}

func init() {
	RegisterCodec(ModuleCdc)
	posts.RegisterCodec(ModuleCdc)
	channel.RegisterCodec(ModuleCdc)
	commitmenttypes.RegisterCodec(ModuleCdc)
}
