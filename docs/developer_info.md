# Developer information

This file collects all information, considerations, design and architecture relevant information relevant to consider if developing PAC.

The original idea behind PAC can be found in the BluePrint dated back in 2013: [Blue print for Praqmatic Automated Changelog (PAC)](BluePrint.md)

## Releasing with a build and delivery pipeline 

* [Our build and delivery pipeline](http://code.praqma.net/ci/view/Open%20Source/view/Praqmatic-Automated-Changelog/)
* Releasing are done from the pipeline manually executing the release job.
 * A release will tag a version of the Github repository - see [releases on Github](https://github.com/Praqma/Praqmatic-Automated-Changelog/releases)
 * A PAC docker image will also be released and pushed with the same version to [`praqma/pac`](https://hub.docker.com/r/praqma/pac/)


Consideration on how we release PAC docker images is summarized in [Release strategy and design for PAC docker image](docker_image_release_strategy.md)

Thoughts on how to report version from PAC itself is explained in [Version reporting in PAC](versioning.md)

## Developer environment

You can use the Docker image to build and test PAC application, that way you do not need to install ruby on your own local machine if you want to extend it. You can run any arbitrary command using this docker command, this example below executes `rake test` and mounts PAC into the container as a volume. Execute this while while your are in the root of your local PAC repository clone:

`docker run --entrypoint=/bin/sh --rm -v $(pwd):/data praqma/pac:snapshot -c rake test`  


## Architecture

The original high-level overview component diagram still stands and should be followed. Basically PAC have core handling the flow, then modules for VCS, task systems, writing reports etc.


![The original design view on PAC](/docs/designview.png)


### Program flow

The general program can be described this way

1. The Vcs module obtains a list of commits object from [Rugged](https://github.com/libgit2/rugged), given user supplied repository and start/finish. 
2. The Vcs turns the [Rugged](https://github.com/libgit2/rugged) commits into PACCommit model objects and returns an array of these in the `PACCommitCollection`. The commit collection functions just like a regular list in Ruby, you can add stuff to it with the `add` method etc. Refer to the source code for information. 
3. The Core module then produces a bare bones `PACTaskCollection` with tasks, the tasks only has the id property.
4. After the bare `PACTaskCollection` has been generated, each task system applies it's decorator(s) to the task adding additonal data.
5. The `PACTaskCollection` is then passed to the Liquid template engine and the outputs are produced

Since [Liquid](https://shopify.github.io/liquid/) has a very strict object model, and no ruby code is allowed inside templates (you can use simple filters and chain theme together with th | (pipe), we need to transform the data in `PACTaskCollection` so that Liquid understands it. For example, symbols as has keys are not allowed in the template.

You are also required to bind Model objects to Liquid locals to be used in templates see the `Notes on Liquid` section.

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

*Module: Logging*

Logging contains some helper functions for printing output in accordance with
the user selected verbosity level. 
The important part is the `Logging.verboseprint(int, args*) function.
Use this anywhere you want to output to console. It will print the provided 
message (args*) if, and only if, the user has specified a verbosity equal to or above the number specified. 

Example: `Logging.verboseprint(1, "Debug message")` will be printet only if pac was called with at least `-v`.

*Module: Vcs*

The `Vcs` module handles the interaction with the chosen VCS (Usually Git). Given a set of parameters it will return a list of commits in the form of `PACCommitCollection` model object.

*Module: Model*

The `Model` module contains all the object models needed. We have the folling, the names make it clear what their responsibilities are
* PACCommit
* PACCommitCollection
* PACTask
* PACTaskCollection

`PACCommit` is just a data structure that encapsulates the `Rugged` commit. It has the `referenced` member which is an indication whether or not this commit had a task reference.

```
class PACCommit
  def initialize(sha, message = nil, date = nil)
	@sha = sha
	@message = message
	@referenced = false
	@date = date
  end
end
``` 

The `PACTaskCollection` has a method to add _n_ tasks to the list. If the task was already added, based on the unique id, then the commits of the two tasks are 
merged, resulting in 1 task, with the extra commits from the other tasks. This happens if a task is referenced in multiple commits. The _uniqueness_ is implemented in the `PACTaskCollection`. The `PACTaskCollection` basically encapsulates a Ruby list, with some additional logic to ensure uniqueness.

```
class PACTaskCollection
  def initialize
    @tasks = []
  end
end 
```

In order to ensure uniqueness, the _==_ (equals) method on the `PACTask` has been overridden, to only take into account the id of the task when determining equality.

The `PACTask` object has a collection of associated commits. It also holds references to the names of which task systems the task applies to. Also labels tied to the commit are also applied to the task, so that the tasks can be sorted by their labels. A `PACTask` can contain multiple labels. 

```
class PACTask
  def initialize(task_id = nil)
    #Lookup key for task management system
    @task_id = task_id      
    #Commits tied to this task
    @commit_collection = PACCommitCollection.new 
    #Data from task management systems(s)
    @attributes = { }
    #Key that determines which system(s) we need to look in for data
    @applies_to = Set.new      
    #Assigned label. Used in templates so that you can group your tasks using labels.
    @label = Set.new
  end
end
```

*Module: Report*

The `Report` module has one class. The generator class that produces the output files. This generator needs to know a complete lists of tasks, the commits involved and a list of user defined templates to generate the report from. Currently we only have Jekyll generation capabilities.  

*Module: Task*

The `Task` module is responsible for applying the appropriate decorators for the task system. It is entirely possible to apply more than one decorator.
The module expects a list of tasks, the id of each task is used to query the task system and add additional info to the task.

#### Model example

The first method we use is the one that traverses your git commit messages, for example, this is the output of the `Core.get_delta` method:

``` 
#<Model::PACCommitCollection:0x007f05540f3df8
 @commits=
  [#<Model::PACCommit:0x007f05540f3cb8
    @date=2015-04-27 10:37:05 +0000,
    @message="Test for multiple\n\nIssue: 1,2\n",
    @referenced=false,
    @sha="fb493078d9f42d79ea0e3a56abca7956a0d47123">,
   #<Model::PACCommit:0x007f05540f3b28
    @date=2015-04-27 10:37:05 +0000,
    @message="Test for empty\n",
    @referenced=false,
    @sha="55857d4e9838d1855b10e4c30b43a433e2db47cd">,
   #<Model::PACCommit:0x007f05540f3948
    @date=2015-04-27 10:37:05 +0000,
    @message="Test for none reference\n\nIssue: none\n",
    @referenced=false,
    @sha="a789b472150f462a8ae291577dcf7557b2b4ca55">,
   #<Model::PACCommit:0x007f05540f37b8
    @date=2015-04-27 10:37:05 +0000,
    @message="Updated readme file again - third commit\n\nIssue: 1\n",
    @referenced=false,
    @sha="cd32697cb7e2d3a7f3b77b5766ec22d31b002367">,
   #<Model::PACCommit:0x007f05540f3628
    @date=2015-04-27 10:37:05 +0000,
    @message=
     "Revert \"Updated readme file\"\n\nThis reverts commit 881b321e68481e0ae5cfab316b4b147e101f844a.\nIssue: 1\n",
    @referenced=false,
    @sha="a7b63f11d24b6f2fd164d35b904386b234667991">,
   #<Model::PACCommit:0x007f05540f3470
    @date=2015-04-27 10:37:05 +0000,
    @message="Updated readme file\n\nIssue: 3\n",
    @referenced=false,
    @sha="881b321e68481e0ae5cfab316b4b147e101f844a">,
   #<Model::PACCommit:0x007f05540f32b8
    @date=2015-04-27 10:37:05 +0000,
    @message="Initial commit - added README\n",
    @referenced=false,
    @sha="f9a66ca6d2e616b1012a1bdeb13f924c1bc9b4b6">]>
```

After this, this collection is passed to the `Core.task_id_list(...)` method and this is how this list look after we've matched commits to tasks each commits gets added to the task(s) it belongs to

``` 
#<Model::PACTaskCollection:0x007f7807f07a18
 @tasks=
  [#<Model::PACTask:0x007f7807f071a8
    @applies_to=#<Set: {"none"}>,
    @attributes={},
    @commit_collection=
     #<Model::PACCommitCollection:0x007f7807f07180
      @commits=
       [#<Model::PACCommit:0x007f7807f144c0
         @date=2015-04-27 10:37:05 +0000,
         @message="Test for multiple\n\nIssue: 1,2\n",
         @referenced=true,
         @sha="fb493078d9f42d79ea0e3a56abca7956a0d47123">,
        #<Model::PACCommit:0x007f7807f14010
         @date=2015-04-27 10:37:05 +0000,
         @message="Updated readme file again - third commit\n\nIssue: 1\n",
         @referenced=true,
         @sha="cd32697cb7e2d3a7f3b77b5766ec22d31b002367">,
        #<Model::PACCommit:0x007f7807f07e50
         @date=2015-04-27 10:37:05 +0000,
         @message=
          "Revert \"Updated readme file\"\n\nThis reverts commit 881b321e68481e0ae5cfab316b4b147e101f844a.\nIssue: 1\n",
         @referenced=true,
         @sha="a7b63f11d24b6f2fd164d35b904386b234667991">]>,
    @label=#<Set: {"none"}>,
    @task_id="1">,
   #<Model::PACTask:0x007f7807f06fa0
    @applies_to=#<Set: {"none"}>,
    @attributes={},
    @commit_collection=
     #<Model::PACCommitCollection:0x007f7807f06f78
      @commits=
       [#<Model::PACCommit:0x007f7807f144c0
         @date=2015-04-27 10:37:05 +0000,
         @message="Test for multiple\n\nIssue: 1,2\n",
         @referenced=true,
         @sha="fb493078d9f42d79ea0e3a56abca7956a0d47123">]>,
    @label=#<Set: {"none"}>,
    @task_id="2">,
   #<Model::PACTask:0x007f7807f05e48
    @applies_to=#<Set: {}>,
    @attributes={},
    @commit_collection=
     #<Model::PACCommitCollection:0x007f7807f05e20
      @commits=
       [#<Model::PACCommit:0x007f7807f14330
         @date=2015-04-27 10:37:05 +0000,
         @message="Test for empty\n",
         @referenced=false,
         @sha="55857d4e9838d1855b10e4c30b43a433e2db47cd">,
        #<Model::PACCommit:0x007f7807f07b30
         @date=2015-04-27 10:37:05 +0000,
         @message="Initial commit - added README\n",
         @referenced=false,
         @sha="f9a66ca6d2e616b1012a1bdeb13f924c1bc9b4b6">]>,
    @label=#<Set: {}>,
    @task_id=nil>,
   #<Model::PACTask:0x007f7807f050d8
    @applies_to=#<Set: {"none"}>,
    @attributes={},
    @commit_collection=
     #<Model::PACCommitCollection:0x007f7807f050b0
      @commits=
       [#<Model::PACCommit:0x007f7807f141a0
         @date=2015-04-27 10:37:05 +0000,
         @message="Test for none reference\n\nIssue: none\n",
         @referenced=true,
         @sha="a789b472150f462a8ae291577dcf7557b2b4ca55">]>,
    @label=#<Set: {"none"}>,
    @task_id="none">,
   #<Model::PACTask:0x007f7807efe990
    @applies_to=#<Set: {"none"}>,
    @attributes={},
    @commit_collection=
     #<Model::PACCommitCollection:0x007f7807efe968
      @commits=
       [#<Model::PACCommit:0x007f7807f07cc0
         @date=2015-04-27 10:37:05 +0000,
         @message="Updated readme file\n\nIssue: 3\n",
         @referenced=true,
         @sha="881b321e68481e0ae5cfab316b4b147e101f844a">]>,
    @label=#<Set: {"none"}>,
    @task_id="3">]>
```

There are another method, the one which applies task systems to this list, internally the only thing that changes is that the `attributes` field get's populated. How this looks depends on the task system, in the example above we just use the default system, which yields and empty hash. The next step simply takes this gross list and add the attributes that the task system contains.

```
  #Apply the task system(s) to each task. Basically populate each task with data from the task system(s)  
  Core.settings[:task_systems].each do |ts|
    everything_ok &= Core.apply_task_system(ts, tasks)
  end
```

A couple of things to note here:

 * `label` on the `PACTask` are the labels assigned via. regular expressions in the config file. It's a `Set` so values are unique, and will be overridden if you provide the same label to two different regular expressions
 * `referenced` on the `PACCommit` is the indicator that this commit has been referenced somewhere by a task.
 * `applies_to` on the `PACTask` indicates which task management system the task was found to belong two. It's a `Set` so value are unique.
 * Note that in all cases, `attributes` is an empty Hash. This variable will be populated by data for all the task systems this `PACTask` applies to. For Jira this is a straight up key/value hash.
 * `task_id` is the key that is used to look up data in the task management system(s) that populate the `attributes` field on the `PACTask`. Method of lookup varies and is unique for each task system. Trac uses an xmlrpc interface with a `ruby` gem, and `Jirà` is just plain `REST` using ruby standard `HTTP` libraries.

Task systems and the user defined regular expressions are applied sequentially to a task, in the order they are listed in the user supplied configuration file. A regex only has 1 label, but the regex can be copied to apply more than 1 label to the same task.

### Notes on Liquid

In order for liquid to produce locals for the template, we implemented a `to_liquid` method. This method should return a hash whose keys can be used 
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

Tests can be easily executed by running `rake test` for unit tests and `rake functional_test` for functional tests.

#### Integration tests for supported tasks systems

 * spin up a container with the task system
 * configure it
 * poor known test data into the system
 * run all relevant tests on the systems (assuming they do not change data)
 * shut down and clean-up containers

Possibly container, configuration and test data could be combined into a new container to save spin-up time. 

Note also that tests relies on know test data, e.g. task references etc. thus they are off-course part of this repository.

Choice of containers:

* **Jira**: The [`blacklabelops/jira`](https://hub.docker.com/r/blacklabelops/jira/) container seems very popular and easy to use, thus this is chosen as base container for Jira. . Compared to our own (Praqma) [staci project](https://github.com/Praqma/staci) that start a complete Atlassian suite, this simple jira container seems easier to start with.
* **Trac**: The [`jmmills/trac`](https://hub.docker.com/r/jmmills/trac/) container was chosen because it had the most pulls of the trac images on the hub, and was the only container that came with the xmlrpc for trac plugin as part of the package.   


The funcional test suites start the needed container with the `test/resources/start_task_system.sh` and the autogenerated `test/resources/stop_task_system.sh`


