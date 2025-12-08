# Makefile for building eBPF programs
# Target: tc_core.bpf.o

BPF_CLANG ?= clang
KERNEL_HEADERS ?= /usr/src/linux-headers-$(shell uname -r)
BPF_CFLAGS ?= -O2 -target bpf -Wall -Werror -Wno-unused-function -I$(KERNEL_HEADERS)/include -I$(KERNEL_HEADERS)/include/uapi

# Files
SRC = tc_core.bpf.c
OBJ = $(SRC:.c=.o)

all: $(OBJ)

%.o: %.c
	$(BPF_CLANG) $(BPF_CFLAGS) -c $< -o $@

clean:
	rm -f $(OBJ)

.PHONY: all clean
