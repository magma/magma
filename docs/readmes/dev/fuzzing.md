---
id: dev_fuzzing
title: Fuzz Testing
hide_title: true
---
# Fuzz Testing

## Writing a Fuzz Test

### The Fuzzer Main

See [AttachAcceptFuzzer.cpp](../../../lte/gateway/c/oai/tasks/nas/emm/msg/AttachAcceptFuzzer.cpp) for an example implementation, and the associated [CMakeLists.txt](../../../lte/gateway/c/oai/tasks/nas/CMakeLists.txt) modifications.

### Fuzzing Engine Inputs

The fuzzer wants "nominal" example inputs from which it can then explore efficiently in search of bad behavior.  You would ideally collect small example binary input files and place them in e.g. AttachAcceptFuzzerInputs directory along side the Fuzz test (these will be later used by our enrollment in [OSS-Fuzz](https://google.github.io/oss-fuzz/)).  Here I have no such example inputs and use random bytes, but still achieved ~6 unique segfaults.


## American Fuzzy Lop Docs

Here we use the [American Fuzzy Lop](https://lcamtuf.coredump.cx/afl/) fuzzing engine to do directed generation of random inputs to our Fuzz test. For explanation on AFL, see the official website. Also useful is [Google AFL info](https://github.com/google/AFL). Note we might have used Clang's built-in fuzzing engine support, if we were building with Clang.

TODO: explore using [AFL++](https://github.com/AFLplusplus/AFLplusplus) instead. But note that on GCC its feature set is highly limited (vs Clang).

## Vagrant VM Spin-Up and Modification

**Presently AFL is not available in our Magma Dev Vagrant VM**. You will need to install it yourself from within post boot.

The following instructions come from [here](https://0x00sec.org/t/fuzzing-projects-with-american-fuzzy-lop-afl/6498):


```shell script
cd magma/lte/gateway
vagrant up magma
vagrant ssh magma
git clone https://github.com/mirrorer/afl.git afl
cd afl
make && sudo make install
cd ..
git clone https://github.com/rc0r/afl-utils.git afl-utils
cd afl-utils
sudo python setup.py install
```

## Building the AFL target

```shell script
cd magma/lte/gateway
make clean
CC=/usr/local/bin/afl-gcc CXX=/usr/local/bin/afl-g++ make
```

This will build the output binary at e.g. `/home/vagrant/build/c/oai/tasks/nas/attach_accept_fuzz`.

## Running AFL Against Target

```shell script
afl-fuzz -d -m 81000000 \
  -i lte/gateway/c/oai/tasks/nas/emm/msg/AttachAcceptFuzzerInputs/ \
  -o findings_dir /home/vagrant/build/c/oai/tasks/nas/attach_accept_fuzz
```

The absurdly large memory allowance (800 TB) is because ASAN on 64 bit consumes nuts virtual memory allocations and AFL will kill the process without allowance. See [Doc](https://afl-1.readthedocs.io/en/latest/notes_for_asan.html) from AFL on this issue. The risk is that runaway memory consumption (non ASAN) won't be detected by the fuzzer.

Go grab a coffee. You are working hard. You deserve it.

## Running AFL with Parallelism

```shell script
afl-fuzz -m 81000000 \
  -i lte/gateway/c/oai/tasks/nas/emm/msg/AttachAcceptFuzzerInputs/ \
  -o sync_dir -M fuzzer01 \
  /home/vagrant/build/c/oai/tasks/nas/attach_accept_fuzz
```

Then from another VM terminal (e.g. tmux session):

```shell script
afl-fuzz -m 81000000 \
  -i lte/gateway/c/oai/tasks/nas/emm/msg/AttachAcceptFuzzerInputs/ \
  -o sync_dir -S fuzzer02 \
  /home/vagrant/build/c/oai/tasks/nas/attach_accept_fuzz
```

## Investigating a Crash

```shell script
gdb /home/vagrant/build/c/oai/tasks/nas/attach_accept_fuzz
(gdb) r findings_dir/crashes/id:000000,sig:06,src:000000,op:int32,pos:2,val:be:+1024
(gdb) backtrace
```

# TODO

## Alternative Fuzzers

Explore [AFL++](https://aflplus.plus/docs/tutorials/libxml2_tutorial/) - is this Clang only?

## Performance Improvements

Review [this](https://barro.github.io/2018/06/afl-fuzz-on-different-file-systems/) article on performance improvements for AFL.
