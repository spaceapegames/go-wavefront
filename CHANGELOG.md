# Changelog

Changelog for go-wavefront.

## [1.8.0]

*Add Chart Attributes*

- Added chartAttributes struct and the required nested structs for drilldown link support

## [1.7.0]

*Feature to configure http client timeout*

- Expose http client timeout parameter

## [1.4.0]

*Add Missing Fields to Dashboards*

- Add missing field from Sources (SecondaryAxis)

## [1.3.0]

*Add Missing Fields to Dashboards*

- A large number of fields previously missing from Dashboard have been implemented

## [1.2.0]

*Improvements to Dashboards*

- Add missing fields from Sources (ScatterPlotSource, Disabled, QuerybuilderEnabled and SourceDescription)
- Add Dynamic and List parameter types

## [1.1.12] - 2017-10-13

- Allow optional Alert fields to be omitted

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
