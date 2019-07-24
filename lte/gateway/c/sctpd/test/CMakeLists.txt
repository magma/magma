add_library(SCTPD_TEST_LIB

)

target_link_libraries(SCTPD_TEST_LIB SCTPD_LIB gmock_main pthread rt)

foreach(sctpd_test sctp_desc event_handler)
  add_executable(${sctpd_test}_test test_${sctpd_test}.cpp)
  target_link_libraries(${sctpd_test}_test SCTPD_TEST_LIB)
  add_test(test_${sctpd_test} ${sctpd_test}_test)
endforeach(sctpd_test)
