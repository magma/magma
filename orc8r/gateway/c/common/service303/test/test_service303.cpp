/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#include <unistd.h>
#include <pthread.h>
#include <prometheus/registry.h>
#include <gtest/gtest.h>
#include <map>
#include <orc8r/protos/metricsd.pb.h>
#include <string>
#include <thread>

#include "orc8r/gateway/c/common/service303/includes/MetricsRegistry.h"
#include "orc8r/gateway/c/common/service303/includes/MetricsSingleton.h"
#include "orc8r/gateway/c/common/service303/includes/MetricsHelpers.h"
#include "orc8r/gateway/c/common/service303/includes/MagmaService.h"
#include "orc8r/gateway/c/common/service_registry/includes/ServiceRegistrySingleton.h"

using grpc::Channel;
using grpc::ChannelCredentials;
using grpc::ClientContext;
using grpc::CreateChannel;
using grpc::InsecureChannelCredentials;
using io::prometheus::client::LabelPair;
using io::prometheus::client::MetricFamily;

using magma::orc8r::MetricsContainer;
using magma::orc8r::ServiceInfo;
using magma::service303::MagmaService;

using ::testing::Test;
namespace magma {

#define NO_LABELS 0
#define NO_BOUNDARIES 0

class Service303Client {
 public:
  explicit Service303Client(const std::shared_ptr<Channel>& channel)
      : stub_(Service303::NewStub(channel)) {}

  /**
   * @param response: a pointer to the ServiceInfo object to populate
   * @return 0 on success, -1 on failure
   */
  int GetServiceInfo(ServiceInfo* response) {
    Void request;
    ClientContext context;
    Status status = stub_->GetServiceInfo(&context, request, response);
    if (!status.ok()) {
      std::cout << "GetServiceInfo fails with code " << status.error_code()
                << ", msg: " << status.error_message() << std::endl;
      return -1;
    }
    return 0;
  }

  /**
   * @param response: the MetricsContainer instance to populate
   * @return 0 on success, -1 on failure
   */
  int GetMetrics(MetricsContainer* response) {
    ClientContext context;
    Void request;
    Status status = stub_->GetMetrics(&context, request, response);
    if (!status.ok()) {
      std::cout << "GetMetrics fails with code " << status.error_code()
                << ", msg: " << status.error_message() << std::endl;
      return -1;
    }
    return 0;
  }

 private:
  std::shared_ptr<Service303::Stub> stub_;
};

class Service303Test : public ::testing::Test {
 public:
  virtual void SetUp() {
    magma_service = std::make_shared<MagmaService>(service_name, version);
    magma_service->Start();
    // Wait for the server to start
    usleep(50000);
    setupClient();
  }

  virtual void TearDown() {
    magma_service->Stop();
    magma_service->WaitForShutdown();
    service303::MetricsSingleton::flush();
    delete service303_client;
  }

  static const MetricFamily& findFamily(const MetricsContainer& container,
                                        const std::string& match) {
    for (auto const& fam : container.family()) {
      if (fam.name().compare(match) == 0) {
        return fam;
      }
    }
    assert(false);
  }

  static const double findGauge(const MetricsContainer& container,
                                const std::string& match) {
    const MetricFamily& fam = findFamily(container, match);
    return fam.metric().Get(0).gauge().value();
  }

 protected:
  const std::string service_name = "test_service";
  const std::string version = "magma@1.2.3";
  std::thread magma_server_thread;
  std::shared_ptr<MagmaService> magma_service = nullptr;
  Service303Client* service303_client;

