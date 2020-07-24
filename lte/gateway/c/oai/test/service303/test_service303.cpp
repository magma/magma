/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */
#include <map>
#include <pthread.h>
#include <string>
#include <unistd.h>

#include <gtest/gtest.h>
#include <prometheus/registry.h>
#include <prometheus/metrics.pb.h>
#include "orc8r/protos/metricsd.pb.h"

#include "bstrlib.h"
#include "MetricsRegistry.h"
#include "MetricsSingleton.h"
#include "service303.h"
#include "Service303Client.h"
#include "ServiceRegistrySingleton.h"

using grpc::Channel;
using grpc::ChannelCredentials;
using grpc::CreateChannel;
using grpc::InsecureChannelCredentials;
using io::prometheus::client::LabelPair;
using io::prometheus::client::MetricFamily;
using magma::Service303Client;
using magma::orc8r::ServiceInfo;
using magma::service303::MetricsRegistry;
using prometheus::BuildCounter;
using prometheus::Registry;
using prometheus::detail::CounterBuilder;

namespace {

// The fixture for testing class Service303.
class Service303Test : public ::testing::Test {
 public:
  static service303_data_t config;
  static Service303Client* client;

  virtual void SetUp() {
    // Code here will be called immediately after the constructor (right
    // before each test).

    auto thread_init = [](void* arg) -> void* {
      service303_data_t* service303_data = (service303_data_t*) arg;
      start_service303_server(service303_data->name, service303_data->version);
      pthread_exit(NULL);
    };
    pthread_t server_thread;
    int rc = pthread_create(
        &server_thread, NULL, thread_init, (void*) &Service303Test::config);
    if (rc) {
      std::cout << "Error:unable to create thread," << rc << std::endl;
      exit(-1);
    }
    // Wait for the server to start
    usleep(50000);
    setupClient();
  }

  virtual void TearDown() {
    // Code here will be called immediately after each test (right
    // before the destructor).
    stop_service303_server();
    magma::service303::MetricsSingleton::flush();
    delete client;
  }

  static const MetricFamily& findFamily(
      const MetricsContainer& container, const std::string& match) {
    for (auto const& fam : container.family()) {
      if (fam.name().compare(match) == 0) {
        return fam;
      }
    }
    assert(false);
  }

  static const double findGauge(
      const MetricsContainer& container, const std::string& match) {
    const MetricFamily& fam = findFamily(container, match);
    return fam.metric().Get(0).gauge().value();
  }

 private:
  static void setupClient() {
    const std::shared_ptr<ChannelCredentials> cred =
        InsecureChannelCredentials();
    std::string service_addr =
        magma::ServiceRegistrySingleton::Instance()->GetServiceAddrString(
            bdata(config.name));
    const std::shared_ptr<Channel> channel = CreateChannel(service_addr, cred);
    client                                 = new Service303Client(channel);
  }
};
// forward declaration and initializing of static member
service303_data_t Service303Test::config;
Service303Client* Service303Test::client = NULL;

}  // namespace

// Tests against Service303::GetServiceInfo()
TEST_F(Service303Test, GetServiceInfo) {
  Service303Client client = *Service303Test::client;

  ServiceInfo service_info;
  int status = client.GetServiceInfo(&service_info);
  EXPECT_EQ(0, status);
  char* name    = bdata(Service303Test::config.name);
  char* version = bdata(Service303Test::config.version);
  EXPECT_EQ(service_info.name(), name);
  EXPECT_EQ(service_info.version(), version);
  EXPECT_EQ(service_info.state(), ServiceInfo::ALIVE);
}

// Tests that Service303 can instrument counters and read them over gRPC.
TEST_F(Service303Test, TestCounters) {
  Service303Client client = *Service303Test::client;

  increment_counter("test_counter", 3, NO_LABELS);
  MetricsContainer metrics_container;
  int status = client.GetMetrics(&metrics_container);
  const MetricFamily& family =
      Service303Test::findFamily(metrics_container, "test_counter");
  EXPECT_EQ(family.name(), "test_counter");
  const io::prometheus::client::Counter& counter =
      family.metric().Get(0).counter();
  EXPECT_EQ(counter.value(), 3);
}

