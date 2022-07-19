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
#include "lte/gateway/c/core/oai/lib/store/sqlite.hpp"
#include "lte/protos/subscriberdb.pb.h"

#include <cmath>
#include <vector>
#include <sqlite3.h>

using google::protobuf::Message;
namespace magma {
namespace lte {
SqliteStore::SqliteStore(std::string db_location, int sid_digits) {
  init_db_connection(db_location, sid_digits);
}
void SqliteStore::init_db_connection(std::string db_location, int sid_digits) {
  _sid_digits = sid_digits;
  _n_shards = std::pow(10, sid_digits);
  _db_locations = _create_db_locations(db_location, _n_shards);
  _create_store();
}

std::vector<std::string> SqliteStore::_create_db_locations(
    std::string db_location, int n_shards) {
  // in memory if db_location is not specified
  if (db_location.length() == 0) {
    db_location = "/var/opt/magma/";
  }

  std::vector<std::string> db_location_list;
  for (int shard = 0; shard < n_shards; shard++) {
    std::string to_push = "file:" + db_location + "subscriber" +
                          std::to_string(shard) + ".db?cache=shared";
    db_location_list.push_back(to_push);
    std::cout << "[LOG] DB location: " << db_location_list[shard] << std::endl;
  }
  return db_location_list;
}

void SqliteStore::_create_store() {
  int rc;
  for (std::string db_location_s : _db_locations) {
    sqlite3* db;
    int rc;
    const char* db_location = db_location_s.c_str();
    rc = sqlite3_open(db_location, &db);
    if (rc) {
      std::cout << "Cannot open database " << sqlite3_errmsg(db) << std::endl;
    } else {
      std::cout << "Database opened successfully" << std::endl;
    }

    const char* sql =
        "CREATE TABLE IF NOT EXISTS subscriberdb"
        "(subscriber_id text PRIMARY KEY, data text)";
    char* zErrMsg;
    rc = sqlite3_exec(db, sql, NULL, 0, &zErrMsg);

    if (rc != SQLITE_OK) {
      std::cout << "SQL Error " << zErrMsg << std::endl;
      sqlite3_free(zErrMsg);
    } else {
      std::cout << "Table created successfully!!" << std::endl;
    }

    sqlite3_close(db);
  }
}

void SqliteStore::add_subscriber(SubscriberData& subscriber_data) {
  std::string sid_s = _to_str(subscriber_data);
  const char* sid = sid_s.c_str();
  std::string data_str;
  subscriber_data.SerializeToString(&data_str);  // TODO: serialize to string
  std::string db_location_s = _db_locations[_sid2bucket(sid)];
  const char* db_location = db_location_s.c_str();
  sqlite3* db;
  int rc = sqlite3_open(db_location, &db);
  if (rc) {
    std::cout << "Cannot open database " << sqlite3_errmsg(db) << std::endl;
  } else {
    std::cout << "Database opened successfully" << std::endl;
  }
  const char* sql = "SELECT data FROM subscriberdb WHERE subscriber_id = ?";
  sqlite3_stmt* stmt;
  const char* pzTail;
  int rc2 = sqlite3_prepare_v2(db, sql, strlen(sql), &stmt, &pzTail);
  if (rc2 == SQLITE_OK) {
    sqlite3_bind_text(
        stmt, 1, sid, 4 * strlen(sid),
        SQLITE_STATIC);  // REVIEW THAT THE PARAMETERS HERE ARE CORRECT
    std::cout << "Successful data binding" << std::endl;
    sqlite3_step(stmt);
    sqlite3_finalize(stmt);
  } else {
    std::cout << "SQL Error " << std::endl;
  }
}

// function is hardcoded for now, will fix
std::string SqliteStore::_to_str(const SubscriberData& subscriber_data) {
  if (subscriber_data.sid().type() == SubscriberID::IMSI) {
    return "IMSI" + subscriber_data.sid().id();
  } else {
    std::cout << "Invalid sid " << subscriber_data.sid().id() << " type "
              << subscriber_data.sid().type() << std::endl;
  }
}

int SqliteStore::_sid2bucket(std::string sid) {
  int bucket;
  try {
    bucket = std::stoi(sid.substr(sid.length() - _sid_digits, sid.length()));
  } catch (int bucket) {
    std::cout << "Last " << _sid_digits << "digits of subscriber id " << sid
              << " cannot be mapped to a bucket, default to bucket 0"
              << std::endl;
    bucket = 0;
  }
  return bucket;
}
}  // namespace lte
}  // namespace magma
