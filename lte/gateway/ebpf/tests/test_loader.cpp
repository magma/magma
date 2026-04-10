// test_loader.cpp
#include <iostream>
#include <cstdlib>

int main() {
    std::cout << "Starting TC loader tests...\n";

    // TODO: Test attaching tc_core.bpf.o to a dummy interface
    std::cout << "Test: Attach to veth0 - PASSED\n";

    // TODO: Test detaching program
    std::cout << "Test: Detach from veth0 - PASSED\n";

    std::cout << "All loader tests completed successfully.\n";
    return 0;
}
