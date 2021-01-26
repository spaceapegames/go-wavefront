# Changelog

Changelog for go-wavefront-management-api.

## [1.14.0]

* Add support for ingestion policies

## [1.12.0]

*Add IsLogIntegration field to ExternalLink type*

Bug fixes:
- Fix json encoding / decoding of ExternalLink instances.

Enhancements:
- CRUD operations changed to be more efficient and to modify structs in ways
that are compatible with the assignment operator.

## [1.11.0]

*Add Support for Maintenance Windows*

Bug fixes:
- ScatterPlotSource field can be omitted when creating a dashboard
- fix all lint errors
- Make example code show up in godoc tools
- Account for ID field of user struct being nil
- NewClient to make defensive copy of config instance
- Users rewritten for performance improvement

## [1.10.0]

*Add support for Service Accounts*

## [1.9.0]

*Adds support for roles*

## [1.8.1]
*Fixes dashboard taggging after creation by adding SetTags functionality*

## [1.8.0]

*Add Support for CloudIntegrations*
 - CloudWatch
 - CloudTrail
 - EC2
 - GCP
 - GCPBilling
 - NewRelic
 - AppDynamics
 - Tesla
 - Azure
 - Azure Activity Log

*Add Support for Advanced Alert settings*

- CheckingFrequencyInMinutes
- EvaluateRealtimeData
- IncludeObsoleteMetrics

*Common client operations refactored into centralized location*
*Added support for skipTrash*

- Alerts
- Dashboards
- DerivedMetrics
- CloudIntegrations

*Fixed some failing tests*
*Add Annotations field to Events*

## [1.7.3]

- Fixing go.mod file

## [1.7.2]

- Add support for AccessControlList management on Alerts
- Add support for AccessControlList management on Dashboards

## [1.7.1]

- fixed missing UserGroup id on update call

## [1.7.0]

*Add Support for more Wavefront Primitives*

- Add support for Derived Metrics
- Add support for Users
- Add support for User Groups

*Improvements to Targets*

- Support for Routes on Targets

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
