# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [0.91.1] - Unreleased
### Fixed
- Allow the creation of realms when explicitly setting a replication class. Fixes a regression
  introduced in v0.91.0.

## [0.91.0] - 2023-05-29
### Changed
- Replace device `metadata` with `attributes`.
- BREAKING: Remove the 0.90.1 Astarte client API and introduce a clean, idiomatic API.
  See [#33](https://github.com/astarte-platform/astarte-go/issues/33).
- BREAKING: add the `isAsync` parameter to InstallInterface and UpdateInterface functions to allow synchronous calls.

### Fixed
- `Raw` properly updates the paginator's state when retrieving paginated data.
- Fix the logic for retrieving data from Appengine API for both time series and data snapshots.

## [0.90.1] - 2021-03-03
### Changed
- Update dependencies

## [0.90.0] - 2021-03-02
### Added
- Initial release
