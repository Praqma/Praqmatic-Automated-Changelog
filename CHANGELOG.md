## Changelog

### 4.0.0

* No longer defaults to output files. Instead writes output to std. out. Pipe to file if you need your output elsewhere
* Greatly simplified the docker image
* No longer supports Mercurial
* No longer supports PDF creation. You can still use external tools to turn the output into a pdf
* Added support for filter-paths in [Issue #130](https://github.com/Praqma/Praqmatic-Automated-Changelog/issues/130)

### 3.x versions

**Incompatible with versions 2.x and earlier - see the [migration guide](docs/Migrating_2.X.X_to_3.X.X.md) for more information**

* Removed all `date` related parameters
* Removed deprecated `--sha` parameter (has been replaced with `from`)

### 2.x versions

**Incompatible with versions 1.x and earlier - see the [migration guide](docs/Migrating_1.X.X_to_2.X.X.md) for more information.**

* Support for report templates
* Support for JIRA

### 1.x versions

* Support for 'none' report - changelog without task system interaction

### 0.x versions

_Initial release and proof-of-concept_
