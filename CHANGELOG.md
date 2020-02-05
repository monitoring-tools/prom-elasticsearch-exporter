# Changelog
All notable changes to this project will be documented in this file.

## [1.2.2] - 2020-01-05
### Changed
- Metric "task" was renamed to "task_group_duration_seconds".
### Added
- Label "cluster" to "task_group_duration_seconds" metric.

## [1.2.1] - 2019-12-24
### Added
- Nodes and indices requests cache hit and miss count 

## [1.2.0] - 2019-11-26
### Added
- Tasks running time metrics added 

## [1.1.0] - 2017-11-16
### Added
- CHANGELOG.md file was added

### Changed
- Multiple ElasticSearch nodes support removed.
  Now one exporter can fetch data from one ElasticSearch node.
