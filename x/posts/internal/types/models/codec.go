package models

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// MsgsCodec is the codec
var ModelsCdc = codec.New()

func init() {
	RegisterModelsCodec(ModelsCdc)
}

// RegisterModelsCodec registers concrete types on the Amino codec
func RegisterModelsCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(PostCreationData{}, "ibc/desmos/PacketDataPostCreation", nil)
}
