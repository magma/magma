// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <boost/algorithm/string/split.hpp>
#include <boost/algorithm/string/trim.hpp>
#include <devmand/channels/cli/CliFlavour.h>
#include <regex>

namespace devmand {
namespace channels {
namespace cli {

using namespace std;
using devmand::channels::cli::CliFlavour;
using devmand::channels::cli::CliInitializer;
using devmand::channels::cli::DefaultPromptResolver;
using devmand::channels::cli::EmptyInitializer;
using devmand::channels::cli::PromptResolver;
using devmand::channels::cli::UbiquitiInitializer;
using devmand::channels::cli::sshsession::SessionAsync;
using folly::Optional;

SemiFuture<Unit> EmptyInitializer::initialize(
    shared_ptr<SessionAsync> session,
    string secret) {
  (void)session;
  (void)secret;
  return folly::makeFuture();
}

SemiFuture<Unit> UbiquitiInitializer::initialize(
    shared_ptr<SessionAsync> session,
    string secret) {
  return session->write("enable\n")
      .thenValue(
          [session, secret](...) { return session->write(secret + "\n"); })
      .thenValue(
          [session](...) { return session->write("terminal length 0\n"); });
}

Future<string> DefaultPromptResolver::resolvePrompt(
    shared_ptr<SessionAsync> session,
    const string& newline,
    shared_ptr<Timekeeper> timekeeper) {
  return session->read().thenValue([=](...) {
    return resolvePrompt(session, newline, delayDelta, timekeeper);
  });
}

Future<string> DefaultPromptResolver::resolvePrompt(
    shared_ptr<SessionAsync> session,
    const string& newline,
    chrono::milliseconds delay,
    shared_ptr<Timekeeper> timekeeper) {
  return resolvePromptAsync(session, newline, delay, timekeeper)
      .thenValue([=](Optional<string> prompt) {
        if (!prompt.hasValue()) {
          return resolvePrompt(
              session, newline, delay + delayDelta, timekeeper);
        } else {
          return folly::makeFuture(prompt.value());
        }
      });
}

Future<Optional<string>> DefaultPromptResolver::resolvePromptAsync(
    shared_ptr<SessionAsync> session,
    const string& newline,
    chrono::milliseconds delay,
    shared_ptr<Timekeeper> timekeeper) {
  return session->write(newline + newline)
      .delayed(delay, timekeeper.get())
      .thenValue([session](...) { return session->read(); })
      .thenValue([=](string output) {
        regex regxp("\\" + newline);
        vector<string> split(
            sregex_token_iterator(output.begin(), output.end(), regxp, -1),
            sregex_token_iterator());

        removeEmptyStrings(split);

        if (split.size() == 2) {
          string s0 = boost::algorithm::trim_copy(split[0]);
          string s1 = boost::algorithm::trim_copy(split[1]);
          if (s0 == s1) {
            return folly::make_optional<string>(s0);
          }
        }
        return Optional<string>();
      });
}

void DefaultPromptResolver::removeEmptyStrings(vector<string>& split) const {
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
    unique_ptr<PromptResolver>&& _resolver,
    unique_ptr<CliInitializer>&& _initializer,
    string _newline,
    regex _baseShowConfig,
    unsigned int _baseShowConfigIdx,
    Optional<char> _singleIndentChar,
    string _configSubsectionEnd)
    : resolver(forward<unique_ptr<PromptResolver>>(_resolver)),
      initializer(forward<unique_ptr<CliInitializer>>(_initializer)),
      newline(_newline),
      baseShowConfig(_baseShowConfig),
      baseShowConfigIdx(_baseShowConfigIdx),
      singleIndentChar(_singleIndentChar),
      configSubsectionEnd(_configSubsectionEnd) {}

shared_ptr<CliFlavour> CliFlavour::create(string flavour) {
  if (flavour == UBIQUITI) {
    return make_shared<CliFlavour>(
        make_unique<DefaultPromptResolver>(),
        make_unique<UbiquitiInitializer>(),
        "\n",
        regex(R"(^((do )?sho?w? runn?i?n?g?-?c?o?n?f?i?g?).*)"),
        1,
        none,
        "exit");
  } else {
    return make_shared<CliFlavour>(
        make_unique<DefaultPromptResolver>(),
        make_unique<EmptyInitializer>(),
        "\n",
        regex(R"(^((do )?sho?w? runn?i?n?g?-?c?o?n?f?i?g?).*)"),
        1,
        ' ',
        "!");
  }
}

Optional<size_t> CliFlavour::getBaseShowConfigIdx(const string cmd) const {
  smatch pieces_match;
  if (regex_match(cmd, pieces_match, baseShowConfig)) {
    return Optional<size_t>((size_t)pieces_match[baseShowConfigIdx].length());
  }
  return Optional<size_t>(none);
}

Optional<char> CliFlavour::getSingleIndentChar() {
  return singleIndentChar;
}

string CliFlavour::getConfigSubsectionEnd() {
  return configSubsectionEnd;
}

vector<string> CliFlavour::splitSubcommands(string subcommands) {
  Optional<char> maybeIndentChar = getSingleIndentChar();
  char indentChar = maybeIndentChar.value_or(' ');
  vector<string> args;
  boost::split(
      args, subcommands, [indentChar](char s) { return s == indentChar; });
  return args;
}

shared_ptr<PromptResolver> CliFlavour::getResolver() {
  return resolver;
}

shared_ptr<CliInitializer> CliFlavour::getInitializer() {
  return initializer;
}

string CliFlavour::getNewline() {
  return newline;
}

} // namespace cli
} // namespace channels
} // namespace devmand
