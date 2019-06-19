package counters

var (
	// ServerInit Counterset for server initialization
	ServerInit = NewOperation("server_init")

	// FilterInit counterset for module initialization
	FilterInit = NewOperation("filter_init", FilterTag)

	// ListenerInit counterset for listener initialization
	ListenerInit = NewOperation("listener_init", ListenerTag)

	// ModuleInit counterset for module initialization
	ModuleInit = NewOperation("module_init", ListenerTag, ModuleTag)

	// DedupPacket RADIUS dedup logic counter
	DedupPacket = NewOperation("radius_dedup", ListenerTag)
)
