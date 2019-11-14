// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <boost/algorithm/string/trim.hpp>
#include <devmand/channels/cli/CliFlavour.h>
#include <regex>

namespace devmand {
namespace channels {
namespace cli {

using devmand::channels::cli::CliFlavour;
using devmand::channels::cli::CliInitializer;
using devmand::channels::cli::PromptResolver;
using devmand::channels::cli::sshsession::SshSessionAsync;
using devmand::channels::cli::EmptyInitializer;
using devmand::channels::cli::UbiquitiInitializer;
using devmand::channels::cli::DefaultPromptResolver;

static const int DEFAULT_MILLIS = 1000;

void EmptyInitializer::initialize(shared_ptr<SshSessionAsync> session) {
    (void) session;
}

void UbiquitiInitializer::initialize(shared_ptr<SshSessionAsync> session) {
  session->write("enable\n")
      .thenValue([=](...) { return session->write("ubnt\n"); })
      .thenValue([=](...) { return session->write("terminal length 0\n"); })
      .get();
}

string DefaultPromptResolver::resolvePrompt(
    shared_ptr<SshSessionAsync> session,
    const string& newline) {
  session->read(DEFAULT_MILLIS).get(); // clear input, converges faster on
                                       // prompt
  for (int i = 1;; i++) {
    int millis = i * DEFAULT_MILLIS;
    session->write(newline + newline).get();
    string output = session->read(millis).get();

    std::regex regxp("\\" + newline);
    std::vector<string> split(
        std::sregex_token_iterator(output.begin(), output.end(), regxp, -1),
        std::sregex_token_iterator());

    removeEmptyStrings(split);

    if (split.size() == 2) {
      string s0 = boost::algorithm::trim_copy(split[0]);
      string s1 = boost::algorithm::trim_copy(split[1]);
      if (s0 == s1) {
        return s0;
      }
    }
  }
}

void DefaultPromptResolver::removeEmptyStrings(std::vector<string>& split) const {
  split.erase(
      remove_if(
          split.begin(),
          split.end(),
          [](string& el) {
            boost::algorithm::trim(el);
            return el.empty();
          }),
      split.end());
}

CliFlavour::CliFlavour(
    PromptResolver * _resolver,
    CliInitializer * _initializer,
    string _newline)
    : resolver(_resolver), initializer(_initializer), newline(_newline) {}

CliFlavour::~CliFlavour() {
    delete resolver;
    delete initializer;
}

    shared_ptr<CliFlavour> CliFlavour::create(string flavour) {
    if (flavour == UBIQUITI) {
        return std::make_shared<CliFlavour>(new DefaultPromptResolver(), new UbiquitiInitializer());
    }

    return std::make_shared<CliFlavour>(new DefaultPromptResolver(), new EmptyInitializer());
}

} // namespace cli
} // namespace channels
} // namespace devmand