cc_library(
	name = "logging",
	srcs = [],
	hdrs = ["magma_logging.h", "magma_logging_init.h"],
)
cc_library(
	name = "utilities",
	srcs = [
		"Utilities.cpp", 
		"GrpcMagmaUtils.cpp", 
		"DiameterCodes.cpp", 
		"SessionID.cpp",
	],
	hdrs = ["Utilities.h", "GrpcMagmaUtils.h", "DiameterCodes.h", "SessionID.h"],
	deps = [":logging"],
)

cc_library(
	name = "credit_key",
	hdrs = ["CreditKey.h"],
	deps = [
		"//lte/protos:policydb_cc_proto",
		"//lte/protos:session_manager_cc_proto",
	],
)
