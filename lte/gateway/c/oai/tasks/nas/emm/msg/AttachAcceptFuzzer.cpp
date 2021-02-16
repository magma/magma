#include <cstddef>
#include <cstdint>
#include <cstdlib>
#include <iostream>
#include <sys/types.h>
#include <sys/stat.h>

#include <fcntl.h>
#include <stdio.h>
#include <unistd.h>

#define FUZZER_LOG "FUZZER_LOG"

extern "C" {
#include "AttachAccept.h"
#include "assertions.h"
#include "log.h"
#include "shared_ts_log.h"
}

extern "C" int LLVMFuzzerTestOneInput(const uint8_t* Data, size_t Size) {
  attach_accept_msg result;
  int success = decode_attach_accept(
      &result, const_cast<uint8_t*>(Data), static_cast<uint32_t>(Size));
  return 0;  // Non-zero return values are reserved for future use.
}

int main(int argc, char** argv) {
  CHECK_INIT_RETURN(
      OAILOG_INIT(FUZZER_LOG, OAILOG_LEVEL_DEBUG, MAX_LOG_PROTOS));

  uint8_t input[128];

  ssize_t length = 0;
  if (argc == 2) {
    int input_fd = open(argv[1], O_RDONLY);
    length       = read(input_fd, input, sizeof(input) / sizeof(input[0]));
    std::cout << "Gathered " << length << " bytes from file with contents '"
              << input << "'" << std::endl;
  } else {
    int input_fd = fileno(stdin);
    length       = read(input_fd, input, sizeof(input) / sizeof(input[0]));
    std::cout << "Gathered " << length
              << " bytes from command line with contents '" << input << "'"
              << std::endl;
  }
  LLVMFuzzerTestOneInput(input, length);
}
