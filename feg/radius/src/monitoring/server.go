package monitoring

// ServerCounters ...
type ServerCounters struct {
	// Init operation counter for server initialization
	Init Operation

	// ListenerInit counterset for listener initialization
	ListenerInit Operation

	// FilterInit counterset for module initialization
	FilterInit Operation

	// ModuleInit counterset for module initialization
	ModuleInit Operation

	// DedupPacket RADIUS dedup logic counter
	DedupPacket Operation
}

// CreateServerCounters ...
func CreateServerCounters() *ServerCounters {
	return &ServerCounters{
		Init:         NewOperation("server_init"),
		ListenerInit: NewOperation("listener_init"),
		FilterInit:   NewOperation("filter_init"),
		ModuleInit:   NewOperation("module_init"),
		DedupPacket:  NewOperation("radius_dedup"),
	}
}
