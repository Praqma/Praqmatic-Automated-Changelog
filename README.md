---
maintainer: JKrag
---

Issue tracking: 
[![Groomed](https://badge.waffle.io/Praqma/Praqmatic-Automated-Changelog.png?label=Status%20-%20workable&title=Groomed)](https://waffle.io/Praqma/Praqmatic-Automated-Changelog) 
[![Up Next](https://badge.waffle.io/Praqma/Praqmatic-Automated-Changelog.png?label=Status%20-%20up%20Next&title=UpNext)](https://waffle.io/Praqma/Praqmatic-Automated-Changelog) 
[![Work In Progress](https://badge.waffle.io/Praqma/Praqmatic-Automated-Changelog.png?label=Status%20-%20in%20progress&title=InProgress)](https://waffle.io/Praqma/Praqmatic-Automated-Changelog)
[![Issues](https://img.shields.io/github/issues/Praqma/Praqmatic-Automated-Changelog.svg)](https://github.com/Praqma/Praqmatic-Automated-Changelog/issues)


# Praqmatic Automated Changelog (PAC)

Tool for creating automated, but pragmatic, changelogs.

PAC collects task references from SCM commit messages and creates changelog reports with additional information extracted from other systems, like your task management system.
Compared to other changelog solutions, PAC is very flexible and customizable. The design allows you to solve the problems of having an unchangeable SCM commit history with incorrect task references.

![The workflow behind PAC for creating changelogs](/docs/process.png)

## Features

* Customizable change reports based on **Liquid templates**
* Collects task references in SCM commits from **Git** or **Mercurial (hg)**
* Extracts data from tasks systems like **Trac** and **JIRA** for the collected tasks
* For task systems returning json data, all data can be used in the templates 
* Supports **MarkDown**, **HTML** and **PDF** as report formats
* Supports extracting data from multiple referenced tasks systems at once
* Supports creating changelog without data from external systems, basing it only on SCM commits
* Easily create reports for several different audiences using data from several sources
* Collects statistic on commits with and without task references

## Getting started

You'll need to have some commits that reference tasks in one way or another, otherwise your changelog will look rather dull.
It could be Jira issues: "_Closing PAC-1337 change help text_" or "_PAC-1042: Fixed bug in date resolver_".

To start generating changelogs you'll need to:

 1. Either install PAC locally or use our provided Docker image
 2. Write a PAC configuration file for your project [see the simple configuration file below](#simple-configuration-file)
 3. Write a template for your change report [see the simple template below](#simple-template)

### Simple configuration file

This is an example of a simple configuration file. It collects task references from commits using the configured regex and create a changelog based on the configured template.
This simple example do not extract data from task systems.

	:general:

	:templates:
	  - { location: templates/default_id_report.md, output: ids.md }

	:task_systems:
	  -
	    :name: none
	    :regex:
	      - { pattern: '/PAC\-(\d+)', label: none }

	:vcs:
	  :type: git
	  :repo_location: '.'

More about configuration in [Configuration](docs/configuration.md).

Help writing regexp using Ruby IRB see this litle howto: [Howto write regexp using IRB](docs/howto_write_regexp_using_irb.md)

### Simple template

This example template simply lists the discovered issues as headers in a Markdown file. 

	# {{title}}
	{% for task in tasks.none %}
	## {{task.task_id}}
	{% endfor %}

More about templates in [Templates](docs/templates.md).


## Usage

Basic usage examples for the PAC Ruby script, run PAC with the `--help` parameter for usage explanation.

### Run PAC on Windows

Detailed instructions can be found in [Installing PAC on Windows](docs/windows_instructions.md).

## Support and maintenance

* PAC is maintained in the scope of [JOSRA](http://www.josra.org/).
* Issue and work tracking is done using [Github issues](https://github.com/Praqma/Praqmatic-Automated-Changelog/issues)
* Support requests and questions can be created as Github issue or send us an mail on support@praqma.net
* Our roadmap is availbe in [roadmap](/roadmap.md)

## Changelog

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

* Trac support
* Markdown, HTML, PDF


## Developer information

For details on design and development info see [Developer information](docs/developer_info.md)

See also [contributing file](/CONTRIBUTING.md).

### Contributors

* Mads Nielsen (mads.nielsen@eficode.com)
* Bue Petersen (bue@praqma.net)
* Andrius Ordojan (and@praqma.net)
* Thierry Lacour (thierry.lacour@eficode.com)
