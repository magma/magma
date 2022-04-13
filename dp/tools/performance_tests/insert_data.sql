do $$
declare
   max_cbsd_id integer := {cbsds_count};
   how_many integer := {logs_count};
   log_from_array TEXT ARRAY DEFAULT ARRAY['CBSD','ACS','DP'];
   log_name_array TEXT ARRAY DEFAULT ARRAY['some_name','some_name1','some_name2','some_name3','some_name4','some_name5','some_other_name','yet_another_name','yet_another_name','different_name'];
   response_code_array INT ARRAY DEFAULT ARRAY[0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15];
   log_from_array_len integer := 3;
   log_name_array_len integer := 10;
   response_code_array_len integer := 16;
   min_date date := '2022-02-04 00:00:00+00';
   max_date date := '2022-02-07 00:00:00+00';
begin
--       Dropping indexes to speed up the inserts, adding them back at the end
      DROP INDEX ix_domain_proxy_logs_cbsd_serial_number;
      DROP INDEX ix_domain_proxy_logs_created_date;
      DROP INDEX ix_domain_proxy_logs_fcc_id;
      DROP INDEX ix_domain_proxy_logs_response_code;
      DROP INDEX ix_domain_proxy_logs_log_name;
      INSERT INTO domain_proxy_logs (
          log_from,
          log_to,
          log_name,
          log_message,
          cbsd_serial_number,
          network_id,
          fcc_id,
          response_code,
          created_date
      )
      SELECT
          log_from_array[floor(random() * log_from_array_len + 1)],
          log_from_array[floor(random() * log_from_array_len + 1)],
          log_name_array[floor(random() * log_name_array_len + 1)],
          'foo',
          'some_cbsd_id' || floor(random() * max_cbsd_id + 1),
          'some_network_id',
          'some_fcc_id' || floor(random() * max_cbsd_id + 1),
          response_code_array[floor(random() * response_code_array_len+1)],
          min_date::timestamp + random() * (max_date::timestamp - min_date::timestamp)
      FROM generate_series(1, how_many);
      CREATE INDEX ix_domain_proxy_logs_cbsd_serial_number ON domain_proxy_logs (cbsd_serial_number);
      CREATE INDEX ix_domain_proxy_logs_created_date ON domain_proxy_logs (created_date);
      CREATE INDEX ix_domain_proxy_logs_fcc_id ON domain_proxy_logs (fcc_id);
      CREATE INDEX ix_domain_proxy_logs_response_code ON domain_proxy_logs (response_code);
      CREATE INDEX ix_domain_proxy_logs_log_name ON domain_proxy_logs (log_name);
end;
$$