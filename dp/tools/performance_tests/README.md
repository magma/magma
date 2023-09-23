# Domain Proxy Database Performance Tests

Testing workflow for Domain Proxy database. The scenario focuses on loading data into domain_proxy_logs table
and running test queries against it simulating frontend generated traffic.

domain_proxy_logs table is the most likely one to cause performance drops due to its record volumes. In real life scenario
**1 log will be stored per radio every < 1 second on average.**

**IMPORTANT**:

Test containers will use local storage as volumes, it is **strongly advisable to prune all dangling local volumes**
after each test.

If you want to inspect all dangling volumes first, run: `docker volume ls -f dangling=true`

To delete dangling volumes, run: `docker volume prune`

This is especially important when running tests on large datasets as they can very quickly take up an entire disk
space when run multiple times.

## Preparing the environment

Below steps will build a **docker container** with a database whose schema will be based on available
migration files.

#### Database default parameters

```
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_PORT=5532
DB_NAME=dp
```

### Building the database and filling it up with data in one go

Go to `$MAGMA_ROOT/dp` and run a `make setup_performance_tests` command either on its own
_(default values will be applied)_ or with the following options:

`LOGS_COUNT` - how many log records should be created (default is 1000)

`CBSDS_COUNT` - how many CBSDs should these logs be associated with (default is 1000).

`LOG_FILE` - where should the test information on how many logs were generated be stored (default is <script_location>/performance_tests.log)

_Note that this does not imply a database relation as domain_proxy_logs table has no relations.
It only indicates a value range in the `cbsd_serial_number_column` (eg. CBSDS_COUNT=1000 means there will be a 1000
`cbsd_serial_numbers` distributed across all log records)._

example:

`make setup_performance_tests LOGS_COUNT=10000000 CBSDS_COUNT=100 LOG_FILE=/path/to/some/file.log`

**Note**: Creating very big datasets will take a significant amount of time due to the presence of indices
in the table. **Generating 100 mln records may take ~30 minutes** depending on the capabilities of the computer.

### Building the database separately

Go to `$MAGMA_ROOT/dp` and run `make build_db`

### Populating the database separately

Go to `$MAGMA_ROOT/dp` and run `make prepare_performance_test_data`
The same options as above may be used to parameterized data generation.

## Running tests

Go to `$MAGMA_ROOT/dp` and run a `make run_performance_tests` command either on its own
_(default values will be applied)_ or with the following options:

`LIMIT` - how many log records should be fetched (default is 100)

`OFFSET` - from which record onward should the data be returned (default is 0)

`LOG_FILE` - where should the test results be stored (default is <script_location>/performance_tests.log)

example:

`make run_performance_tests OFFSET=10 LOG_FILE=/path/to/some/file.log`

**Note**: Above options are meant to simulate pagination in the UI. Please use common sense in parameterizing them
especially with large datasets.

**Note**: The tests will be run against the database created in previous steps. Please **do not attempt to redirect
them to any other database**.

## Checking results

Test results will be stored in the file specified by the `LOG_FILE` option.
Example resultset might look something like this:

