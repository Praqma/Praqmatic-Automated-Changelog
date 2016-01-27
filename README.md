# Praqmatic Automated Changelog (PAC)

Tool for creating automated, but pragmatic, changelogs.

Compared to other ways of extracting changes, typically from the SCM commit messages, this is based around getting information from another source by looking up references to such in the SCM commit message.

Currently proof-of-concepts handles

* git dvcs
* hg dvcs

You can output in any format you like using the liquid templating language. We have the option to turn html into pdf with the pdf switch in the template setup.

## Settings file example (Jira)

Below is an example of an example that uses Jira.

	:general:
	  date_template: '%Y-%m-%d'

	:templates:
	  - { location: templates/default_id_report.md, output: ids.md }
	  - { location: templates/default.md, output: default.md }
	  - { location: templates/default_html.html, pdf: true, output: default.html }

	:task_systems:
	  - 
	    :name: none
	    :regex:
	      - { pattern: '/Issue:\s*(\d+)/i', label: none }
	      - { pattern: '/Issue:\s*(none)/i', label: none }
	      - { pattern: '/(#\d+)/', label: none }
	      - { pattern: '/us:(\d+)/', label: none }
	    :delimiter: '/,|\s/'
	  -
	    :name: jira
	    :query_string: "http//your.server.hostname/jira/rest/api/latest/#{task_id}"
	    :usr: "user"  
	    :pw: "password"
	    :regex:
	    - { pattern: '/PRJ-(\d+)/i', label: jira }      
	  -
	    :name: trac
	    :trac_url: "https://my.trac.site"
	    :trac_usr: "user"
	    :trac_pwd: "pass"
	    :regex:
	    - { pattern: '/Ticket-(\d+)/i', label: trac }  

## Special case with none task system

We also implemented a setting if you do not wish to talk to your external task tracking, but still want to get a list of the id's associated with your task management system. Replace
`:fogbugz:` or `:trac:` with `:none:`

	:none:
		regex:
		 - '/.*Issue:\s*(?<id>[\d+|[,|\s]]+).*/im'
		delimiter: '/,|\s/'

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
	
### Try the ID report feature

If you want to try out the ID report feature, you can use one of the test resources we already have and uses for testing.

Here is how...

1. unzip the test repositoyr `test/resources/idReportTestRepository.zip` into `test/resources/`
1. Run from the root of your current directory (where the README file you currently are reading, not the unzipped one) the following command: `./pac.rb -d 2012-01-01 --settings=test/resources/idReportTestRepository_settings.yml`
1. Look at the ids.md file created in your current directory  
	

## Usage examples

