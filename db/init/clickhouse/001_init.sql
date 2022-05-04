CREATE TABLE log_queue (
    timestamp DateTime,
    message_type String,
    message String
) ENGINE = Kafka 
SETTINGS kafka_broker_list = 'kafka:9092',
    kafka_topic_list = 'log',
    kafka_group_name = 'clickhouse1',
    kafka_format = 'JSONEachRow';

CREATE TABLE stripe_log_table
(
    timestamp DateTime,
    message_type String,
    message String
)
ENGINE = StripeLog;

CREATE MATERIALIZED VIEW log_consumer TO stripe_log_table
    AS SELECT *
    FROM log_queue;