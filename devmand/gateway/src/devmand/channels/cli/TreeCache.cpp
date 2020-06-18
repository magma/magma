// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/cli/TreeCache.h>

namespace devmand::channels::cli {

// public api start

TreeCache::TreeCache(shared_ptr<CliFlavour> _cliFlavour)
    : cliFlavour(_cliFlavour) {}

void TreeCache::clear() {
  treeCache.clear();
}

Optional<pair<string, vector<string>>> TreeCache::parseCommand(string cmd) {
  Optional<pair<string, string>> split = splitSupportedCommand(cmd);
  if (split) {
    vector<string> subcommands;
    if (split.value().second.size() > 0) {
      subcommands = cliFlavour->splitSubcommands(split.value().second);
    }
    return make_pair(split.value().first, subcommands);
  }
  return none;
}

void TreeCache::update(string output) {
  boost::replace_all(output, "\r", "");
  treeCache = readConfigurationToMap(output, 0);
}

Optional<string> TreeCache::getSection(pair<string, vector<string>> cmd) {
  if (cmd.second.size() == 0) {
    throw runtime_error("Command does not have subcommands");
  }
  return findMatchingSubset(cmd.second, treeCache);
}

bool TreeCache::isEmpty() {
  return treeCache.size() == 0;
}

// public api end

// used in StringUtils::getFirstLine
static const regex firstLineRegex = regex("^([^\n]*)\n");
class StringUtils {
 private:
  // trim from start (in place)
  static inline void ltrim(std::string& s) {
    s.erase(
        s.begin(),
        std::find_if(
            s.begin(),
            s.end(),
            std::not1(std::ptr_fun<int, int>(std::isspace))));
  }

  // trim from end (in place)
  static inline void rtrim(std::string& s) {
    s.erase(
        std::find_if(
            s.rbegin(),
            s.rend(),
            std::not1(std::ptr_fun<int, int>(std::isspace)))
            .base(),
        s.end());
  }

 public:
  // trim from both ends (in place)
  static inline void trim(std::string& s) {
    ltrim(s);
    rtrim(s);
  }

  static string getFirstLine(string input) {
    smatch sm;
    if (regex_search(input, sm, firstLineRegex)) {
      return sm[1];
    } else {
      // no \n, return everything
      return input;
    }
  }
};

size_t TreeCache::size() {
  return treeCache.size();
}

// parsing start
map<vector<string>, string> TreeCache::readConfigurationToMap(
    string showRunningOutput,
    unsigned int indentationLevel) {
  // currently only one level of sections is supported, thus indentationLevel ==
  // 0
  regex regexMatch = regex(createSectionPattern(indentationLevel));
  smatch sm;
  map<vector<string>, string> result;

  while (hasNextSection(regexMatch, sm, showRunningOutput)) {
    string section = sm[1];
    string firstLine = StringUtils::getFirstLine(section);
    vector<string> args = cliFlavour->splitSubcommands(firstLine);
    result.insert(make_pair(args, section));
  }
  // last piece of content that is not a section is ignored
  return result;
}

map<vector<string>, string> TreeCache::readConfigurationToMap(
    string showRunningOutput) {
  return readConfigurationToMap(showRunningOutput, 0);
}

bool TreeCache::hasNextSection(regex& regexMatch, smatch& sm, string& content) {
  if (not sm.empty()) {
    // move content after previous section
    content = "\n";
    content += sm.suffix();
  } else if (content.size() > 0 && content.substr(0, 1) != "\n") {
    content = "\n" + content;
  }
  return regex_search(content, sm, regexMatch);
}

string TreeCache::createSectionPattern(unsigned int indentationLevel) {
  return createSectionPattern(
      cliFlavour->getSingleIndentChar(),
      cliFlavour->getConfigSubsectionEnd(),
      indentationLevel);
}

string TreeCache::createSectionPattern(
    Optional<char> maybeIndentChar,
    string configSubsectionEnd,
    unsigned int indentationLevel) {
  string match = "\\S[^]*?";
  if (maybeIndentChar.hasValue()) {
    // CONFIG_SUBSECTION_PATTERN = ^%s{LVL}\S.*?^%s{LVL}END
    string nlIndent = "";
    nlIndent += maybeIndentChar.value();
    nlIndent += "{" + to_string(indentationLevel) + "}";
    return "\n(" + nlIndent + match + "\n" + nlIndent + configSubsectionEnd +
        "\n)";
  } else {
    // CONFIG_SUBSECTION_NO_INDENT_PATTERN = ^\S.*?^END
    return "\n(" + match + "\n" + configSubsectionEnd + "\n)";
  }
}
// parsing end
// reading start
Optional<pair<string, string>> TreeCache::splitSupportedCommand(string cmd) {
  const Optional<size_t>& maybePosition = cliFlavour->getBaseShowConfigIdx(cmd);
  if (maybePosition.hasValue()) {
    string first = cmd.substr(0, maybePosition.value()); // space removed
    string second = cmd.substr(maybePosition.value());
    StringUtils::trim(second);
    return make_pair(first, second);
  } else {
    // not a show running config command
    return none;
  }
}

Optional<string> TreeCache::findMatchingSubset(
    vector<string> subcommands,
    map<vector<string>, string> treeCache) {
  // simplistic case: find whole key in map
  if (treeCache.count(subcommands) == 1) {
    return treeCache.at(subcommands);
  }
  return none;
}
// reading end
} // namespace devmand::channels::cli