_Note there is a docker container available also, which makes the tools and environment setup easier: [Praqma/docker-pac](https://github.com/Praqma/docker-pac)_

Show help

    ./pac.rb -h
    
Get commits using tags. tail tag provided

    ./pac.rb -t Release-1.0 --settings=./pacfogbugz_pac_settings.yml

Get commits using tags, head and tail tags provided

    ./pac.rb -t Release-1.0 Release-2.0 --settings=./pacfogbugz_pac_settings.yml

Get commits using time

    ./pac.rb -d 2013-10-01 --settings=./pacfogbugz_pac_settings.yml 
    
Get all commits since latest point release. Given a specified pattern. Default is 'tags'

	./pac.rb -t LATEST --settings=./pacfogbugz_pac_settings.yml --pattern='tags/Release-1.*'

Note that when using tags with a pattern that matches multiple tags, the latest is always used. The comparison is always compared to HEAD.

## Using the Praqma/docker-pac container

* [Praqma/docker-pac imagefile](https://github.com/Praqma/docker-pac)
* [praqma/pac image](https://registry.hub.docker.com/u/praqma/pac/)

The following usage example are actually based on test repository we supply with `idReportTestRepository` in `test/resources/idReportTestRepository.zip` so you can easy try the following yourself.

* First create the `testing-PAC` folder somewhere on your local machine and accessible from docker.
* enable display as described just above
* clone latest [PAC](https://github.com/Praqma/Praqmatic-Automated-Changelog) (tagged release version) to the `testing-PAC`-folder as `PAC-0.9.0`

Your should now have the following:

```
testing-PAC
testing-PAC/PAC-0.9.0
```

* unzip the testing repository `test/resources/idReportTestRepository.zip` into `testing-PAC/idReportTestRepository`
	* `unzip PAC-0.9.0/test/resources/idReportTestRepository.zip`


Then you have:


```
testing-PAC
testing-PAC/PAC-0.9.0
testing-PAC/idReportTestRepository
```

* copy and edit the PAC configuration file to match repostiroy location:
  * `cp PAC-0.9.0/test/resources/idReportTestRepository_settings.yml idReportTestRepository/`
  * edit the line `repo_location:` to match `repo_location: "idReportTestRepository"`
* Then run docker pac container v2 like this from the `testing-PAC` directory:
  `docker run -v $(pwd):/data -v /tmp/.X11-unix:/tmp/.X11-unix:ro -e DISPLAY=$DISPLAY praqma/pac:v2 ruby PAC-0.9.0/pac.rb -s f9a66ca6d2e6 --settings=idReportTestRepository/idReportTestRepository_settings.yml`

_and you will get ids.md (the ID report) as output_ and see that PAC is able to run and use the toolstack supplied.


## Prerequisites
If you are going to be using the tool to generate PDF files which we use kramdown and pdfkit to generate you'll need to run the following command and a linux machine

`sudo apt-get install libxslt-dev libxml2-dev`

`sudo apt-get install wkhtmltopdf`


Also you'll need to install the gems specified in the Gemfile in order to get it working. At the current state Linux support is much better than windows

_Note there is a docker container available also, which makes the tools and environment setup easier: [Praqma/docker-pac](https://github.com/Praqma/docker-pac)_



## Contributors
* Hugo Leote (hleote@praqma.net)
* Martin Georgiev (mvgeorgiev@praqma.net)
* Mads Nielsen (man@praqma.net)
* Bue Petersen (bue@praqma.net)



## Developer info

### Program flow

The general program can be described this way

1. The Vcs module obtains a list of commits, given user supplied repository and start/finish. 
2. The Vcs turns raw commits into PACCommit model objects and returns an array of these in the `PACCommitCollection`
3. The Core module then produces a bare bones `PACTaskCollection` with tasks, the tasks only has the id property.
4. After the bare `PACTaskCollection` has been generated, each task system applies it's decorator(s) to the task adding additonal data.
5. The `PACTaskCollection` is then passed to the Liquid template engine and the outputs are produced

 
### Object model 
The principal model in PAC consists of the following ruby `Modules` 
* Core
* Vcs
  * GitVcs
  * MercurialVcs
* Model
  * PACTask
  * PACTaskCollection
  * PACCommit
  * PACCommitCollection
* Report
  * Generator
* Task
  * JiraTaskSystem
  * TracTaskSystem

*Module: Core*

The responsibility of `Core` is to combine the data found in the `Vcs` module and using this information together with the information from the `Task` module
to produce a gross list of tasks discovered. 

The `Core` module is used in the _Main_ method of the pac.rb script.   

*Module: Vcs*

The `Vcs` module handles the interaction with the chosen VCS (Usually Git). Given a set of parameters it will return a list of commits in the form of `PACCommitCollection` model object.

*Module: Model*

The `Model` module contains all the object models needed. We have the folling, the names make it clear what their responsibilities are
* PACCommit
* PACCommitCollection
* PACTask
* PACTaskCollection

The `PACTaskCollection` has a method to add _n_ tasks to the list. If the task was already added, based on the unique id, then the commits of the two tasks are 
merged, resulting in 1 task, with the extra commits from the other tasks. This happens if a task is referenced in multiple commits. The _uniqueness_ is implemented in the `PACTaskCollection`. 

In order to ensure that, the _==_ (equals) method on the `PACTask` has been overriden, to only take into account the id of the task when determining equality.

The `PACTask` object has a collection of associated commits. It also holds references to the names of which task systems the task applies to. Also labels tied to the commit are also applied to the task, so that the tasks can be sorted by their labels. A `PACTask` can contain multiple labels. 

*Module: Report*

The `Report` module has one class. The generator class that produces the output files. This generator needs to know a complete lists of tasks, the commits involved and a list of user defined templates to generate the report from. Currently we only have Jekyll generation capabilities.  

*Module: Task*

The `Task` module is responsible for applying the appropriate decorators for the task system. It is entirely possible to apply more than one decorator.
The module expects a list of tasks, the id of each task is used to query the task system and add additional info to the task.


### Notes on Liquid

In order for liquid to produce locals for the template, you need to implement a `to_liquid` method. This method should return a hash whose keys can be used 
as parameters in the template.

Here is an example

    def to_liquid
      { 
        'task_id' => @task_id, 
        'commits' => @commit_collection, 
        'attributes' => attributes,
        'label' => label
      }
    end


### Tests

Tests can be easily executed by running `rake test` for unit tests and `rake functional_test` for functional tests. Rake has been 