#include <devmand/devices/cli/ParsingUtils.h>

namespace devmand {
namespace devices {
namespace cli {

using namespace std;

function<ydk::uint64(string)> toUI64 = [](auto s) { return stoull(s); };
function<ydk::uint16(string)> toUI16 = [](auto s) { return stoi(s); };

folly::Optional<string> extractValue(
    const string& output,
    const regex& pattern,
    const uint& groupToExtract) {
  std::stringstream ss(output);
  std::string line;

  while (std::getline(ss, line, '\n')) {
    boost::algorithm::trim(line);
    std::smatch match;
    if (std::regex_match(line, match, pattern) and
        match.size() > groupToExtract and match[groupToExtract].length() > 0) {
      return folly::Optional<string>(match[groupToExtract]);
    }
  }

  return folly::Optional<string>();
}

void parseValue(
    const string& output,
    const regex& pattern,
    const uint& groupToExtract,
    const std::function<void(string)>& setter) {
  const folly::Optional<string>& optValue =
      extractValue(output, pattern, groupToExtract);
  if (optValue) {
    setter(optValue.value());
  }
}

} // namespace cli
} // namespace devices
}