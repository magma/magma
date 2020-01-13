#define LOG_WITH_GLOG

#include <magma_logging.h>

#include <devmand/channels/cli/Cli.h>
#include <devmand/channels/cli/TreeCacheCli.h>
#include <devmand/test/TestUtils.h>
#include <devmand/test/cli/TreeCacheTestData.h>
#include <devmand/test/cli/utils/Log.h>
#include <gtest/gtest.h>
#include <magma_logging.h>

namespace devmand::channels::cli {
using namespace std;

static const char* showRunningCommand = "show running-config";

static const shared_ptr<CliFlavour> ubiquitiFlavour =
    CliFlavour::create(UBIQUITI);

static const shared_ptr<CliFlavour> ciscoFlavour = CliFlavour::create("cisco");

class TreeCacheTest : public ::testing::Test {
 protected:
  unique_ptr<TreeCache> tested_ubiquiti;
  unique_ptr<TreeCache> tested_cisco;

  void SetUp() override {
    devmand::test::utils::log::initLog();

    tested_ubiquiti = make_unique<TreeCache>(ubiquitiFlavour);
    tested_ubiquiti->update(testdata::SH_RUN_UBIQUITI);
    EXPECT_EQ(15, tested_ubiquiti->size());

    tested_cisco = make_unique<TreeCache>(ciscoFlavour);
    tested_cisco->update(testdata::SH_RUN_CISCO);
    EXPECT_EQ(3, tested_cisco->size());
  }
};

TEST_F(TreeCacheTest, createSectionPattern_noIndent) {
  auto content_chars = R"template(section1
foo
exit
section2
bar
baz
exit
)template";
  string content(content_chars);

  string sectionPattern = tested_ubiquiti->createSectionPattern(0);
  EXPECT_EQ("\n(\\S[^]*?\nexit\n)", sectionPattern);
  regex regexMatch = regex(sectionPattern);
  smatch sm;

  EXPECT_TRUE(tested_ubiquiti->hasNextSection(regexMatch, sm, content));
  EXPECT_EQ(sm[1], "section1\nfoo\nexit\n");

  EXPECT_TRUE(tested_ubiquiti->hasNextSection(regexMatch, sm, content));
  EXPECT_EQ(sm[1], "section2\nbar\nbaz\nexit\n");

  EXPECT_FALSE(tested_ubiquiti->hasNextSection(regexMatch, sm, content));
}

TEST_F(TreeCacheTest, createSectionPattern_noNewLineAtStart) {
  string content = "section1\nfoo\nexit\n";
  regex regexMatch = regex(tested_ubiquiti->createSectionPattern(0));
  smatch sm;

  EXPECT_TRUE(tested_ubiquiti->hasNextSection(regexMatch, sm, content));
  EXPECT_EQ(sm[1], "section1\nfoo\nexit\n");
  EXPECT_FALSE(tested_ubiquiti->hasNextSection(regexMatch, sm, content));
}

TEST_F(TreeCacheTest, createSectionPattern_noNewLine) {
  string content = "section1";
  regex regexMatch = regex(tested_ubiquiti->createSectionPattern(0));
  smatch sm;

  EXPECT_FALSE(tested_ubiquiti->hasNextSection(regexMatch, sm, content));
}

TEST_F(TreeCacheTest, createSectionPattern_indent) {
  string content(testdata::SH_RUN_TWO_IFC);

  string sectionPattern = tested_ubiquiti->createSectionPattern(' ', "!", 0);

  regex regexMatch = regex(sectionPattern);
  smatch sm;

  EXPECT_TRUE(tested_ubiquiti->hasNextSection(regexMatch, sm, content));
  EXPECT_EQ(
      sm[1],
      "interface Loopback99\n"
      " description bla\n"
      "!\n");
  EXPECT_TRUE(tested_ubiquiti->hasNextSection(regexMatch, sm, content));
  EXPECT_EQ(
      sm[1],
      "interface Bundle-Ether103.100\n"
      " description TOOL_TEST\n"
      " ethernet cfm\n"
      "  mep domain DML3 service 504 mep-id 1\n"
      "   cos 6\n"
      "  !\n"
      " !\n"
      "!\n");
  EXPECT_FALSE(tested_ubiquiti->hasNextSection(regexMatch, sm, content));
  EXPECT_FALSE(tested_ubiquiti->hasNextSection(regexMatch, sm, content));
}

