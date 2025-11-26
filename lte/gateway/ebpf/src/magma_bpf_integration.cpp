#include "bpf_loader.h"
#include <iostream>
#include <vector>
#include <string>
#include <csignal>

class MagmaBpfIntegration {
public:
    MagmaBpfIntegration(const std::vector<std::string>& ifaces, const std::string& bpf_file)
        : interfaces(ifaces), bpf_object_file(bpf_file) {}

    bool initialize() {
        if (bpf_loader_init(&loader) != 0) {
            std::cerr << "Failed to initialize BPF loader" << std::endl;
            return false;
        }

        if (bpf_loader_load(&loader, bpf_object_file.c_str()) != 0) {
            std::cerr << "Failed to load BPF object: " << bpf_object_file << std::endl;
            return false;
        }

        // Attach BPF program to all interfaces
        for (const auto& iface : interfaces) {
            if (bpf_loader_attach_tc(&loader, iface.c_str()) != 0) {
                std::cerr << "Failed to attach BPF to interface: " << iface << std::endl;
                return false;
            }
            std::cout << "Attached BPF to interface: " << iface << std::endl;
        }

        return true;
    }

    void cleanup() {
        for (const auto& iface : interfaces) {
            bpf_loader_detach_tc(&loader, iface.c_str());
        }
        bpf_loader_cleanup(&loader);
        std::cout << "Cleaned up BPF loader and detached programs" << std::endl;
    }

private:
    struct bpf_loader loader;
    std::vector<std::string> interfaces;
    std::string bpf_object_file;
};

// Global pointer for signal handling
static MagmaBpfIntegration* global_integration = nullptr;

void signal_handler(int signum) {
    if (global_integration) {
        std::cout << "Caught signal " << signum << ", cleaning up..." << std::endl;
        global_integration->cleanup();
    }
    exit(signum);
}

int main(int argc, char** argv) {
    if (argc < 2) {
        std::cerr << "Usage: " << argv[0] << " <bpf_object_file>" << std::endl;
        return 1;
    }

    std::string bpf_file = argv[1];
    std::vector<std::string> interfaces = {"br-magma", "gtp0", "sgi0"};

    MagmaBpfIntegration integration(interfaces, bpf_file);
    global_integration = &integration;

    // Setup signal handlers for graceful cleanup
    std::signal(SIGINT, signal_handler);
    std::signal(SIGTERM, signal_handler);

    if (!integration.initialize()) {
        std::cerr << "Failed to initialize Magma BPF integration" << std::endl;
        return 1;
    }

    std::cout << "Magma BPF integration running. Press Ctrl+C to exit." << std::endl;

    // Main loop: in production, this could read events, sync maps, etc.
    while (true) {
        std::this_thread::sleep_for(std::chrono::seconds(1));
    }

    integration.cleanup();
    return 0;
}
