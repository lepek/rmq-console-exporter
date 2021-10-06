# RMQ-CONSOLE-EXPORTER
This is an agent that exposes a `/metrics` endpoint which return RabbitMQ Server related metrics in 
Prometheus string format.

## Description
When the `/metrics` endpoint is accessed a metric collection process starts: 
- One or more binaries are executed
- The stdout of those commands is parsed and transformed into metrics
- The metrics are returned and exposed in Prometheus string format

## Usage
```bash
./rmq-console-exporter --help                                                                                                                                          ✔
Usage of ./rmq-console-exporter:
  -log_level string
    	Log Level: debug, info, error, etc (default "info")
  -output_buffer int
    	Output Buffer[lines] (default 100000)
  -port int
    	Port to expose metrics (default 2112)
  -prefix string
    	Metrics prefix (default "rmq_")
  -queue_parser string
    	Queue Parser to use: json or tabular (default "json")
  -timeout int
    	Timeout[Ms] for each collector (default 30000)
```

## Sample Output
```http request
bash-4.2$ curl -s http://127.0.0.1:2112/metrics
# HELP rmq_consumer_utilisation Fraction of the time (between 0.0 and 1.0) that the queue is able to immediately deliver messages to consumers. This can be less than 1.0 if consumers are limited by network congestion or prefetch count.
# TYPE rmq_consumer_utilisation gauge
rmq_consumer_utilisation{queue="10_128_4_241:5672.sage-xds-service.35k_0204_kc.LATEST.perf",state="running"} 1
# HELP rmq_consumers Number of consumers.
# TYPE rmq_consumers gauge
rmq_consumers{queue="10_128_4_241:5672.sage-xds-service.35k_0204_kc.LATEST.perf",state="running"} 1
# HELP rmq_memory Bytes of memory allocated by the runtime for the queue, including stack, heap and internal structures.
# TYPE rmq_memory gauge
rmq_memory{queue="10_128_4_241:5672.sage-xds-service.35k_0204_kc.LATEST.perf",state="running"} 55756
# HELP rmq_message_bytes_ready Like message_bytes but counting only those messages ready to be delivered to clients.
# TYPE rmq_message_bytes_ready gauge
rmq_message_bytes_ready{queue="10_128_4_241:5672.sage-xds-service.35k_0204_kc.LATEST.perf",state="running"} 0
# HELP rmq_message_bytes_unacknowledged Like message_bytes but counting only those messages delivered to clients but not yet acknowledged.
# TYPE rmq_message_bytes_unacknowledged gauge
rmq_message_bytes_unacknowledged{queue="10_128_4_241:5672.sage-xds-service.35k_0204_kc.LATEST.perf",state="running"} 0
# HELP rmq_messages_ready Number of messages ready to be delivered to clients.
# TYPE rmq_messages_ready gauge
rmq_messages_ready{queue="10_128_4_241:5672.sage-xds-service.35k_0204_kc.LATEST.perf",state="running"} 0
# HELP rmq_messages_unacknowledged Like message_bytes but counting only those messages ready to be delivered to clients.
# TYPE rmq_messages_unacknowledged gauge
rmq_messages_unacknowledged{queue="10_128_4_241:5672.sage-xds-service.35k_0204_kc.LATEST.perf",state="running"} 0
# HELP rmq_command_runtime_seconds Runtime of the command executed to collect the metrics.
# TYPE rmq_command_runtime_seconds gauge
rmq_command_runtime_seconds{command_executed="rabbitmqctl list_queues --formatter json name state messages_ready message_bytes_ready messages_unacknowledged message_bytes_unacknowledged memory consumers consumer_utilisation head_message_timestamp"} 0.5431952
```

## Exposed Metrics

### Queue Metrics

#### Metrics
- `messages_ready`: Number of messages ready to be delivered to clients.
- `message_bytes_ready`: Like message_bytes but counting only those messages ready to be delivered to clients.
- `messages_unacknowledged`: Like message_bytes but counting only those messages ready to be delivered to clients.
- `message_bytes_unacknowledged`: Like message_bytes but counting only those messages delivered to clients but not yet 
  acknowledged.
- `memory`: Bytes of memory allocated by the runtime for the queue, including stack, heap and internal structures.
- `consumers`: Number of consumers.
- `consumer_utilisation`: Fraction of the time (between 0.0 and 1.0) that the queue is able to immediately deliver 
  messages to consumers. This can be less than 1.0 if consumers are limited by network congestion or prefetch count.
- `head_message_timestamp`: The timestamp property of the first message in the queue, if present. 
  Timestamps of messages only appear when they are in the paged-in state.

#### Labels
- `queue`: The name of the queue with non-ASCII characters escaped as in C.
- `state`: The state of the queue. Normally "running", but may be "{syncing, message_count}" if the queue is synchronising.

### Agent metrics

#### Metrics
- `command_runtime`: Runtime of the command executed to collect the metrics.

#### Labels
- `command_executed`: Full command executed with arguments.

## Changelog
### 0.1
- OMSDPM-5724: First release. It provides RabbitMQ queue metrics from the `rabbitmqctl` command output. 
  It can get the metrics from tabular output and from json output (using `--formatter json`).
  It also provides an additional metric regarding the collection itself.
  Since the collection can take a long time, to avoid hammering the RMQ server, 
  the agent prevents concurrent runs of the same collection.
  It has been tested on RMQ Server 3.6, 3.7, 3.8 and 3.9.  


