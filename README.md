---
maintainer: JKrag
---

Issue tracking:
[![Issues](https://img.shields.io/github/issues/Praqma/Praqmatic-Automated-Changelog.svg)](https://github.com/Praqma/Praqmatic-Automated-Changelog/issues)

# Praqmatic Automated Changelog (PAC)

Tool for creating automated, but pragmatic, changelogs.

PAC collects task references from SCM commit messages and creates changelog reports with additional information extracted from other systems, like your task management system.
Compared to other changelog solutions, PAC is very flexible and customizable. The design allows you to solve the problems of having an unchangeable SCM commit history with incorrect task references.

## Features

* Customizable change reports based on **Liquid templates**
* Collects task references in SCM commits from **Git**
* Extracts data from tasks systems which return json and has support for basic authentication.
* For task systems returning json data, all data can be used in the templates 
* Supports **MarkDown**, **HTML** as report formats. You can always convert this to PDF through browser or other tools
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

## Developer information

For details on design and development info see [Developer information](docs/developer_info.md)

See also [contributing file](/CONTRIBUTING.md).

### Contributors

* Mads Nielsen (mads.nielsen@eficode.com)
* Bue Petersen (bue@praqma.net)
* Andrius Ordojan (and@praqma.net)
* Thierry Lacour (thierry.lacour@eficode.com)
