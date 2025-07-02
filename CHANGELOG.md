# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.0.2] - 2025-06-30

### Changed

- Replaced `string` constants with `EntityType` and `StreamResponseType` custom types for improved type safety.

### Fixed

- Fixed a potential deadlock in the `Chat` method's SSE handling by improving context cancellation management. This prevents goroutine leaks if the context is cancelled while streaming a response.

## [0.0.1] - 2025-06-20

### Added

- Initial release of the `gofabric` Go client library for the Fabric API.

[Unreleased]: https://github.com/sherif-fanous/gofabric/compare/v0.0.2...HEAD
[0.0.2]: https://github.com/sherif-fanous/gofabric/compare/v0.0.1...v0.0.2
[0.0.1]: https://github.com/sherif-fanous/gofabric/releases/tag/v0.0.1