```angular2html
[2022-02-09 13:54:59,179] Running test queries
[2022-02-09 13:54:59,184] ---Query "SELECT domain_proxy_logs.log_from, domain_proxy_logs.log_to, domain_proxy_logs.created_date, domain_proxy_logs.fcc_id, domain_proxy_logs.response_code, domain_proxy_logs.cbsd_serial_number, domain_proxy_logs.log_name, domain_proxy_logs.log_message FROM domain_proxy_logs  WHERE (network_id = '1') ORDER BY created_date DESC LIMIT 10000;" 0.0023 seconds to execute ---
[2022-02-09 13:54:59,185] ---Query "SELECT domain_proxy_logs.log_to, domain_proxy_logs.fcc_id, domain_proxy_logs.log_name, domain_proxy_logs.log_message, domain_proxy_logs.created_date, domain_proxy_logs.log_from, domain_proxy_logs.cbsd_serial_number, domain_proxy_logs.response_code FROM domain_proxy_logs  WHERE (network_id = '1') ORDER BY created_date DESC LIMIT 10000;" 0.0022 seconds to execute ---
[2022-02-09 13:54:59,187] ---Query "SELECT domain_proxy_logs.log_name, domain_proxy_logs.cbsd_serial_number, domain_proxy_logs.log_from, domain_proxy_logs.log_message, domain_proxy_logs.fcc_id, domain_proxy_logs.log_to, domain_proxy_logs.response_code, domain_proxy_logs.created_date FROM domain_proxy_logs  WHERE (network_id = '1' AND log_from = 'CBSD') ORDER BY created_date DESC LIMIT 10000;" 0.0021 seconds to execute ---
[2022-02-09 13:54:59,188] ---Query "SELECT domain_proxy_logs.log_to, domain_proxy_logs.response_code, domain_proxy_logs.fcc_id, domain_proxy_logs.log_message, domain_proxy_logs.log_name, domain_proxy_logs.created_date, domain_proxy_logs.cbsd_serial_number, domain_proxy_logs.log_from FROM domain_proxy_logs  WHERE (network_id = '1' AND log_to = 'CBSD') ORDER BY created_date DESC LIMIT 10000;" 0.0019 seconds to execute ---
[2022-02-09 13:54:59,189] ---Query "SELECT domain_proxy_logs.log_from, domain_proxy_logs.log_to, domain_proxy_logs.created_date, domain_proxy_logs.cbsd_serial_number, domain_proxy_logs.log_name, domain_proxy_logs.log_message, domain_proxy_logs.response_code, domain_proxy_logs.fcc_id FROM domain_proxy_logs  WHERE (network_id = '1' AND fcc_id = 'some_fcc_id1') ORDER BY created_date DESC LIMIT 10000;" 0.0011 seconds to execute ---
[2022-02-09 13:54:59,191] ---Query "SELECT domain_proxy_logs.log_message, domain_proxy_logs.response_code, domain_proxy_logs.created_date, domain_proxy_logs.log_to, domain_proxy_logs.cbsd_serial_number, domain_proxy_logs.fcc_id, domain_proxy_logs.log_from, domain_proxy_logs.log_name FROM domain_proxy_logs  WHERE (network_id = '1' AND cbsd_serial_number = 'some_cbsd_id1') ORDER BY created_date DESC LIMIT 10000;" 0.0016 seconds to execute ---
[2022-02-09 13:54:59,194] ---Query "SELECT domain_proxy_logs.log_name, domain_proxy_logs.cbsd_serial_number, domain_proxy_logs.response_code, domain_proxy_logs.log_to, domain_proxy_logs.created_date, domain_proxy_logs.log_from, domain_proxy_logs.log_message, domain_proxy_logs.fcc_id FROM domain_proxy_logs  WHERE (network_id = '1' AND log_name = 'some_name') ORDER BY created_date DESC LIMIT 10000;" 0.0022 seconds to execute ---
[2022-02-09 13:54:59,196] ---Query "SELECT domain_proxy_logs.fcc_id, domain_proxy_logs.created_date, domain_proxy_logs.cbsd_serial_number, domain_proxy_logs.log_name, domain_proxy_logs.response_code, domain_proxy_logs.log_to, domain_proxy_logs.log_from, domain_proxy_logs.log_message FROM domain_proxy_logs  WHERE (network_id = '1' AND response_code = 0) ORDER BY created_date DESC LIMIT 10000;" 0.0023 seconds to execute ---
[2022-02-09 13:54:59,199] ---Query "SELECT domain_proxy_logs.log_to, domain_proxy_logs.created_date, domain_proxy_logs.fcc_id, domain_proxy_logs.log_from, domain_proxy_logs.log_name, domain_proxy_logs.cbsd_serial_number, domain_proxy_logs.log_message, domain_proxy_logs.response_code FROM domain_proxy_logs  WHERE (network_id = '1' AND created_date <= '2022-02-04 00:14:02+00') ORDER BY created_date DESC LIMIT 10000;" 0.0022 seconds to execute ---
[2022-02-09 13:54:59,200] ---Query "SELECT domain_proxy_logs.fcc_id, domain_proxy_logs.log_from, domain_proxy_logs.log_message, domain_proxy_logs.created_date, domain_proxy_logs.response_code, domain_proxy_logs.log_to, domain_proxy_logs.log_name, domain_proxy_logs.cbsd_serial_number FROM domain_proxy_logs  WHERE (network_id = '1' AND fcc_id = 'some_fcc_id1' AND cbsd_serial_number = 'some_cbsd_id1' AND log_name = 'some_name' AND created_date >= '2022-02-04 00:14:02+00' AND created_date <= '2022-02-04 00:19:02+00') ORDER BY created_date DESC LIMIT 10000;" 0.0011 seconds to execute ---
[2022-02-09 13:54:59,201] ---Query "SELECT domain_proxy_logs.log_message, domain_proxy_logs.response_code, domain_proxy_logs.created_date, domain_proxy_logs.log_from, domain_proxy_logs.log_to, domain_proxy_logs.cbsd_serial_number, domain_proxy_logs.log_name, domain_proxy_logs.fcc_id FROM domain_proxy_logs  WHERE (network_id = '1' AND created_date >= '2022-02-04 00:14:02+00' AND created_date <= '2022-02-04 00:19:02+00') ORDER BY created_date DESC LIMIT 10000;" 0.001 seconds to execute ---
```
