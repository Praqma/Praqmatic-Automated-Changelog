---
maintainer: JKrag
---

# Praqmatic Automated Changelog (PAC)

Tool for creating automated, but pragmatic, changelogs.

PAC collects task references from SCM commit messages and creates changelog reports with additional information extracted from other systems, like your task management system.
Compared to other changelog solutions, PAC is very flexible and customizable. The design allows you to solve the problems of having an unchangeable SCM commit history with incorrect task references.

## Features

* Customizable change reports based on **Liquid templates**
* Collects task references in SCM commits from Git
* Extracts data from tasks systems that has an open rest interface that returns `Json` (Basic authenticaton supported)
* For task systems returning json data, all data can be used in the templates 
* Supports any textual output format. I.e you can write a template to produce css, markdown, html etc..
* Supports extracting data from multiple referenced tasks systems at once
* Supports creating changelog without data from external systems, basing it only on SCM commits
* Easily create reports for several different audiences using data from several sources
* Collects statistic on commits with and without task references

PAC has a flexible internal design, making it easy to extend with support for additional report formats, task management systems and so on.

## Demo

**Generates changelogs using data from GitHub**

    ./github_example.sh

This demo, using docker, generates some reports through a custom template which uses data from PAC's own repository where we mount the current directory and run PAC against it. If succesful, you should see a file in `reports/github-report.html`. 

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
	  :strict: false
	  :ssl_verify: true

	:templates:
	  - { location: templates/default_id_report.md, output: ids.md }

	:task_systems:
	  -
	    :name: none
	    :regex:
	      - { pattern: '/PAC\-(\d+)', label: none }

	:vcs:
	  :repo_location: '.'

Help writing regexp using Ruby IRB see this litle howto: [Howto write regexp using IRB](docs/howto_write_regexp_using_irb.md)

### Simple template

This example template simply lists the discovered issues as headers in a Markdown file. 

	# {{title}}
	{% for task in tasks.none %}
	## {{task.task_id}}
	{% endfor %}

## Usage

Basic usage examples for the PAC Ruby script:

Show help

    pac -h
    
Get commits using tags from "Release-1.0" tag to "HEAD":

    pac from Release-1.0 --settings=./default_settings.yml

    pac from-latest-tag "Release-1.0" --settings=./default_settings.yml

Get commits using tags from "Release-1.0" to "Release-2.0"

    pac from Release-1.0 --to Release-2.0 --settings=./default_settings.yml 

    pac from-latest-tag Release-1.0 --to Release-2.0 --settings=./default_settings.yml

Get commits using latest tag of any name: 

	pac from-latest-tag "*" --settings=./default_settings.yml

The above getting started is only a simple example, so to utilize all the features in PAC you can dive into the following sections.

*Verbosity:*  As most unix tools, pac supports `-v` and `-q` to indicate that it should be more or less verbose in its output. Both switches can be provided multiple times (e.g. `-vv`) for increased effect. One can be used to cancel out the other, which can be useful for instance for reverting behaviour added in a bash alias.

We recommend using the PAC docker image, as described below in [Running PAC](#running-pac). The basic usage examples then becomes like described in [Usage](docs/using_the_pac_docker_image.md#usage) in [Using the PAC Docker image](docs/using_the_pac_docker_image.md).

### Configuring PAC

All available configurations option for PAC is described in the [Configuration](docs/configuration.md) section.

### Writing templates

Information on how to write templates for the changelog and use the extracted data can be found in the [Templates](docs/templates.md) section.

### Running PAC on Linux

We recommend to use our supplied [`praqma/pac`](https://hub.docker.com/r/praqma/pac/) Docker image so you avoid configuring a Ruby environmnet yourself.

### Run PAC on Windows

Detailed instructions can be found in [Running PAC on Windows](docs/windows_instructions.md).

## Support and maintenance

* Issue and work tracking is done using [Github issues](https://github.com/Praqma/Praqmatic-Automated-Changelog/issues)
* Support requests and questions can be created as Github issue or send us an mail on support@praqma.net

## Changelog

### 4.x.x 

**Incompatible with versions 3.x.x and below** 

Refocus effort on optimitzing the docker image, removing dependencies and focusing on pure rest/json based tasks systems. 

* The `jira` task system has been renamed `json`
* Additional ssl options added. `--ssl-verify` and `--ssl-no-verify` options. Defaults to `true` to turn on ssl verification
* New authorization section. Supports token auth with GitHub token authentication and Basic Authentication for other systems 
* Also able to read .netrc files and using this information in Basic Authentication
* The configuration file is now interpreted using ERB. This means that you can use ERB code blocks like `<%= ...expression %>` in the configuration file. Useful if you want to have certain things happinging in a test environment and other things going on in a production envrionment. **WARNING** This can be a potential security issue

```

:task_systems:
  - 
    :name: json
    :regex:
      - { pattern: '/#(\d+)/', label: github }
    :query_string: "https://api.github.com/repos/Praqma/Praqmatic-Automated-Changelog/issues/#{task_id}"
    # Option 0: No auth. Just remove the :auth: section
    # Option 1: GitHub token authentication (Store in Environment variable)
    :auth:
      :github:
        :token: <%= ENV['GITHUB_API_TOKEN'] %> #
    # Option 2: Use basic authentication with values from somewhere
    :auth:
      :basic:
        :username: <%= ENV['BASIC_USR'] %>
        :password: <%= ENV['BASIC_PW'] %>
    # Option 3: Use .netrc file. Remember to mount the files to default location or known location when using docker.
    :auth:
      :basic:
        :netrc: /data/.netrc 
```

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

### Maintainers

* Mads Nielsen (man@praqma.net)
* Thierry Lacour (thi@praqma.net)
* Jan Krag (jak@praqma.net)
* Claus Schneider (cls@praqma.net)
