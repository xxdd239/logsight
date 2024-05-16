
truncate table ts_radius_accounting;
SELECT create_hypertable('ts_radius_accounting', 'acct_start_time');

ALTER TABLE ts_radius_accounting
    SET (
        timescaledb.compress,
        timescaledb.compress_orderby = 'acct_start_time desc',
        timescaledb.compress_segmentby = 'id, username, acct_session_id'
        );

SELECT set_chunk_time_interval('ts_radius_accounting', INTERVAL '24 hours');
-- SELECT remove_compression_policy('ts_radius_accounting');
SELECT add_compression_policy('ts_radius_accounting', INTERVAL '120d');
SELECT add_retention_policy('ts_radius_accounting', INTERVAL '1 year');



truncate table ts_syslog;
SELECT create_hypertable('ts_syslog', 'timestamp');

ALTER TABLE ts_syslog
    SET (
        timescaledb.compress,
        timescaledb.compress_orderby = 'timestamp desc',
        timescaledb.compress_segmentby = 'id'
        );

SELECT set_chunk_time_interval('ts_syslog', INTERVAL '24 hours');
-- SELECT remove_compression_policy('ts_syslog');
SELECT add_compression_policy('ts_syslog', INTERVAL '15d');
SELECT add_retention_policy('ts_syslog', INTERVAL '6 months');


