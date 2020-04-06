# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## Changed

* Make variable creation process asynchronous. [#31]

## [1.2.0] - 2020-04-01

### Added

* Support multiple varfiles when create/updating variables. [#28]

### Fixed

* Incorrect encoding of HCL lists. [#27]

## [1.1.1] - 2020-03-31

### Fixed

* Fix the problem preventing to update an existing workspace, or to create a new
  workspace with the`--force` flag. [#23]

## [1.1.0] - 2020-03-30

### Added

* Ability to detect HCL variables from a varfile. [#14]
* Ability to delete workspaces. [#15]
* Ability to delete variables. [#20]

### Fixed

* Fix problem preventing to create a new workspace. [#18]
* Fix the installer. [#19]

## [1.0.0] - 2020-03-25

Initial version with support for managing:

* Workspaces
* Variables

[//]: # (Release links)
[1.0.0]: https://github.com/rgreinho/tfe-cli/releases/tag/1.0.0
[1.1.0]: https://github.com/rgreinho/tfe-cli/releases/tag/1.1.0
[1.1.1]: https://github.com/rgreinho/tfe-cli/releases/tag/1.1.1
[1.2.0]: https://github.com/rgreinho/tfe-cli/releases/tag/1.2.0

[//]: # (Issue/PR links)
[#14]: https://github.com/rgreinho/tfe-cli/pull/14
[#15]: https://github.com/rgreinho/tfe-cli/pull/15
[#18]: https://github.com/rgreinho/tfe-cli/pull/18
[#19]: https://github.com/rgreinho/tfe-cli/pull/19
[#20]: https://github.com/rgreinho/tfe-cli/pull/20
[#23]: https://github.com/rgreinho/tfe-cli/pull/23
[#27]: https://github.com/rgreinho/tfe-cli/pull/27
[#28]: https://github.com/rgreinho/tfe-cli/pull/28
[#31]: https://github.com/rgreinho/tfe-cli/pull/31
