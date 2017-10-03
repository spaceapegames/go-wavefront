# Changelog

Changelog for go-wavefront.

## [1.1.12] - 2017-10-03

For Events process annotations before unmarshelling JSON

## [1.1.11] - 2017-09-14

*Add the ability to manage alert targets*

- Support for Alert Targets (notificants)

## [1.1.0] - 2017-08-17

*Add the ability to manage dashboards*

- Support for dashboards

## [1.0.0] - 2017-07-17

*Complete re-write of libraries. Breaking API changes*

- Re-write of library code to make compatible with the Wavefront v2 API.
- Support for Alerts, Querying, Search, Events.
- Writer now supports metric tagging.
- Remove CLI, restructure code, sanitise data-structures, make more idiomatic.
