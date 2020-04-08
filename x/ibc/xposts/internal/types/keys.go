package types

const (
	// ModuleName defines the IBC posts name
	ModuleName = "xposts"

	// Version defines the current version the IBC posts
	// module supports
	Version = "ics20-1"

	// PortID that IBC posts module binds to
	PortID = "posts"

	// StoreKey is the store key string for IBC posts
	StoreKey = ModuleName

	// RouterKey is the message route for IBC posts
	RouterKey = ModuleName

	// QuerierRoute is the querier route for IBC posts
	QuerierRoute = ModuleName
)
