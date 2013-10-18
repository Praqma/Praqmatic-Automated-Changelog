# Praqmatic Automated Changelog (PAC)

Simple proof-of-concept for a automated, but pragmatic, changelog system.

Compared to other ways of extracting changes, typically from the SCM commit messages, this is based around getting information from another source by looking up references to such in the SCM commit message.

Currently proof-of-concepts handles

* git dvcs, Trac case management system (unix and windows)
* git dvcs, Fogbugz case management system (unix and windows)
* hg dvcs, Trac case managenent system (unix)
* hg dvcs, Fogbugz case managenebt system (unix)

The current script can output in three different formats

* html
* pdf
* markdown (default, always on)

## Settings file example (Fogbugz and git)

Below is an example of an example that uses FogBugz and git. If you want to use trac replace `:fogbugz:` with `:trac:`

    :general:
      date_template: "%Y-%m-%d"
      changelog_name: "Changelog"
      changelog_formats:
        - "html"
      changelog_css:
    
    :fogbugz:
      fogbugz_url: "https://my.fogbugz.site"
      fogbugz_usr: my@companymail.net
      fogbugz_pwd: p455w0rd
      fogbugz_fields: "sTitle,sStatus,sUrl,sCategory,sTags,sPriority,sReleaseNotes"
    
    :vcs:
      type: git
      repo_location: "/home/myuser/myproject/repo"
      usr:
      pwd:


## Usage examples
Show help

    ./pac.rb -h
    
Get commits using tags. tail tag provided

    ./pac.rb -t Release-1.0 --settings=./pacfogbugz_pac_settings.yml --outpath=/usr/share/changelog #Using tags.

Get commits using tags, head and tail tags provided

    ./pac.rb -t Release-1.0 Release-2.0 --settings=./pacfogbugz_pac_settings.yml --outpath=/usr/share/changelog

Get commits using time

    ./pac.rb -d 2013-10-01 --settings=./pacfogbugz_pac_settings.yml --outpath=/usr/share/changelog #Using date as selector.
    
Get all commits since latest point release. Given a specified pattern. Default is 'tags'

	./pac.rb -t LATEST --settings=./pacfogbugz_pac_settings.yml --outpath=/usr/share/changelog --pattern='tags/Release-1.*'

## Prerequisites
If you are going to be using the tool to generate PDF files which we use kramdown and pdfkit to generate you'll need to run the following command and a linux machine

`sudo apt-get install libxslt-dev libxml2-dev`

`sudo apt-get install wkhtmltopdf`


Also you'll need to install the gems specified in the Gemfile in order to get it working. At the current state Linux support is much better than windows

## Contributors
* Hugo Leote (hleote@praqma.net)
* Martin Georgiev (mvgeorgiev@praqma.net)
* Mads Nielsen (man@praqma.net)
* Bue Petersen (bue@praqma.net)