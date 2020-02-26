set(GTPV2C_DIR  ${CMAKE_CURRENT_SOURCE_DIR}/nwgtpv2c-0.11/src)

add_library(LIB_GTPV2C
    ${GTPV2C_DIR}/NwGtpv2cTrxn.c
    ${GTPV2C_DIR}/NwGtpv2cTunnel.c
    ${GTPV2C_DIR}/NwGtpv2cMsg.c
    ${GTPV2C_DIR}/NwGtpv2cMsgIeParseInfo.c
    ${GTPV2C_DIR}/NwGtpv2cMsgParser.c
    ${GTPV2C_DIR}/NwGtpv2c.c
)
target_link_libraries(LIB_GTPV2C
    COMMON
    LIB_BSTR
)
target_include_directories(LIB_GTPV2C PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}/nwgtpv2c-0.11/include/
    ${CMAKE_CURRENT_SOURCE_DIR}/nwgtpv2c-0.11/shared/
    ${CMAKE_CURRENT_SOURCE_DIR}/nwgtpv2c-0.11/src/
)
