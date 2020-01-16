
#pragma once

#include <boost/thread/mutex.hpp>
#include <devmand/channels/cli/Cli.h>
#include <devmand/channels/cli/CliFlavour.h>
#include <devmand/channels/cli/SshSessionAsync.h>
#include <folly/futures/Future.h>

namespace devmand {
namespace channels {
namespace cli {

using folly::Executor;
using folly::SemiFuture;
using folly::Unit;
using std::function;
using std::shared_ptr;
using std::string;

class LoggingCli : public Cli {
 private:
  string id;
  shared_ptr<Cli> cli;
  shared_ptr<folly::Executor> executor;

  SemiFuture<string> executeSomething(
      const Command& cmd,
      const string&& loggingPrefix,
      const function<SemiFuture<string>(shared_ptr<Cli> cli)>& innerFunc);

 public:
  LoggingCli(
      string _id,
      const shared_ptr<Cli>& _cli,
      const shared_ptr<folly::Executor> executor);

  SemiFuture<folly::Unit> destroy() override;

  ~LoggingCli() override;

  folly::SemiFuture<std::string> executeRead(const ReadCommand cmd) override;

  folly::SemiFuture<std::string> executeWrite(const WriteCommand cmd) override;
};

} // namespace cli
} // namespace channels
} // namespace devmand
