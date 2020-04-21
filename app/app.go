package app

import (
	"io"
	"os"

	codecstd "github.com/cosmos/cosmos-sdk/codec/std"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	port "github.com/cosmos/cosmos-sdk/x/ibc/05-port"
	transfer "github.com/cosmos/cosmos-sdk/x/ibc/20-transfer"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	"github.com/desmos-labs/desmos/x/commons"
	"github.com/desmos-labs/desmos/x/posts"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	ibcposts "github.com/desmos-labs/desmos/x/ibc/xposts"
	"github.com/desmos-labs/desmos/x/magpie"
)

const (
	appName = "desmos"
)

var (
	// DefaultCLIHome represents the default home directory for the application CLI
	DefaultCLIHome = os.ExpandEnv("$HOME/.desmoscli")

	// DefaultNodeHome sets the folder where the application data and configuration will be stored
	DefaultNodeHome = os.ExpandEnv("$HOME/.desmosd")

	// ModuleBasics is in charge of setting up basic module elements
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(
			paramsclient.ProposalHandler, distr.ProposalHandler, upgradeclient.ProposalHandler,
		),
		params.AppModuleBasic{},
		slashing.AppModuleBasic{},
		ibc.AppModuleBasic{},
		upgrade.AppModuleBasic{},

		// Custom modules
		magpie.AppModuleBasic{},
		posts.AppModuleBasic{},

		// IBC modules
		transfer.AppModuleBasic{},
		ibcposts.AppModuleBasic{},
	)

	// Module account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:           nil,
		distr.ModuleName:                nil,
		staking.BondedPoolName:          {auth.Burner, auth.Staking},
		staking.NotBondedPoolName:       {auth.Burner, auth.Staking},
		gov.ModuleName:                  {auth.Burner},
		transfer.GetModuleAccountName(): {auth.Minter, auth.Burner},
	}

	// module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{
		distr.ModuleName: true,
	}
)

// MakeCodec generates the necessary codecs for Amino
func MakeCodec() *codec.Codec {
	var cdc = codec.New()

	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	codec.RegisterEvidences(cdc)
	authvesting.RegisterCodec(cdc)

	return cdc.Seal()
}

// Verify app interface at compile time
var _ simapp.App = (*DesmosApp)(nil)

// DesmosApp extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities arKen't needed for testing.
type DesmosApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	// keys to access the substores
	keys    map[string]*sdk.KVStoreKey
	tkeys   map[string]*sdk.TransientStoreKey
	memKeys map[string]*sdk.MemoryStoreKey

	// subspaces
	subspaces map[string]params.Subspace

	// Keepers
	AccountKeeper    auth.AccountKeeper
	BankKeeper       bank.Keeper
	CapabilityKeeper *capability.Keeper
	StakingKeeper    staking.Keeper
	SlashingKeeper   slashing.Keeper
	DistrKeeper      distr.Keeper
	GovKeeper        gov.Keeper
	UpgradeKeeper    upgrade.Keeper
	ParamsKeeper     params.Keeper
	IBCKeeper        *ibc.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly

	// Custom modules
	MagpieKeeper magpie.Keeper
	PostsKeeper  posts.Keeper

	// IBC modules
	TransferKeeper transfer.Keeper
	IBCPostsKeeper ibcposts.Keeper

	// Module Manager
	mm *module.Manager

	// Simulation manager
	sm *module.SimulationManager
}

