# Blue print for Praqmatic Automated Changelog (PAC)

_A short description of an automated changelog system and some initial design considerations_.

Original idea from 2013/2014.

# Motivation

## A good changelog and a customized release note
A changelog is a description of changes between two different revisions of a software component. Such a list of changes is often generated bases on the VCS commit messages and revision history. For example it can be a list of patches to the source code, or it can be a list of commit messages with dates or authors.

An important aspect of a changelog is that it should be targeted to a specific groups of readers. It could be project managers or it could be end-user and customers. This might requires two very different worded documents, though they are in the end based on the exact same information from the changelog.
Often the document for customers are called _release note_.

Determining changes between two different revision of a software component requires some kind of traceability in your software development. 
Changes narrows down to what are changed in your VCS. As this might not be descriptive enough we often tries to combine and relate it with information from the task management system. Often developers must supply further and elaborated information as well.
All this is much easier with some kind of traceability between commits and tasks. For example a mention of one or more  tasks in the commit message related to that commit.

Combining information will require manual work, because most of the information often needs some kind of rewriting or explanation to the non technical readers. Commits are targeted developers and often not a god explanation. Tasks are supposed to be more elaborate and at a higher level, but also often include developers internal notes.
So basically there is lots of information, but rewriting is often needed to make some group of readers understand the change log (or release note).

## Continuous changelog
In a continuous delivery world we are always ready to release, and as a release also include a changelog, or a release note, and other kind of documentation, these documents must also always be available continuously.

We recognize that the good release note requires some manual work and effort as above suggested, so this mean we must ensure most work are done continuously and up-front on that part as well:

  1) make sure the needed information for changelog is always available
  2) combine several sources of information with traceability
  3) make it easy and fully automated to generate the needed document(s)


In an agile team doing continuous integration there will be many small commits and they will always reference a task. Such a commit, with a task reference, and preferably also a short description written for developers by developers will give the basic information to the changelog. Information that will always be there in the VCS.

With traces like a reference to a task, we can start to collect more information than the commit message itself.
Information can be gathered from the task management system, even maybe from special fields in such a task, and we can start working with that information. Supply more fields with information, rewrite etc. which is often not possible with the commit messages.

Finally we need the generation of the changelog to run automated, only supplying configuration for example about output format and maybe information filtering options.



## Yet another system?
A small research on what is already out there for git, showed that there is lots of ways to create changelog in nice ways from git commit messages, but there isn't a solution to extract the information from other resources based on a references in the commits.

See "Research reference URLs" below

To be able to collect more information that the commit message is interesting for three reasons:

  §1 We will be able to edit the changelog easily, by editing the referenced information residing in another system, and redoing a changelog creation. It you don't like the idea that the changelog can change, it is always possible to store it in the artifact system together with the released artifacts.

  §2 We will avoid creating to many restrictions on commit messages, that might keep the developers from committing often if they find the requirement for commit messages to cumbersome. Most of the found change logs solutions is based on nicely created commit messages, with formatting rules etc. This will typically have an effect on how often developers commit. We will recommend to commit often, even thought the commit message is only a reference to a ticket or case and a short say description. Preparing nice commit messages and rewriting history before developers share their code, can have a negative effect on continuously integrating. Even though we consider only one nicely written commit message when rebasing to say master branch, it is still a time consuming task to write long commit message and they are not easily editable (§1).

  §3 There is often another group of readers for a release note, than for commit messages, so having several aggregated (manually or automated) types of information is interesting.


### Research reference URLs

http://trac-hacks.org/wiki/ChangeLogMacro

http://trac-hacks.org/wiki/TracTicketChangelogPlugin

http://changelog.complete.org/archives/694-trac-git

http://www.turnkeylinux.org/trac

http://blog.webfaction.com/2011/05/trac-and-git-two-new-best-friends/

http://nileshgr.com/2012/08/08/trac-and-git-the-right-way

http://trac-hacks.org/wiki/plugin

http://unethicalblogger.com/2008/11/23/git-integration-with-hudson-and-trac.html

http://stackoverflow.com/questions/14984276/automated-changelog-using-jira-github-and-jenkins

http://stackoverflow.com/questions/9078474/jenkins-how-to-save-changelog-for-build?rq=1

http://stackoverflow.com/questions/3523534/good-ways-to-manage-a-changelog-using-git

http://grahamweldon.com/posts/view/automatically-generating-changelogs-from-git-for-cakephp

http://jenkins.361315.n4.nabble.com/Toward-automated-changelog-generation-td4663249.html

http://stackoverflow.com/questions/14960674/changelog-generation-from-github-issues

https://coderwall.com/p/5cv5lg

http://git.661346.n2.nabble.com/Generating-GNU-style-Changelog-from-git-commits-td6279979.html

http://blog.cryos.net/archives/202-Git-and-Automatic-ChangeLog-Generation.html

http://git.savannah.gnu.org/gitweb/?p=gnulib.git;a=blob_plain;f=build-aux/gitlog-to-changelog;hb=HEAD



# Use cases
## Current state changelog
The project manager want to know which new features are available in release candidate 1.2.0-RC77. Most tasks in the sprint are closed but not all changes have reached the release candidate QA level yet.
He goes to Jenkins an execute the RC changelog job. A document is short after available with a summary of features from closed tasks referenced by commit since last release tag. The summary of a feature comes from the subject- and summary- field on the tasks in the task management system.


## Release note
As a build step on the release job, a release note document is generated and stored together with the released software binaries.
The release document is several pages long, and contain information in sections about bugfixes, new features and improvements.
Each section contain information from each closes task referenced by a commit since last release tag.
As it is a release note, the customer subject- and customer summary-field is automatically used for the document. Developers internal comments on the tasks etc. are not included.


