package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// MsgsCodec is the codec
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

func RegisterCodec(cdc *codec.Codec) {
	RegisterModelsCodec(cdc)
	RegisterMessagesCodec(cdc)
}
