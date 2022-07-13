#include "lte/protos/subscriberdb.pb.h"
#include "lte/gateway/c/core/oai/lib/store/sqlite.hpp"
#include <cmath>
#include <vector>
#include <sqlite3.h>

using google::protobuf::Message;

SqliteStore::SqliteStore() { init_db_connection(); }

void SqliteStore::init_db_connection(std::string db_location,
                                     int sid_digits = 2) {
  _sid_digits = sid_digits;
  _n_shards = std::pow(10, sid_digits);
  _db_locations = _create_db_locations(db_location, _n_shards);
  _create_store();
}

std::vector<std::string> SqliteStore::_create_db_locations(
    std::string db_location, int n_shards) {
  // in memory if db_location is not specified
  if (db_location.length() == 0) {
    db_location = "/var/opt/magma/"
  }

  vector<std::string> db_location_list;
  for (int shard = 0; shard < n_shards; shard++) {
    std::string to_push = 'file:' + db_location + 'subscriber' +
                          std::to_string(shard) + ".db?cache=shared";
    db_location_list.push_back(string to_push);
    std::cout << "[LOG] DB location: " << db_location_list[shard] << endl;
  }
  return db_location_list;
}

void SqliteStore::_create_store() {
  int rc;
  for (int i = 0; i < _db_locations.length(); i++) {
    sqlite* db;
    int rc;
    rc = sqlite3_open(_db_locations[i], &db);
    if (rc) {
      std::cout << "Cannot open database " << sqlite3_errmsg(db) << endl;
    } else {
      std::cout << "Database opened successfully" << endl;
    }

    std::string sql =
        "CREATE TABLE IF NOT EXISTS subscriberdb"
        "(subscriber_id text PRIMARY KEY, data text)";
    rc = sqlite3_exec(db, sql, callback, 0,
                      zErrMsg);  // TODO: define callback function, figure out
                                 // what zErrMsg should look like.

    if (rc != SQLITE_OK) {
      std::cout << "SQL Error " << zErrMsg << endl;
      sqlite3_free(zErrMsg);
    } else {
      std::cout << "Table created successfully!!" << endl;
    }

    sqlite_close(db);
  }
}

void SqliteStore::add_subscriber(const SubscriberData& subscriber_data) {
  sid = to_str(subscriber_data.sid);
  data_str = subscriber_data.SerializeToString();
  db_location = _db_locations[_sid2bucket(sid)];
  int rc = sqlite3_open(db_location);
  if (rc) {
    std::cout << "Cannot open database " << sqlite3_errmsg(db) << endl;
  } else {
    std::cout << "Database opened successfully" << endl;
  }
  std::string sql = "SELECT data FROM subscriberdb WHERE subscriber_id = ?";
  sqlite3_stmt* stmt;
  std::string pzTail;
  int rc2 = sqlite3_prepare_v2(db, sql, sql.length(), &stmt, &pzTail);
  if (rc2 == SQLITE_OK) {
    sqlite_bind_text(stmt, 1, sid);
    std::cout << "Successful data binding" << endl;
    sqlite3_step(stmt);
    sqlite3_finalize(stmt);
  }
  if (rc != SQLITE_OK) {
    std::cout << "SQL Error: " << zErrMsg << endl;
  } else {
    std::cout << "Add subscriber successful" << endl;
  }
}

std::string SqliteStore::to_str(sid_pb) {
  if (sid_pb.type == SubscriberID.IMSI) {
    return "IMSI" + sid_pb.id
  }
}

int SqliteStore::_sid2bucket(std::string sid) {
  int bucket;
  try {
    bucket = std::stoi(sid.substr(sid.length() - _sid_digits, sid.length()));
  } catch (int bucket) {
    std::cout << "Last " << _sid_digits << "digits of subscriber id " << sid
              << " cannot be mapped to a bucket, default to bucket 0" << endl;
    bucket = 0
  }
  return bucket;
}