## Sprint cases and statistics.
An agile team only work on planned tasks in current sprint. This means there should not be any commit with a relation to a task reference, and the commit message should contain a task reference then.
Extracting an automated changelog, with the statistic configuration enabled, will summary all referenced tasks, their type like feature, bugfix etc.
Count can be shown on closed, opened etc.
Further a summary of commits without a reference is shown sorted by developers, date and the commit message.
That is a very visible way to show missing references in commits, if you do not want to enforce the policy.


## Internal release note
During an investigation on a regression problem, where we tries to narrow down a long period for introducing a bug in the temperature sensor, we extract an internal changelog.
The changelog a list of all commits, their tasks, status, type, and the subject- and summary-field that in a more non-technical manner explain the change.
Both developers and technical staff can understand the changelog, and pinpoint interesting commit to investigate further or specific software versions to compare.


## Command line examples
We imagine the changelog can be called from command line, and then be used on both Jenkins build steps or by developers.
Examples could be:

	./PAC --tag 1.0.0 --tag 1.1.2
	./PAC --date 2013-05-31 --date 2013-06-11
	./PAC --vcs 27ab8ff --vcs 87e8ef

In all cases two arguments are used for respectively a start- and a end-references. We can use tags, dates or VCS revision references like Git SHAs.
The end reference can be omitted, for which we assume latest tag, today's date or latest commit.
Optionally other arguments like a configuration file is accepted.

If no arguments are given, the arguments comes from a configuration file at a default location.



# Design proposal

![The flow from task, to commit, to changelog](docs/process.png "The flow from task, to commit, to changelog")

The picture shows part of the process where we generate a changelog from Jenkins and commits is used to find task references, which again is used for finding information to the changelog before the document written.

![Modular design where interfaces can be specified in configuration](docs/designview.png "Modular design where interfaces can be specified in configuration")

A modular design is crucial so we can easily support a new VCS or document type.
We imagine a modular design with a ___PAC core___ responsible for the overall flow, configuration and input argument parsing etc. Then a plug-able module for the primary tasks that hooks up on the core through the interface.

## Interface design and module considerations

### Interfaces in general
PAC core works with a data format that can handle data for every task. We need to define such data structure.

For every interface, the PAC core side will receive data in that format. On the other side of the interface the module is responsible for understanding the PAC core format.

## VCS interface and module
The VCS module will on request from PAC core gather the needed information and return it to PAC in a general data format.

An implementation of a VCS module must be able to receive a request containing a definition of start- and end-range for the commits that must be gathered.

The VCS module must hand PAC core the references in commits in a format common for all VCS modules, thus the VCS module is responsible for understanding the commit messages references.
It can be a regular expression in the VCS module configuration.

## Document interface and module
PAC hands the document module a common defined document format (could be markdown) from where the document module is responsible for the translation into a format the document module supports.

The document module will probably mostly be about gluing already existing software together.

## Task info module and interface
The PAC core will receive a change list from VCS module, also containing a task reference list. The task reference list is defined in common format for all task modules, so that every task module then have the responsibility to understand the PAC task reference list and translate it to task system request through a task system API or database query.


## Magic words and commit message format
We prefer a solution, where a githook is enforcing every commit to have a task reference. Such a githook is not considered part of the PAC, but it will influence both the VCS and task info module.

With such a githook there is a rule for writing commit messages, and especially how to reference a task in the message. For example "Ticket#666: Interface implementation for invoice line".
In the example the task reference found will be task no. 666, so the VCS module will be responsible for understanding the commit message and extract the task number.
As this may vary a lot between projects, it is suggested that such parsing is part of a configuration. For example the Git VCS module will accept a parsing rule as a configuration option.



# Technologies

We are quite set that the system should be implemented in Ruby. It is an obvious choice, though not the only one.

  * Ruby have lots of libs for databases, REST APIs which will be useful as we need to glue lots of systems together
  * Ruby is a good choice for parsing and working with text and regular expressions
  * Ruby have been used much already for CI/CD scripting, tooling etc.
  * It should be cross platform, which Ruby easily can support

Regarding distribution of the Praqmatic Automated Changelog there will be several options like just plain source code, Rake or Gem if using Ruby.


# Road map
## Prototype
A prototype should implement a PAC core, with at Git VCS module, a document module for pdf output and a task info module for Trac task management system.

Included in the prototype must be simple input parameter parsing, accepting SHA, tags and dates. An internal PAC data structure must be suggested, and the PAC internal document data structure should be suggested (is markdown a solution?).

The Trac information module must extract subject and summary field for referenced tickets, but comment fields with internal developer commit must be left out.
Ticket type like feature or improvement can be left out in prototype.

Prototype is developed and delivered with a CI/CD setup, with a large set of tests targeted at core functionality and the basic usage work flows.


## 1. iterations
The next iteration after the prototype will supply further features:
  * Task module must handle task types and task status (closed, open ...)
  * Filtering functionality for documents based on task reference criterias such as status of tasks and type of tasks
  * A simple Mercurial VCS module must be implemented to show prototype design flexibility
  * A new document module outputting to HTML or another format must be implemented to show design flexibility
  * A new task module for Fogbugz should show design and architecture flexibility for the prototype.
  * Design and architecture improvement and fixes


## Jenkins plugin
There will come a time where it is obvious to put a Jenkins plugin on top of the changelog.

## GUI
A GUI, maybe Ruby Qt can later be considered.
