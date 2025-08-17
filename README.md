# Pubsub to Kafka

P2K is a shipper which reads messages from a GCP Pub/Sub subscription and write them into a Kafka

# PubSub Contract

go-P2K expect that the publish events are json structured.

## Required fields

The publish event requires a field named `_kafka` with the bellow structure. This fields indicates to which configured cluster and topic the evet will be sent and after processing it it will be removed from the event.

| Field            | Type                   | Required? | Description                                                                                                   |
| ---------------- | ---------------------- | --------- | ------------------------------------------------------------------------------------------------------------- |
| \_kafka.topic    | string                 | yes       | Kafka's topic where to send the event                                                                         |
| \_kafka.clusters | comma seperated string | no        | Cluster names where events should be sent. If not specified the event will be sent to all configured clusters |

Example:

```
{
  "_kafka": {
    "topic": "test-topic",
    "clusters": "docker-kafka1, kafka2"
  }
}
```

# Configurations

## Configurations File

| Field                    | Type             | Description                                                            |
| ------------------------ | ---------------- | ---------------------------------------------------------------------- |
| pubsub.project           | string           | GCP Project ID where the Pub/Sub is                                    |
| pubsub.subscription      | string           | Pub/Sub's Subscripntion name to connect                                |
| kafka.clusters           | array of object  | List of Kafka configurations, listed bellow                            |
| kafka.clusters.name      | string           | Name of Kafka cluster, used by messages                                |
| kafka.clusters.endpoints | array of strings | Comman seperated vist of Kafka brokers                                 |
| kafka.useTls             | bool             | Toogle the use of TLS in connection to Kafka                           |
| kafka.caFile             | string           | Path to CA file                                                        |
| kafka.crtFile            | string           | Path to TLS certifica file                                             |
| kafka.keyFile            | string           | Path to TLS key file                                                   |
| kafka.checkCrt           | bool             | Enables Certificate verification. Set to false only for debug proposes |

Example:

```json
{
  "pubsub": {
    "project": "acme",
    "subscription": "sub-events"
  },
  "kafka": {
    "clusters": [
      {
        "name": "kafka",
        "endpoints": "localhost:9094"
      }
    ],
    "useTls": false,
    "caFile": "ca file",
    "crtFile": "crt file",
    "keyFile": "key file",
    "checkCrt": true
  }
}
```

## Environment Variables

| Environment Variable | Type/Value                     |
| -------------------- | ------------------------------ |
| PUBSUB_PROJECT       | string                         |
| PUBSUB_SUBSCRIPTION  | string                         |
| KAFKA_CLUSTERS       | Kafkas configuration in json   |
| KAFKA_USE_TLS        | bool (yes/no, true/false, 1/0) |
| KAFKA_CA_FILE        | string                         |
| KAFKA_CRT_FILE       | string                         |
| KAFKA_KEY_FILE       | string                         |
| KAFKA_CHECK_CERT     | bool (yes/no, true/false, 1/0) |
