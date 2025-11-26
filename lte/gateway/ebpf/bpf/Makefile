# Makefile for building eBPF programs
# Target: tc_core.bpf.o

BPF_CLANG ?= clang
BPF_CFLAGS ?= -O2 -target bpf -Wall -Werror
BPF_LD_FLAGS ?=

# Files
SRC = tc_core.bpf.c
OBJ = $(SRC:.c=.o)

all: $(OBJ)

%.o: %.c
	$(BPF_CLANG) $(BPF_CFLAGS) -c $< -o $@

clean:
	rm -f $(OBJ)

.PHONY: all clean
