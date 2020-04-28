package posts

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/capability"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channelexported "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	porttypes "github.com/cosmos/cosmos-sdk/x/ibc/05-port/types"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

// type check to ensure the interface is properly implemented
var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the posts module.
type AppModuleBasic struct{}

// Name returns the posts module's name.
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers the posts module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the auth
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return cdc.MustMarshalJSON(DefaultGenesis())
}

// ValidateGenesis performs genesis state validation for the posts module.
func (AppModuleBasic) ValidateGenesis(_ codec.JSONMarshaler, _ json.RawMessage) error {
	return nil
}

// RegisterRESTRoutes registers the REST routes for the posts module.
func (AppModuleBasic) RegisterRESTRoutes(_ context.CLIContext, _ *mux.Router) {
	// TODO
}

// GetTxCmd returns the root tx command for the posts module.
func (AppModuleBasic) GetQueryCmd(_ *codec.Codec) *cobra.Command {
	// TODO
	return nil
}

// GetQueryCmd returns the root query command for the posts module.
func (AppModuleBasic) GetTxCmd(_ *codec.Codec) *cobra.Command {
	// TODO
	return nil
}

//____________________________________________________________________________

// AppModule implements an application module for the posts module.
type AppModule struct {
	AppModuleBasic
	keeper Keeper
}

// NewAppModule creates a new AppModule Object
func NewAppModule(k Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
	}
}

// RegisterInvariants performs a no-op.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Route returns the message routing key for the posts module.
func (am AppModule) Route() string {
	return RouterKey
}

// NewHandler returns an sdk.Handler for the posts module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute returns the posts module's querier route name.
func (am AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler returns the posts module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	// TODO
	return nil
}

// InitGenesis performs genesis initialization for the posts module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the auth
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONMarshaler) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the posts module.
func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {
}

// EndBlock returns the end blocker for the posts module. It returns no validator
// updates.
func (am AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

//____________________________________________________________________________

// Implement IBCModule callbacks
func (am AppModule) OnChanOpenInit(
	ctx sdk.Context,
	order channelexported.Order,
	connectionHops []string,
	portID string,
	channelID string,
	chanCap *capability.Capability,
	counterparty channeltypes.Counterparty,
	version string,
) error {
	// TODO: Enforce ordering, currently relayers use ORDERED channels

	// Require portID is the portID the posts module is bound to
	boundPort := am.keeper.GetPort(ctx)
	if boundPort != portID {
		return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid port: %s, expected %s", portID, boundPort)
	}

	if version != Version {
		return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid version: %s, expected %s", version, "ics20-1")
	}

	// Claim channel capability passed back by IBC module
	if err := am.keeper.ClaimCapability(ctx, chanCap, ibctypes.ChannelCapabilityPath(portID, channelID)); err != nil {
		return sdkerrors.Wrap(channel.ErrChannelCapabilityNotFound, err.Error())
	}

	return nil
}

func (am AppModule) OnChanOpenTry(
	ctx sdk.Context,
	_ channelexported.Order,
	_ []string,
	portID,
	channelID string,
	chanCap *capability.Capability,
	counterparty channeltypes.Counterparty,
	version,
	counterpartyVersion string,
) error {
	// TODO: Enforce ordering, currently relayers use ORDERED channels

	// Require portID is the portID transfer module is bound to
	boundPort := am.keeper.GetPort(ctx)
	if boundPort != portID {
		return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid port: %s, expected %s", portID, boundPort)
	}

	if version != Version {
		return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid version: %s, expected %s", version, "ics20-1")
	}

	if counterpartyVersion != Version {
		return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid counterparty version: %s, expected %s", counterpartyVersion, "ics20-1")
	}

	// Claim channel capability passed back by IBC module
	if err := am.keeper.ClaimCapability(ctx, chanCap, ibctypes.ChannelCapabilityPath(portID, channelID)); err != nil {
		return sdkerrors.Wrap(channel.ErrChannelCapabilityNotFound, err.Error())
	}

	return nil
}

func (am AppModule) OnChanOpenAck(_ sdk.Context, _, _ string, counterpartyVersion string) error {
	if counterpartyVersion != Version {
		return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid counterparty version: %s, expected %s", counterpartyVersion, "ics20-1")
	}
	return nil
}

func (am AppModule) OnChanOpenConfirm(_ sdk.Context, _, _ string) error {
	return nil
}

func (am AppModule) OnChanCloseInit(_ sdk.Context, _, _ string) error {
	return nil
}

func (am AppModule) OnChanCloseConfirm(_ sdk.Context, _, _ string) error {
	return nil
}

func (am AppModule) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet) (*sdk.Result, error) {
	var data PostCreationPacketData
	if err := ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal post creation packet data: %s", err.Error())
	}
	return handleCreationPacketData(ctx, am.keeper, packet, data)
}

func (am AppModule) OnAcknowledgementPacket(_ sdk.Context, _ channeltypes.Packet, _ []byte) (*sdk.Result, error) {
	return nil, nil
}

func (am AppModule) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet) (*sdk.Result, error) {
	return nil, nil
}
