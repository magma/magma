#include <fstream>
#include <algorithm>

#include "mock_spgw_op.h"
#include "spgw_state.h"

#include <google/protobuf/text_format.h>

namespace magma {
namespace lte {

std::vector<std::string> load_file_into_vector_of_line_content(
    const std::string& file_name) {
  std::fstream file_content(file_name.c_str(), std::ios_base::in);
  std::string line;
  std::vector<std::string> vector_of_lines;
  if (file_content) {
    while (std::getline(file_content, line)) {
      line.erase(std::remove(line.begin(), line.end(), '\n'), line.end());
      vector_of_lines.push_back(line);
    }
  } else {
    std::cerr << "couldn't open file: " << file_name << std::endl;
  }

  file_content.close();
  return vector_of_lines;
}

int mock_read_spgw_ue_state_db(const std::vector<std::string>& ue_samples) {
  for (auto name_of_sample : ue_samples) {
    oai::SpgwUeContext ue_proto = oai::SpgwUeContext();
    std::fstream input(name_of_sample.c_str(), std::ios::in | std::ios::binary);
    if (!ue_proto.ParseFromIstream(&input)) {
      std::cerr << "Failed to parse the sample: " << name_of_sample
                << std::endl;
      return -1;
    }

    spgw_ue_context_t* ue_context_p =
        (spgw_ue_context_t*) calloc(1, sizeof(spgw_ue_context_t));
    SpgwStateConverter::proto_to_ue(ue_proto, ue_context_p);
  }
  return 0;
}

}  // namespace lte
}  // namespace magma