// NewDesmosApp is a constructor function for DesmosApp
func NewDesmosApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	skipUpgradeHeights map[int64]bool, home string, baseAppOptions ...func(*bam.BaseApp),
) *DesmosApp {
	// First define the top level codec that will be shared by the different modules
	// TODO: Remove cdc in favor of appCodec once all modules are migrated.
	cdc := codecstd.MakeCodec(ModuleBasics)
	appCodec := codecstd.NewAppCodec(cdc)

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)
	keys := sdk.NewKVStoreKeys(
		auth.StoreKey, bank.StoreKey, staking.StoreKey,
		distr.StoreKey, slashing.StoreKey,
		gov.StoreKey, params.StoreKey, ibc.StoreKey, transfer.StoreKey,
		upgrade.StoreKey, capability.StoreKey,

		// Custom modules
		magpie.StoreKey, posts.StoreKey,

		// IBC modules
		ibcposts.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(params.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capability.MemStoreKey)

	// Here you initialize your application with the store keys it requires
	var app = &DesmosApp{
		BaseApp:   bApp,
		cdc:       cdc,
		keys:      keys,
		tkeys:     tkeys,
		memKeys:   memKeys,
		subspaces: make(map[string]params.Subspace),
	}

	// Init params keeper and subspaces
	app.ParamsKeeper = params.NewKeeper(appCodec, keys[params.StoreKey], tkeys[params.TStoreKey])
	app.subspaces[auth.ModuleName] = app.ParamsKeeper.Subspace(auth.DefaultParamspace)
	app.subspaces[bank.ModuleName] = app.ParamsKeeper.Subspace(bank.DefaultParamspace)
	app.subspaces[staking.ModuleName] = app.ParamsKeeper.Subspace(staking.DefaultParamspace)
	app.subspaces[distr.ModuleName] = app.ParamsKeeper.Subspace(distr.DefaultParamspace)
	app.subspaces[slashing.ModuleName] = app.ParamsKeeper.Subspace(slashing.DefaultParamspace)
	app.subspaces[gov.ModuleName] = app.ParamsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())

	// set the BaseApp's parameter store
	bApp.SetParamStore(app.ParamsKeeper.Subspace(bam.Paramspace).WithKeyTable(std.ConsensusParamsKeyTable()))

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capability.NewKeeper(appCodec, keys[capability.StoreKey], memKeys[capability.MemStoreKey])
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibc.ModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(transfer.ModuleName)
	scopedPostsKeeper := app.CapabilityKeeper.ScopeToModule(ibcposts.ModuleName)

	// Add keepers
	app.AccountKeeper = auth.NewAccountKeeper(
		appCodec, keys[auth.StoreKey], app.subspaces[auth.ModuleName], auth.ProtoBaseAccount, maccPerms,
	)
	app.BankKeeper = bank.NewBaseKeeper(
		appCodec, keys[bank.StoreKey], app.AccountKeeper, app.subspaces[bank.ModuleName], app.BlacklistedAccAddrs(),
	)
	stakingKeeper := staking.NewKeeper(
		appCodec, keys[staking.StoreKey], app.AccountKeeper, app.BankKeeper, app.subspaces[staking.ModuleName],
	)
	app.DistrKeeper = distr.NewKeeper(
		appCodec, keys[distr.StoreKey], app.subspaces[distr.ModuleName], app.AccountKeeper, app.BankKeeper,
		&stakingKeeper, auth.FeeCollectorName, app.ModuleAccountAddrs(),
	)
	app.SlashingKeeper = slashing.NewKeeper(
		appCodec, keys[slashing.StoreKey], &stakingKeeper, app.subspaces[slashing.ModuleName],
	)
	app.UpgradeKeeper = upgrade.NewKeeper(skipUpgradeHeights, keys[upgrade.StoreKey], appCodec, home)

	// Register the proposal types
	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(upgrade.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper))
	app.GovKeeper = gov.NewKeeper(
		appCodec, keys[gov.StoreKey], app.subspaces[gov.ModuleName], app.AccountKeeper, app.BankKeeper,
		&stakingKeeper, govRouter,
	)

	// Register the staking hooks
	// NOTE: StakingKeeper above is passed by reference, so that it will contain these hooks
	app.StakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	// Create IBC Keeper
	app.IBCKeeper = ibc.NewKeeper(
		app.cdc, keys[ibc.StoreKey], app.StakingKeeper, scopedIBCKeeper,
	)

	// Register custom modules
	app.MagpieKeeper = magpie.NewKeeper(app.cdc, keys[magpie.StoreKey])
	app.PostsKeeper = posts.NewKeeper(app.cdc, keys[posts.StoreKey])

	// Create IBC modules
	app.TransferKeeper = transfer.NewKeeper(
		app.cdc, keys[transfer.StoreKey],
		app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper,
		app.AccountKeeper, app.BankKeeper, scopedTransferKeeper,
	)
	transferModule := transfer.NewAppModule(app.TransferKeeper)

	app.IBCPostsKeeper = ibcposts.NewKeeper(
		app.cdc, keys[ibcposts.StoreKey], app.PostsKeeper,
		app.IBCKeeper.ChannelKeeper, app.IBCKeeper.PortKeeper, scopedPostsKeeper,
	)
	ibcPostsModule := ibcposts.NewAppModule(app.IBCPostsKeeper, app.PostsKeeper)

	// Create static IBC router, add posts route, then set and seal it
	ibcRouter := port.NewRouter()
	ibcRouter.AddRoute(transfer.ModuleName, transferModule)
	ibcRouter.AddRoute(ibcposts.ModuleName, ibcPostsModule)
	app.IBCKeeper.SetRouter(ibcRouter)

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.AccountKeeper),
		bank.NewAppModule(app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(*app.CapabilityKeeper),
		gov.NewAppModule(app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		slashing.NewAppModule(app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		distr.NewAppModule(app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		staking.NewAppModule(app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		params.NewAppModule(app.ParamsKeeper),

		// Custom modules
		magpie.NewAppModule(app.MagpieKeeper, app.AccountKeeper),
		posts.NewAppModule(app.PostsKeeper, app.AccountKeeper, app.BankKeeper),

		// IBC Modules
		transferModule,
		ibcPostsModule,
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	app.mm.SetOrderBeginBlockers(
		upgrade.ModuleName, distr.ModuleName, slashing.ModuleName,
		staking.ModuleName, ibc.ModuleName,
	)
	app.mm.SetOrderEndBlockers(gov.ModuleName, staking.ModuleName)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		auth.ModuleName, distr.ModuleName, staking.ModuleName, bank.ModuleName,
		slashing.ModuleName, gov.ModuleName,
		ibc.ModuleName, genutil.ModuleName, transfer.ModuleName,

		// Custom modules
		magpie.ModuleName, posts.ModuleName,

		// IBC modules
		ibcposts.ModuleName,
	)

	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions
	app.sm = module.NewSimulationManager(
		auth.NewAppModule(app.AccountKeeper),
		bank.NewAppModule(app.BankKeeper, app.AccountKeeper),
		gov.NewAppModule(app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		staking.NewAppModule(app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		distr.NewAppModule(app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		slashing.NewAppModule(app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		params.NewAppModule(app.ParamsKeeper),

		// Custom modules
		posts.NewAppModule(app.PostsKeeper, app.AccountKeeper, app.BankKeeper),
		magpie.NewAppModule(app.MagpieKeeper, app.AccountKeeper),
	)
	app.sm.RegisterStoreDecoders()

	// Initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// Initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(ante.NewAnteHandler(
		app.AccountKeeper, app.BankKeeper, *app.IBCKeeper, auth.DefaultSigVerificationGasConsumer,
	))
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	// Initialize and seal the capability keeper so all persistent capabilities
	// are loaded in-memory and prevent any further modules from creating scoped
	// sub-keepers.
	ctx := app.BaseApp.NewContext(true, abci.Header{})
	app.CapabilityKeeper.InitializeAndSeal(ctx)

	return app
}

// SetupConfig sets up the given config as it should be for Desmos
func SetupConfig(config *sdk.Config) {
	config.SetBech32PrefixForAccount(
		commons.Bech32MainPrefix,
		commons.Bech32MainPrefix+sdk.PrefixPublic,
	)
	config.SetBech32PrefixForValidator(
		commons.Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixOperator,
		commons.Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixOperator+sdk.PrefixPublic,
	)
	config.SetBech32PrefixForConsensusNode(
		commons.Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixConsensus,
		commons.Bech32MainPrefix+sdk.PrefixValidator+sdk.PrefixConsensus+sdk.PrefixPublic,
	)

	// 852 is the international dialing code of Hong Kong
	// Following the coin type registered at https://github.com/satoshilabs/slips/blob/master/slip-0044.md
	config.SetCoinType(852)
}

// BeginBlocker application updates every begin block
func (app *DesmosApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *DesmosApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer application update.md at chain initialization
func (app *DesmosApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	return app.mm.InitGenesis(ctx, app.cdc, genesisState)
}

// LoadHeight loads a particular height
func (app *DesmosApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *DesmosApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[auth.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// BlacklistedAccAddrs returns all the app's module account addresses black listed for receiving tokens.
func (app *DesmosApp) BlacklistedAccAddrs() map[string]bool {
	blacklistedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blacklistedAddrs[auth.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	return blacklistedAddrs
}

// Codec returns the application's sealed codec.
func (app *DesmosApp) Codec() *codec.Codec {
	return app.cdc
}

// SimulationManager implements the SimulationApp interface
func (app *DesmosApp) SimulationManager() *module.SimulationManager {
	return app.sm
}
