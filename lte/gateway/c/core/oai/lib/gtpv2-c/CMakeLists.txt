set(NWGTPV2C_DIR ${PROJECT_SOURCE_DIR}/lib/gtpv2-c/nwgtpv2c-0.11/src)
set(NWGTPV2C_IE_FORMATTER_DIR ${PROJECT_SOURCE_DIR}/lib/gtpv2-c/gtpv2c_ie_formatter/src)

add_library(LIB_GTPV2C
    ${NWGTPV2C_DIR}/NwGtpv2cTrxn.c
    ${NWGTPV2C_DIR}/NwGtpv2cTunnel.c
    ${NWGTPV2C_DIR}/NwGtpv2cMsg.c
    ${NWGTPV2C_DIR}/NwGtpv2cMsgIeParseInfo.c
    ${NWGTPV2C_DIR}/NwGtpv2cMsgParser.c
    ${NWGTPV2C_DIR}/NwGtpv2c.c
    ${NWGTPV2C_IE_FORMATTER_DIR}/gtpv2c_ie_formatter.c
    )
target_link_libraries(LIB_GTPV2C
    COMMON
    LIB_BSTR
    )

target_include_directories(LIB_GTPV2C PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}/nwgtpv2c-0.11/include/
    ${CMAKE_CURRENT_SOURCE_DIR}/nwgtpv2c-0.11/shared/
    ${CMAKE_CURRENT_SOURCE_DIR}/nwgtpv2c-0.11/src/
    ${PROJECT_SOURCE_DIR}/tasks/nas/ies/
    ${CMAKE_CURRENT_SOURCE_DIR}/
    ${CMAKE_CURRENT_SOURCE_DIR}/gtpv2c_ie_formatter/shared/
    )