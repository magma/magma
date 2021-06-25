cc_library(
	name = "logging",
	srcs = [],
	hdrs = ["magma_logging.h", "magma_logging_init.h"],
)
cc_library(
	name = "utilities",
	srcs = ["Utilities.cpp", "GrpcMagmaUtils.cpp", "DiameterCodes.cpp", "SessionID.cpp"],
	hdrs = ["Utilities.h", "GrpcMagmaUtils.h", "DiameterCodes.h", "SessionID.h"],
	deps = [":logging"],
)
