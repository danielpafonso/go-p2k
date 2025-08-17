# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0]

### Add

- Config endpoints
- TLS option for kafka
- Script to generate certificates
- Environment variable `KAFKA_CHECK_CERT`
- Type check when parsing \_kafka field

### Changed

- Podman Kafka version to 4.0.0

## [0.1.0]

### Add

- Health Check and Metrics endpoints
- Initial Metrics
- Multiple Kafka Cluster publisher, with Deadletter to console
- Environment Variable to overwrite Kafka configuration, `KAFKA_CLUSTERS`
- Pub/Sub message parsing from json

### Changed

- Use log instead of fmt for runtime log messages

### Removed

- Environment variables `KAFKA_ENDPOINTS` and `KAFKA_TOPIC`

## [0.0.1]

### Add

- Load configurations and overwrite with Environment variables
- Pub/Sub subscription client
- Kafka Producer
- Main loop logic

---

[unreleased]: https://github.com/danielpafonso/go-p2k/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/danielpafonso/go-p2k/releases/tag/v0.2.0
[0.1.0]: https://github.com/danielpafonso/go-p2k/releases/tag/v0.1.0
[0.0.1]: https://github.com/danielpafonso/go-p2k/releases/tag/v0.0.1
