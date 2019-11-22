#include <ydk/types.hpp>
#include <folly/futures/Future.h>
#include <gtest/gtest.h>
namespace devmand {
namespace test {
namespace cli {
class ModelRegistryTest : public ::testing::Test {};
TEST_F(ModelRegistryTest, caching) {
std::make_unique<std::string>();
}
} // namespace cli
} // namespace test
} // namespace devmand
