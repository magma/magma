#include "mock_spgw_op.h"
#include "spgw_state.h"

#include <google/protobuf/text_format.h>


namespace magma {
namespace lte {

int mock_read_spgw_ue_state_db() {
  int number_of_ue_samples = 1;
  for (int i = 0; i < number_of_ue_samples; ++i) {
    oai::SpgwUeContext ue_proto = oai::SpgwUeContext();
    std::string name_of_sample = "/home/vagrant/magma/lte/gateway/c/core/oai/test/spgw_task/data/sample.bin.v10";
    std::fstream input(name_of_sample.c_str(), std::ios::in | std::ios::binary);
    if (!ue_proto.ParseFromIstream(&input)) {
      std::cerr << "Failed to parse the sample: " << name_of_sample << std::endl;
      return -1;
    }
    
  
  spgw_ue_context_t* ue_context_p =
      (spgw_ue_context_t*) calloc(1, sizeof(spgw_ue_context_t));
  SpgwStateConverter::proto_to_ue(ue_proto, ue_context_p);
  }
  return 0;
}

int mock_spgw_app_init(){
    mock_read_spgw_ue_state_db();
    return 0;
}

}
}