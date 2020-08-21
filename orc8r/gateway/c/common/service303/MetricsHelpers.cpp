#include "MetricsHelpers.h"

#include "MetricsSingleton.h"

namespace magma {
namespace service303 {

void increment_counter(
    const char* name, double increment, size_t n_labels, ...) {
  va_list ap;
  va_start(ap, n_labels);
  MetricsSingleton::Instance().IncrementCounter(name, increment, n_labels, ap);
  va_end(ap);
}

void increment_gauge(const char* name, double increment, size_t n_labels, ...) {
  va_list ap;
  va_start(ap, n_labels);
  MetricsSingleton::Instance().IncrementGauge(name, increment, n_labels, ap);
  va_end(ap);
}

void decrement_gauge(const char* name, double decrement, size_t n_labels, ...) {
  va_list ap;
  va_start(ap, n_labels);
  MetricsSingleton::Instance().DecrementGauge(name, decrement, n_labels, ap);
  va_end(ap);
}

void set_gauge(const char* name, double value, size_t n_labels, ...) {
  va_list ap;
  va_start(ap, n_labels);
  MetricsSingleton::Instance().SetGauge(name, value, n_labels, ap);
  va_end(ap);
}

void observe_histogram(
    const char* name, double observation, size_t n_labels, ...) {
  va_list ap;
  va_start(ap, n_labels);
  MetricsSingleton::Instance().ObserveHistogram(
      name, observation, n_labels, ap);
  va_end(ap);
}

}
}