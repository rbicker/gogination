# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.0.3] - 2021-01-06
### Fixed
- fixed non-working filters due to incorrect nesting of bson.D objects 

## [0.0.2] - 2021-01-04
### Fixed
- renamed error ErrExceptedStruct to ErrExpectedStruct
- removed $oid filter

## [0.0.1] - 2021-01-04
### Added
- added NextFilter function to determine a filter for next document based on given document, filter and orderBy