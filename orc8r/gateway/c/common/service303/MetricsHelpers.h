#pragma once

#include <stdio.h>
#include <stdarg.h>

namespace magma {
namespace service303 {

void increment_counter(
    const char* name, double increment, size_t n_labels, ...);

void increment_gauge(const char* name, double increment, size_t n_labels, ...);

void decrement_gauge(const char* name, double decrement, size_t n_labels, ...);

void set_gauge(const char* name, double value, size_t n_labels, ...);

void observe_histogram(
    const char* name, double observation, size_t n_labels, ...);


}
}