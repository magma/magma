// src/pipelined_hook.cpp
#include <iostream>
#include <vector>
#include <string>
#include <csignal>
#include <thread>
#include <chrono>

extern "C" {
#include "bpf_loader.h"
#include "bpf_utils.h"
}

static volatile bool g_exiting = false;

static void sig_handler(int s) {
    g_exiting = true;
}

class PipelinedHook {
public:
    PipelinedHook(const std::string &bpf_obj_path, const std::vector<std::string> &ifs)
        : bpf_path(bpf_obj_path), interfaces(ifs) {
        bpf_loader_init(&loader);
    }

    ~PipelinedHook() {
        cleanup();
    }

    bool load_and_attach() {
        if (bpf_loader_load(&loader, bpf_path.c_str()) != 0) {
            std::cerr << "Failed to load BPF object: " << bpf_path << std::endl;
            return false;
        }

        for (const auto &ifname : interfaces) {
            // Ensure clsact is present
            if (bpf_utils_create_clsact(ifname.c_str()) != 0) {
                std::cerr << "Warning: failed to create clsact on " << ifname << std::endl;
            }
            if (bpf_loader_attach_tc(&loader, ifname.c_str()) != 0) {
                std::cerr << "Failed to attach BPF program to interface: " << ifname << std::endl;
                return false;
            }
            std::cout << "Attached BPF to " << ifname << std::endl;
        }
        return true;
    }

    void detach_all() {
        for (const auto &ifname : interfaces) {
            int idx = bpf_utils_ifindex(ifname.c_str());
            if (idx <= 0) continue;
            bpf_utils_tc_detach(idx);
            // Optionally remove clsact:
            if (bpf_utils_delete_clsact(ifname.c_str()) != 0) {
                std::cerr << "Warning: failed to delete clsact on " << ifname << std::endl;
            }
            std::cout << "Detached BPF from " << ifname << std::endl;
        }
    }

    void cleanup() {
        detach_all();
        bpf_loader_cleanup(&loader);
    }

    // Run loop - typically will be a blocking daemon process.
    void run_loop() {
        std::cout << "PipelinedHook running. Press Ctrl+C to quit.\n";
        while (!g_exiting) {
            std::this_thread::sleep_for(std::chrono::seconds(1));
        }
        std::cout << "PipelinedHook terminating...\n";
    }

private:
    struct bpf_loader loader;
    std::string bpf_path;
    std::vector<std::string> interfaces;
};

// Simple CLI:
// ./pipelined_hook <bpf_object> load|detach|run
int main(int argc, char **argv) {
    if (argc < 3) {
        std::cerr << "Usage: " << argv[0] << " <bpf_object.o> <load|detach|run>\n";
        return 2;
    }

    std::string bpf_obj = argv[1];
    std::string cmd = argv[2];
    std::vector<std::string> ifs = {"br-magma", "gtp0", "sgi0"}; // default interfaces

    signal(SIGINT, sig_handler);
    signal(SIGTERM, sig_handler);

    PipelinedHook hook(bpf_obj, ifs);

    if (cmd == "load") {
        if (!hook.load_and_attach()) return 1;
        hook.cleanup();
        return 0;
    } else if (cmd == "detach") {
        hook.detach_all();
        hook.cleanup();
        return 0;
    } else if (cmd == "run") {
        if (!hook.load_and_attach()) return 1;
        hook.run_loop();
        hook.cleanup();
        return 0;
    } else {
        std::cerr << "Unknown command: " << cmd << "\n";
        return 2;
    }
}
