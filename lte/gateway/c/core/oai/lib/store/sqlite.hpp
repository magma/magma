#include <string>
#include <vector>
#include "lte/protos/subscriberdb.pb.h"

using google::protobuf::Message;

namespace magma {
namespace lte {
class SqliteStore {
 public:
  SqliteStore(std::string db_location, int sid_digits);

  // Initialize data store
  void init_db_connection(std::string db_location, int sid_digits = 2);

  // Add subscriber
  void add_subscriber(
      const SubscriberData& subscriber_data);  // TODO: add the parameters

  // Delete subscriber
  void delete_subscriber();  // TODO: add the parameters

 private:
  int _sid_digits;
  int _n_shards;
  std::vector<std::string> _db_locations;
  std::vector<std::string> _create_db_locations(std::string db_location,
                                                int _n_shards);
  void _create_store();
  std::string _to_str(const SubscriberData& subscriber_data);
  // Map subscriber ID to bucket
  int _sid2bucket(std::string sid);
};
}  // namespace lte
}  // namespace magma
