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
      regex:
    	- '/[Case|\[Case\]|fixed]\s(?<id>([0-9]+))+/i'
    	- '/(?<id>JENKINS-[0-9]+)/i'
    
    :vcs:
      type: git
      repo_location: "/home/myuser/myproject/repo"
      usr:
      pwd:

## Special case with none task system

We also implemented a setting if you do not wish to talk to your external task tracking, but still want to get a list of the id's associated with your task management system. Replace
`:fogbugz:` or `:trac:` with `:none:`

	:none:
		regex:
		 - '/.*Issue:\s*(?<id>[\d+|[,|\s]]+).*/im'
		delimiter:
		- '/,|\s/'

The above example will match issue reference, in the commit message and return them as capturing group 'id'. If your regexp returns a match, that needs to be split you can use the _optional_ delimiter regexp to split the 'id' match data.

In the source code we use the regexp you give in the configuration like this:

	ids = /.*Issue:\s*(?<id>[\d+|[,|\s]]+).*/im.match(commit_message)[:id]
	list_of_ids = ids.split(/,|\s/)

(in the real implementation we also handles if there are not matches or the split doesn't give more ids)

__Note:__ You typically want the delimiter to be part of your regexp also when collecting references. The delimiter is just needed to be repeated to ensure the program can determine what is delimiter and what not.

When you construct your regexp, the IRB is helpful. 

Example one, checking the above example regexp works as we expect:

	irb(main):147:0> commit_message="Commit message header
	irb(main):148:0" 
	irb(main):149:0" More lines in commit message.
	irb(main):150:0" Issue: 12"
	=> "Commit message header\n\nMore lines in commit message.\nIssue: 12"
	
	irb(main):166:0> ids = /.*Issue:\s*(?<id>[\d+|[,|\s]]+).*/im.match(commit_message)[:id]
	=> "12"
	irb(main):167:0> list_of_ids = ids.split(/,|\s/)
	=> ["12"]

Another example on how the regexp matches:

	irb(main):174:0> commit_message="Commit message header short
	irb(main):175:0" 
	irb(main):176:0" Issue: 14, 45, 3"
	
	irb(main):178:0> ids = /.*Issue:\s*(?<id>[\d+|[,|\s]]+).*/im.match(commit_message)[:id]
	=> "14, 45, 3"
	irb(main):179:0> list_of_ids = ids.split(/,|\s/)
	=> ["14", "", "45", "", "3"]
	

	

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