TEST_F(TreeCacheTest, readConfigurationToMap_ubiquiti) {
  string content(testdata::SH_RUN_INT_GI4);
  map<vector<string>, string> actual =
      tested_ubiquiti->readConfigurationToMap(content);
  EXPECT_EQ(actual.size(), 1);
  map<vector<string>, string> expected;

  vector<string> key = {"interface", "0/14"};
  string value =
      "interface 0/14\n"
      "description 'Ruckus-AP'\n"
      "switchport mode access\n"
      "switch access vlan 100\n"
      "exit\n";
  expected.insert(make_pair(key, value));
  EXPECT_EQ(actual, expected);
}

TEST_F(TreeCacheTest, readConfigurationToMap_cisco) {
  string content(testdata::SH_RUN_TWO_IFC);
  map<vector<string>, string> actual =
      tested_cisco->readConfigurationToMap(content);
  EXPECT_EQ(actual.size(), 2);

  map<vector<string>, string> expected;
  {
    vector<string> key = {"interface", "Loopback99"};
    string value =
        "interface Loopback99\n"
        " description bla\n"
        "!\n";
    expected.insert(make_pair(key, value));
  }
  {
    vector<string> key = {"interface", "Bundle-Ether103.100"};
    string value =
        "interface Bundle-Ether103.100\n"
        " description TOOL_TEST\n"
        " ethernet cfm\n"
        "  mep domain DML3 service 504 mep-id 1\n"
        "   cos 6\n"
        "  !\n"
        " !\n"
        "!\n";
    expected.insert(make_pair(key, value));
  }
  EXPECT_EQ(actual, expected);
}

TEST_F(TreeCacheTest, splitSupportedCommand) {
  Optional<pair<string, string>> split = tested_ubiquiti->splitSupportedCommand(
      "show running-config interface 0/14");
  EXPECT_TRUE(split);
  EXPECT_EQ("show running-config", split->first);
  EXPECT_EQ("interface 0/14", split->second);
}

// public api tests:

TEST_F(TreeCacheTest, parseUnknownCommand) {
  auto nothing = tested_ubiquiti->parseCommand("foo");
  EXPECT_FALSE(nothing);
}

TEST_F(TreeCacheTest, getWholeConfig_expectException) {
  Optional<pair<string, vector<string>>> maybeCmd =
      tested_ubiquiti->parseCommand(showRunningCommand);
  EXPECT_TRUE(maybeCmd);
  EXPECT_EQ(maybeCmd->second.size(), 0);
  EXPECT_THROW(tested_ubiquiti->getSection(maybeCmd.value()), runtime_error);
}

TEST_F(TreeCacheTest, getParticularIfc) {
  Optional<pair<string, vector<string>>> maybeCmd =
      tested_ubiquiti->parseCommand("sh run interface 0/14");
  EXPECT_TRUE(maybeCmd);
  EXPECT_EQ(maybeCmd->second.size(), 2);
  Optional<string> wholeConfig = tested_ubiquiti->getSection(maybeCmd.value());
  EXPECT_TRUE(wholeConfig);
  EXPECT_EQ(testdata::SH_RUN_INT_GI4, wholeConfig.value());
}

TEST_F(TreeCacheTest, parseSH_RUN_UBNT_REAL) {
  tested_ubiquiti->clear();
  tested_ubiquiti->update(testdata::SH_RUN_UBNT_REAL);
  EXPECT_EQ(tested_ubiquiti->size(), 11);
  EXPECT_EQ(
      tested_ubiquiti->toString(),
      "(11)["
      "{!},"
      "{!Current}{Configuration:},"
      "{interface}{0/10},"
      "{interface}{0/11},"
      "{interface}{0/7},"
      "{interface}{0/8},"
      "{interface}{0/9},"
      "{interface}{vlan}{33},"
      "{ip}{ssh}{protocol}{1}{2},"
      "{line}{ssh},{line}{telnet},"
      "]");
}

} // namespace devmand::channels::cli
