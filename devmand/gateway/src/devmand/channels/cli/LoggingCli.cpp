#include <devmand/channels/cli/LoggingCli.h>
#include <magma_logging.h>

namespace devmand::channels::cli {
using namespace folly;
using namespace std;

LoggingCli::LoggingCli(
    string _id,
    const shared_ptr<Cli>& _cli,
    const shared_ptr<Executor> _executor)
    : id(_id), cli(_cli), executor(_executor) {}

SemiFuture<Unit> LoggingCli::destroy() {
  MLOG(MDEBUG) << "[" << id << "] "
               << "destroy: started";
  // call underlying destroy()
  SemiFuture<Unit> innerDestroy = cli->destroy();
  MLOG(MDEBUG) << "[" << id << "] "
               << "destroy: done";
  return innerDestroy;
}

LoggingCli::~LoggingCli() {
  MLOG(MDEBUG) << "[" << id << "] "
               << "~Lcli: started";
  destroy().get();
  executor = nullptr;
  cli = nullptr;
  MLOG(MDEBUG) << "[" << id << "] "
               << "~Lcli: done";
}

SemiFuture<string> LoggingCli::executeRead(const ReadCommand cmd) {
  return executeSomething(cmd, "LCli.executeRead", [cmd](shared_ptr<Cli> _cli) {
    return _cli->executeRead(cmd);
  });
}

SemiFuture<string> LoggingCli::executeWrite(const WriteCommand cmd) {
  return executeSomething(
      cmd, "LCli.executeWrite", [cmd](shared_ptr<Cli> _cli) {
        return _cli->executeWrite(cmd);
      });
}

SemiFuture<string> LoggingCli::executeSomething(
    const Command& cmd,
    const string&& loggingPrefix,
    const function<SemiFuture<string>(shared_ptr<Cli> cli)>& innerFunc) {
  MLOG(MDEBUG) << "[" << id << "] (" << cmd << ") " << loggingPrefix << "('"
               << cmd << "') called";
  SemiFuture<string> inner =
      innerFunc(cli); // we expect that this method does not block
  MLOG(MDEBUG) << "[" << id << "] (" << cmd << ") "
               << "Obtained future from underlying cli";
  return move(inner)
      .via(executor.get())
      .thenTry([id = this->id, cmd](const Try<string>&& t) {
        // Logging callback
        if (t.hasException()) {
          MLOG(MWARNING) << "[" << id << "] "
                         << "Command: \"" << cmd << "\""
                         << " Failed with: " << t.exception().what();
        } else {
          MLOG(MINFO) << "[" << id << "] "
                      << "Command: \"" << cmd << "\""
                      << " Invoked successfully with output: \""
                      << Command::escape(t.value()) << "\"";
        }
        return t;
      })
      .semi();
}
} // namespace devmand::channels::cli
