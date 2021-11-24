#include <iostream>
#include <fstream>
#include <string>
#include <vector>

#include "spgw_state_converter.h"

namespace magma {
namespace lte {

#define DEFAULT_SPGW_CONTEXT_DATA_PATH                                         \
  "/home/vagrant/magma/lte/gateway/c/core/oai/test/spgw_task/data/"

std::vector<std::string> load_file_into_vector_of_line_content(
    const std::string& data_folder_path, const std::string& file_name);
status_code_e mock_read_spgw_ue_state_db(
    const std::vector<std::string>& ue_samples);
}  // namespace lte
}  // namespace magma