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

PAC has a flexible internal design, making it easy to extend with support for additional report formats, task management systems and so on.

See our [roadmap](/roadmap.md) for future plans.

## Demo

Take PAC for a quick spin with our demo scripts (requires Linux Bash, zip and docker).

**Generate changelogs using only SCM commits**

    ./demo_setup_docker.sh

This demo generates some reports based only on SCM commits of a small test repository.
After running the demo, look for the reports in `demorepo/idReportTestRepository/default.[pdf|html|md|` and `demorepo/idReportTestRepository/ids.md`

**Generates changelogs using data from JIRA**

    ./demo_setup_docker_jira.sh

This demo generates some reports through a custom template which uses data from SCM commits of a small test repository and a spun-up JIRA instance.

**Try a manual demo using the pac-manuscript repository**

The [pac-manuscript](https://github.com/praqma-test/pac-manuscript) repository is a short demonstration of how PAC works with Jira, where each step of gradually introduces more of the PAC configuration and features.
The readme in the repository explain how to use it, but the pac-maunscript repository is a imaginary example on a source code repository matching some Jira issues where the series of commits are introducing new PAC features.

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
	  date_template: "%Y-%m-%d"

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

Basic usage examples for the PAC Ruby script:

Show help

    ./pac.rb -h
    
Get commits using tags from "Release-1.0" tag to "HEAD":

    ./pac.rb -t Release-1.0 --settings=./default_settings.yml

    ./pac.rb from-latest-tag "Release-1.0" --settings=./default_settings.yml

Get commits using tags from "Release-1.0" to "Release-2.0"

    ./pac.rb -t Release-1.0 Release-2.0 --settings=./default_settings.yml 

    ./pac.rb from-latest-tag Release-1.0 --to Release-2.0 --settings=./default_settings.yml

Get commits using latest tag of any name: 

	./pac.rb from-latest-tag "*" --settings=./default_settings.yml

Get commits using time

    ./pac.rb -d 2013-10-01 --settings=./default_settings.yml

    ./pac.rb from-date 2013-10-01 --settings=./default_settings.yml      

The above getting started is only a simple example, so to utilize all the features in PAC you can dive into the following sections.

*Verbosity:*  As most unix tools, pac supports `-v` and `-q` to indicate that it should be more or less verbose in its output. Both switches can be provided multiple times (e.g. `-vv`) for increased effect. One can be used to cancel out the other, which can be useful for instance for reverting behaviour added in a bash alias.

We recommend using the PAC docker image, as described below in [Running PAC](#running-pac). The basic usage examples then becomes like described in [Usage](docs/using_the_pac_docker_image.md#usage) in [Using the PAC Docker image](docs/using_the_pac_docker_image.md).

### Configuring PAC

All available configurations option for PAC is described in the [Configuration](docs/configuration.md) section.

### Writing templates

Information on how to write templates for the changelog and use the extracted data can be found in the [Templates](docs/templates.md) section.

### Take PAC for a spin

You can try PAC using our PAC Docker container and a zipped github repository we use for testing. See [Try PAC with a test repo and PAC docker image](docs/try_pac_with_test_repo_and_docker.md)

### Running PAC

We recommend to use our supplied [`praqma/pac`](https://hub.docker.com/r/praqma/pac/) Docker image so you avoid configuring a Ruby environmnet yourself.

* See [Using the PAC docker image](docs/using_the_pac_docker_image.md).

If you like to configure your own Ruby environment and run PAC as simple Ruby script (`pac.rb`) follow instruction below for Linux (Ubuntu) or Windows.

* PAC requires Ruby version 2 or later. Currently tested with version 2.3.0 of Ruby.


### Run PAC on Linux (Ubuntu)

Configure your Linux Ruby environment to run PAC and get PAC from sources:

Prerequisites:

 * Ruby version 2 (you can see specific version in the [PAC docker image file](Dockerfile))
 * The `bundler` Ruby Gem
 * Native libraries - for Ubuntu they are: `sudo apt-get install cmake libxslt-dev libxml2-dev wkhtmltopdf`    

Then get and use PAC:
              
1. Clone the pac repository to your local machine: `git clone https://github.com/Praqma/Praqmatic-Automated-Changelog.git pac`
2. Optionally check-out the `latest` tag or a specific release tag if you don't want bleeding edge.
3. Change directory to `pac` (the git clone) and run the command `bundle install` to install all the used Ruby Gems.
4. Optionally add a symlink to `pac.rb` in your path, for example: `ln -s [your install directory]/pac.rb /usr/bin/pac` 

That's it. Test your installation by executing pac: `pac`. If you get a help screen the installation was successful.

### Run PAC on Windows

Detailed instructions can be found in [Installing PAC on Windows](docs/windows_instructions.md).


## Support and maintenance

* PAC is maintained in the scope of [JOSRA](http://www.josra.org/).
* Issue and work tracking is done using [Github issues](https://github.com/Praqma/Praqmatic-Automated-Changelog/issues)
* Support requests and questions can be created as Github issue or send us an mail on support@praqma.net
* Our roadmap is availbe in [roadmap](/roadmap.md)

## Changelog

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

* Mads Nielsen (man@praqma.net)
* Bue Petersen (bue@praqma.net)
* Andrius Ordojan (and@praqma.net)
* Thierry Lacour (thi@praqma.net)
