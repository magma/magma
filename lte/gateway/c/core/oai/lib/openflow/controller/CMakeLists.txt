add_compile_options(-std=c++11)

set(OUTPUT_DIR "${CMAKE_CURRENT_BINARY_DIR}")

set(CMAKE_CXX_FLAGS "-Wno-write-strings -Wno-literal-suffix")

add_library(LIB_OPENFLOW_CONTROLLER
    ControllerMain.cpp
    OpenflowController.cpp
    PagingApplication.cpp
    ControllerEvents.cpp
    BaseApplication.cpp
    OpenflowMessenger.cpp
    GTPApplication.cpp
    IMSIEncoder.cpp
    )
target_link_libraries(LIB_OPENFLOW_CONTROLLER
    COMMON
    fluid_base fluid_msg
    LIB_BSTR
    TASK_SGW
    )
target_include_directories(LIB_OPENFLOW_CONTROLLER PUBLIC
    $ENV{MAGMA_ROOT}
    ${CMAKE_CURRENT_SOURCE_DIR}
    )
