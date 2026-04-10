# Makefile for building eBPF GTP programs (libbpf compatible)

CLANG ?= clang
LLC ?= llc

ARCH := $(shell uname -m | sed 's/x86_64/x86/')

KERNEL_HEADERS ?= /usr/src/linux-headers-$(shell uname -r)

BPF_CFLAGS := -O2 -g -Wall -Werror \
              -target bpf \
              -D__TARGET_ARCH_$(ARCH) \
              -I$(KERNEL_HEADERS)/include \
              -I$(KERNEL_HEADERS)/include/uapi \
              -I$(KERNEL_HEADERS)/include/generated \
              -I.

# eBPF source files
BPF_SRCS := \
    gtp_decap.c \
    gtp_encap.c \
    gtp_mark.c

# Output objects
BPF_OBJS := $(BPF_SRCS:.c=.o)

.PHONY: all clean

# Build all handlers
all: $(BPF_OBJS)

# Compile eBPF programs
%.o: %.c
    $(CLANG) $(BPF_CFLAGS) -c $< -o $@

# Clean build outputs
clean:
    rm -f $(BPF_OBJS)