package types

const (
	// ModuleName defines the module name
	ModuleName = "lottery"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_lottery"

	// Version defines the current version the IBC module supports
	Version = "lottery-1"

	// PortID is the default port id that module binds to
	PortID = "lottery"
)

var (
	// PortKey defines the key to store the port ID in store
	PortKey = KeyPrefix("lottery-port-")
)

const (
	LotteryKey      = "Lottery/value/"
	LotteryCountKey = "Lottery/count/"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