 private:
  void setupClient() {
    const std::shared_ptr<ChannelCredentials> cred =
        InsecureChannelCredentials();
    std::string service_addr =
        magma::ServiceRegistrySingleton::Instance()->GetServiceAddrString(
            service_name);
    const std::shared_ptr<Channel> channel = CreateChannel(service_addr, cred);
    service303_client = new Service303Client(channel);
  }
};

// Tests against Service303::GetServiceInfo()
TEST_F(Service303Test, test_get_service_info) {
  ServiceInfo service_info;
  int status = service303_client->GetServiceInfo(&service_info);
  EXPECT_EQ(0, status);
  EXPECT_EQ(service_info.name(), service_name);
  EXPECT_EQ(service_info.version(), version);
  EXPECT_EQ(service_info.state(), ServiceInfo::ALIVE);
}

// Tests that Service303 can instrument counters and read them over gRPC.
TEST_F(Service303Test, test_counters) {
  increment_counter("test_counter", 3, NO_LABELS);
  MetricsContainer metrics_container;
  EXPECT_EQ(0, service303_client->GetMetrics(&metrics_container));
  const MetricFamily& family =
      Service303Test::findFamily(metrics_container, "test_counter");
  EXPECT_EQ(family.name(), "test_counter");
  const io::prometheus::client::Counter& counter =
      family.metric().Get(0).counter();
  EXPECT_EQ(counter.value(), 3);
}

// Tests that Service303 can instrument gauges and read them over gRPC.
TEST_F(Service303Test, test_gauges) {
  // Increment gauge with labels
  increment_gauge("test_gauge", 3, 2, "key", "value", "test", "test");
  MetricsContainer metrics_container;
  EXPECT_EQ(0, service303_client->GetMetrics(&metrics_container));
  const MetricFamily* family =
      &Service303Test::findFamily(metrics_container, "test_gauge");
  EXPECT_EQ(family->name(), "test_gauge");
  auto* metric = &family->metric().Get(0);
  EXPECT_EQ(metric->gauge().value(), 3);
  EXPECT_EQ(metric->label().size(), 2);

  // Increment another gauge
  increment_gauge("test_gauge", 3, NO_LABELS);
  EXPECT_EQ(0, service303_client->GetMetrics(&metrics_container));
  family = &Service303Test::findFamily(metrics_container, "test_gauge");
  EXPECT_EQ(family->name(), "test_gauge");
  auto* gauge = &family->metric().Get(0).gauge();
  EXPECT_EQ(gauge->value(), 3);

  // Decrement back to zero
  decrement_gauge("test_gauge", 3, NO_LABELS);
  EXPECT_EQ(0, service303_client->GetMetrics(&metrics_container));
  family = &Service303Test::findFamily(metrics_container, "test_gauge");
  EXPECT_EQ(family->name(), "test_gauge");
  gauge = &family->metric().Get(0).gauge();
  EXPECT_EQ(gauge->value(), 0);

  // Set the gauge to 10
  set_gauge("test_gauge", 10, NO_LABELS);
  EXPECT_EQ(0, service303_client->GetMetrics(&metrics_container));
  family = &Service303Test::findFamily(metrics_container, "test_gauge");
  EXPECT_EQ(family->name(), "test_gauge");
  gauge = &family->metric().Get(0).gauge();
  EXPECT_EQ(gauge->value(), 10);
}

// Tests that Service303 can instrument histograms and read them over gRPC.
TEST_F(Service303Test, test_histograms) {
  // First observation in a histogram without buckets
  observe_histogram("test_hist", 3, NO_LABELS, NO_BOUNDARIES);
  MetricsContainer metrics_container;
  EXPECT_EQ(0, service303_client->GetMetrics(&metrics_container));
  const MetricFamily* family =
      &Service303Test::findFamily(metrics_container, "test_hist");
  EXPECT_EQ(family->name(), "test_hist");
  EXPECT_EQ(family->metric().size(), 1);
  auto* histogram = &family->metric().Get(0).histogram();
  EXPECT_EQ(histogram->sample_count(), 1);
  EXPECT_EQ(histogram->sample_sum(), 3);
  EXPECT_EQ(histogram->bucket().size(), 1);
  EXPECT_EQ(histogram->bucket().Get(0).cumulative_count(), 1);
  EXPECT_EQ(histogram->bucket().Get(0).upper_bound(),
            std::numeric_limits<double>::infinity());

  // Adding another observation with buckets won't add another metric or
  // more buckets but it will add another observation to the metric
  observe_histogram("test_hist", 7, NO_LABELS, 3, 1, 10, 100);
  EXPECT_EQ(0, service303_client->GetMetrics(&metrics_container));
  family = &Service303Test::findFamily(metrics_container, "test_hist");
  EXPECT_EQ(family->name(), "test_hist");
  EXPECT_EQ(family->metric().size(), 1);
  histogram = &family->metric().Get(0).histogram();
  EXPECT_EQ(histogram->sample_count(), 2);
  EXPECT_EQ(histogram->sample_sum(), 10);
  EXPECT_EQ(histogram->bucket().size(), 1);
  EXPECT_EQ(histogram->bucket().Get(0).cumulative_count(), 2);
  EXPECT_EQ(histogram->bucket().Get(0).upper_bound(),
            std::numeric_limits<double>::infinity());

  // Since we define a new metric using labels, it should construct with the
  // boundary definitions
  observe_histogram("test_hist", 50, 1, "key", "value", 3, 1., 10., 100.);
  EXPECT_EQ(0, service303_client->GetMetrics(&metrics_container));
  family = &Service303Test::findFamily(metrics_container, "test_hist");
  EXPECT_EQ(family->metric().size(), 2);
  EXPECT_EQ(family->name(), "test_hist");
  histogram = &family->metric().Get(0).histogram();
  EXPECT_EQ(histogram->sample_count(), 1);
  EXPECT_EQ(histogram->sample_sum(), 50);
  EXPECT_EQ(histogram->bucket().size(), 4);
  EXPECT_EQ(histogram->bucket().Get(0).cumulative_count(), 0);
  EXPECT_EQ(histogram->bucket().Get(0).upper_bound(), 1);
  EXPECT_EQ(histogram->bucket().Get(1).cumulative_count(), 0);
  EXPECT_EQ(histogram->bucket().Get(1).upper_bound(), 10);
  EXPECT_EQ(histogram->bucket().Get(2).cumulative_count(), 1);
  EXPECT_EQ(histogram->bucket().Get(2).upper_bound(), 100);
  EXPECT_EQ(histogram->bucket().Get(3).cumulative_count(), 0);
  EXPECT_EQ(histogram->bucket().Get(3).upper_bound(),
            std::numeric_limits<double>::infinity());
}

TEST_F(Service303Test, test_timing_metrics) {
  MetricsContainer metrics_container;
  EXPECT_EQ(0, service303_client->GetMetrics(&metrics_container));
  double start_time = Service303Test::findGauge(
      metrics_container,
      std::to_string(magma::orc8r::MetricName::process_start_time_seconds));
  double uptime = Service303Test::findGauge(
      metrics_container,
      std::to_string(magma::orc8r::MetricName::process_cpu_seconds_total));

  EXPECT_GT(start_time, 0);
  EXPECT_GE(uptime, 0);

  sleep(1);
  EXPECT_EQ(0, service303_client->GetMetrics(&metrics_container));
  double new_start = Service303Test::findGauge(
      metrics_container,
      std::to_string(magma::orc8r::MetricName::process_start_time_seconds));
  double new_uptime = Service303Test::findGauge(
      metrics_container,
      std::to_string(magma::orc8r::MetricName::process_cpu_seconds_total));

  EXPECT_DOUBLE_EQ(start_time, new_start);
  EXPECT_GT(new_uptime, uptime);
}

TEST_F(Service303Test, test_memory_metrics) {
  MetricsContainer metrics_container;
  EXPECT_EQ(0, service303_client->GetMetrics(&metrics_container));
  double physical_mem = Service303Test::findGauge(
      metrics_container,
      std::to_string(magma::orc8r::MetricName::process_resident_memory_bytes));
  double virtual_mem = Service303Test::findGauge(
      metrics_container,
      std::to_string(magma::orc8r::MetricName::process_virtual_memory_bytes));

  // Sanity check memory metrics
  EXPECT_GT(physical_mem, 0);
  EXPECT_GT(virtual_mem, 0);
}

// test that metric names and labels can be converted to their enum form
TEST_F(Service303Test, test_enum_conversions) {
  // enum values (from metricsd.proto):
  // s1_setup => 500, result => 0
  increment_counter("mme_new_association", 3, 2, "result", "success", "gateway",
                    "1234");

  MetricsContainer metrics_container;
  EXPECT_EQ(0, service303_client->GetMetrics(&metrics_container));
  const MetricFamily& family =
      Service303Test::findFamily(metrics_container, "500");
  const io::prometheus::client::Counter& counter =
      family.metric().Get(0).counter();
  EXPECT_EQ(counter.value(), 3);

  const LabelPair& first_label = family.metric().Get(0).label().Get(0);
  // test converted enum
  EXPECT_EQ(first_label.name(), "0");
  EXPECT_EQ(first_label.value(), "success");

  const LabelPair& second_label = family.metric().Get(0).label().Get(1);
  // test non converted enum
  EXPECT_EQ(second_label.name(), "gateway");
  EXPECT_EQ(second_label.value(), "1234");
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}
}  // namespace magma
