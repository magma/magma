add_compile_options(-std=c++11)

add_executable(service_registry_test test_service_registry.cpp)

target_link_libraries(service_registry_test
  SERVICE_REGISTRY CONFIG  gtest gtest_main pthread grpc++ dl yaml-cpp
    )

add_test(test_service_registry service_registry_test)