// Tests that Service303 can instrument gauges and read them over gRPC.
TEST_F(Service303Test, TestGauges) {
  Service303Client client = *Service303Test::client;

  // Increment gauge with labels
  increment_gauge("test_gauge", 3, 2, "key", "value", "test", "test");
  MetricsContainer metrics_container;
  int status = client.GetMetrics(&metrics_container);
  const MetricFamily* family =
      &Service303Test::findFamily(metrics_container, "test_gauge");
  EXPECT_EQ(family->name(), "test_gauge");
  auto* metric = &family->metric().Get(0);
  EXPECT_EQ(metric->gauge().value(), 3);
  EXPECT_EQ(metric->label().size(), 2);

  // Increment another gauge
  increment_gauge("test_gauge", 3, NO_LABELS);
  status = client.GetMetrics(&metrics_container);
  family = &Service303Test::findFamily(metrics_container, "test_gauge");
  EXPECT_EQ(family->name(), "test_gauge");
  auto* gauge = &family->metric().Get(0).gauge();
  EXPECT_EQ(gauge->value(), 3);

  // Decrement back to zero
  decrement_gauge("test_gauge", 3, NO_LABELS);
  status = client.GetMetrics(&metrics_container);
  family = &Service303Test::findFamily(metrics_container, "test_gauge");
  EXPECT_EQ(family->name(), "test_gauge");
  gauge = &family->metric().Get(0).gauge();
  EXPECT_EQ(gauge->value(), 0);

  // Set the gauge to 10
  set_gauge("test_gauge", 10, NO_LABELS);
  status = client.GetMetrics(&metrics_container);
  family = &Service303Test::findFamily(metrics_container, "test_gauge");
  EXPECT_EQ(family->name(), "test_gauge");
  gauge = &family->metric().Get(0).gauge();
  EXPECT_EQ(gauge->value(), 10);
}

// Tests that Service303 can instrument histograms and read them over gRPC.
TEST_F(Service303Test, TestHistograms) {
  Service303Client client = *Service303Test::client;

  // First observation in a histogram without buckets
  observe_histogram("test_hist", 3, NO_LABELS, NO_BOUNDARIES);
  MetricsContainer metrics_container;
  int status = client.GetMetrics(&metrics_container);
  const MetricFamily* family =
      &Service303Test::findFamily(metrics_container, "test_hist");
  EXPECT_EQ(family->name(), "test_hist");
  EXPECT_EQ(family->metric().size(), 1);
  auto* histogram = &family->metric().Get(0).histogram();
  EXPECT_EQ(histogram->sample_count(), 1);
  EXPECT_EQ(histogram->sample_sum(), 3);
  EXPECT_EQ(histogram->bucket().size(), 1);
  EXPECT_EQ(histogram->bucket().Get(0).cumulative_count(), 1);
  EXPECT_EQ(
      histogram->bucket().Get(0).upper_bound(),
      std::numeric_limits<double>::infinity());

  // Adding another observation with buckets won't add another metric or
  // more buckets but it will add another observation to the metric
  observe_histogram("test_hist", 7, NO_LABELS, 3, 1, 10, 100);
  status = client.GetMetrics(&metrics_container);
  family = &Service303Test::findFamily(metrics_container, "test_hist");
  EXPECT_EQ(family->name(), "test_hist");
  EXPECT_EQ(family->metric().size(), 1);
  histogram = &family->metric().Get(0).histogram();
  EXPECT_EQ(histogram->sample_count(), 2);
  EXPECT_EQ(histogram->sample_sum(), 10);
  EXPECT_EQ(histogram->bucket().size(), 1);
  EXPECT_EQ(histogram->bucket().Get(0).cumulative_count(), 2);
  EXPECT_EQ(
      histogram->bucket().Get(0).upper_bound(),
      std::numeric_limits<double>::infinity());

  // Since we define a new metric using labels, it should construct with the
  // boundary definitions
  observe_histogram("test_hist", 50, 1, "key", "value", 3, 1., 10., 100.);
  status = client.GetMetrics(&metrics_container);
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
  EXPECT_EQ(
      histogram->bucket().Get(3).upper_bound(),
      std::numeric_limits<double>::infinity());
}

TEST_F(Service303Test, TestTimingMetrics) {
  Service303Client client = *Service303Test::client;

  MetricsContainer metrics_container;
  int status        = client.GetMetrics(&metrics_container);
  double start_time = Service303Test::findGauge(
      metrics_container,
      std::to_string(magma::orc8r::MetricName::process_start_time_seconds));
  double uptime = Service303Test::findGauge(
      metrics_container,
      std::to_string(magma::orc8r::MetricName::process_cpu_seconds_total));

  EXPECT_GT(start_time, 0);
  EXPECT_GE(uptime, 0);

  sleep(1);
  client.GetMetrics(&metrics_container);
  double new_start = Service303Test::findGauge(
      metrics_container,
      std::to_string(magma::orc8r::MetricName::process_start_time_seconds));
  double new_uptime = Service303Test::findGauge(
      metrics_container,
      std::to_string(magma::orc8r::MetricName::process_cpu_seconds_total));

  EXPECT_DOUBLE_EQ(start_time, new_start);
  EXPECT_GT(new_uptime, uptime);
}

TEST_F(Service303Test, TestMemoryMetrics) {
  Service303Client client = *Service303Test::client;

  MetricsContainer metrics_container;
  int status          = client.GetMetrics(&metrics_container);
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
TEST_F(Service303Test, TestEnumConversions) {
  Service303Client client = *Service303Test::client;

  // enum values (from metricsd.proto):
  // s1_setup => 500, result => 0
  increment_counter(
      "mme_new_association", 3, 2, "result", "success", "gateway", "1234");

  MetricsContainer metrics_container;
  int status = client.GetMetrics(&metrics_container);
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
  // Test service config
  Service303Test::config.name    = cstr2bstr("test_service");
  Service303Test::config.version = cstr2bstr("0.0.1");

  return RUN_ALL_TESTS();
}
