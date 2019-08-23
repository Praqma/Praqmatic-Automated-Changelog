# Templates

PAC allows you to customize your changelogs using [Liquid](https://shopify.github.io/liquid/) templates.

Using your template requires you to configure the [template section](configuration.md#templates) in the PAC configuration file.


Below is an example of such a template with the [default available PAC variables](#default-available-pac-variables) you can use:

	# PAC Changelog
	{% for task in tasks.referenced %}
	## {{task.task_id}}
	{% for commit in task.commits %}
	- {{commit.shortsha}}: {{commit.header}}
	{% endfor %}
	{% endfor %}
	## Unspecified
	{% for commit in tasks.unreferenced %}
	- {{commit.shortsha}}: {{commit.header}} 
	{% endfor %}

	## Nones
	{% for task in tasks.none %}
	- {{task.task_id}}
	{% endfor %}

	## Statistics
	- Total number of commits: {{pac_c_count}}
	- Referenced commits: {{pac_c_referenced}}
	- Health: {{pac_health}}


The above example is a mix of Liquid and simple MarkDown.

* Anything inside `{{}}` is an [object in Liquid](https://shopify.github.io/liquid/basics/introduction/#objects) and replaced by content. E.g. one of the PAC variables like `title`.
* The `%` are [Liquid tags](https://shopify.github.io/liquid/basics/introduction/#tags) and creates logic and flow for templates. E.g. iterating over `tasks` in the list `tasks.referenced`. _Remember to end the tag `{% endfor %}`_ 

So the above template have the following layout:

* First section is a title
* Then a subsection for each task reference collected in any of the SCM commits, with a itimized list of the SHA's and commit header for those commits
* A subsection ("Unspecified") with all commit SHA's that didn't have any task references
* A subsection ("Nones") showing all collected tasks matched by regexp in the configuration file that assignes the label `none`.
* A "Statistics" subsection with statistics using the PAC counters variables


## Default available PAC variables

PAC always make the following variable available to use in a template.

* `tasks.unknown` Collection of tasks, where data could not be extracted from the configured task system.
* `tasks.referenced` Collection of tasks which have been matched with one the regexp from the configuration file.
* `tasks.label` Collection of tasks with the assigned label `label`. In the example above used in the "Nones" section. Read more about labels in the [Labels](label_configuration.md) section.
* `tasks.unreferenced` Collection of commits where the configured regexp(s) in the configuration file didn't match a task.

* `pac_c_count` Total number of commits PAC considered.
* `pac_c_referenced` Number of commits with task references.
* `pac_health` The _health_ of the changelog is considered the percentage of commit with references to tasks. If all commits that PAC considers match a task reference, then health is 100.

* `task.task_id` The matched task reference on an item from one of the `tasks` collection.
* `task.commits` Collection of commits. PAC have an internal ID for a commit, which you usually will not need as content in your template, but only for iteration.
* `commit.[header | shortsha]` A `commit` is an item from the `tasks.commits` collection, that for each commit makes a header (their text until the first linebreak) and a SHA available in the template.

You can find more examples on PAC templates in the `templates` folder.

## Extra PAC variables

If you have defined the `:properties:` section in your configuration file or are running PAC with the `--properties` switch. You will also have the following variables available in your templates:

- `properties.*` The `*` should be substituted with the name of your variable.

This can be useful if you want to inject and use environment variables in your template, or you have a common set of changelog templates that rely on some form of dymanic content based on your build environment.  


## Using debug to inspect raw task system data

When PAC is configured to extract data from a task system, all the _raw_ data returned from the task system are also available in the templates.

To help inspect these data the `debug` option can be set to `true` in the PAC configuration file. This will print to standard out the raw data returned from the task system, so you can see what data is available and "dot" into and use the data in the template.

Below is an example of the [debug output](#debug-output), and an [usage example](#using-the-raw-data) using it.

## Debug output

Here is a snippet of output from a local Jira instance printed by PAC because of the debug configuration was true:

	"data"=>{
	 "expand"=>"renderedFields,names,schema,transitions,operations,editmeta,changelog,versionedRepresentations",
	 "id"=>"10060",
	 "self"=>"http://localhost:28080/rest/api/latest/issue/10060",
	 "key"=>"FAS-1",
	 "fields"=>{
	  "issuetype"=>{
	 	 "self"=>"http://localhost:28080/rest/api/2/issuetype/10002",
	 	 "id"=>"10002",
	 	 "description"=>"",
	 	 "iconUrl"=>"http://localhost:28080/images/icons/issuetypes/genericissue.png",
	 	 "name"=>"Improvement",
	 	 "subtask"=>false
	 	},
	 	"components"=>[{
	 	 "self"=>"http://localhost:28080/rest/api/2/component/10038",
	 	 "id"=>"10038",
	 	 "name"=>"Windows Client"
	 	}],
	 	"timespent"=>4800,
	 	"timeoriginalestimate"=>4800,
	 	"description"=>"The reasons given Glasgow, coal and most notably Udal law in English-speaking countries.
	 	...
	 	"summary"=>"Thatcher's government since 1952). Scotland The Scots pine marten.",

## Using the raw data

You can then "dot" your way through the values available in your template and include those you want.
The raw data are availble on the in `attributes.data` on the the default PAC variable item `task`.

In the [debug output](#debug-output) above we could use: 

* `task.attributes.data.fields.description` is the JIRA issue description
* `task.attributes.data.fields.summary` is the JIRA issue summary

Pretty powerfull - right?

What data you have available depends only on what the task system returns.

Below is an example using some of the raw data through `attributes.data`. 

	# PAC Changelog
	{% for task in tasks.referenced %}
	## {{task.task_id}} {task.attributes.data.fields.description}}
	 
	{{task.attributes.data.fields.summary}}
	
	{% for commit in task.commits %}
	- {{commit.shortsha}}: {{commit.header}}
	{% endfor %}
	{% endfor %}

This template have the following layout:

* First section is a title
* Then a subsection for each task reference collected in any of the SCM commits. 
  * The subsection includes the JIRA issue number and the description
  * In the subsection is the complete summary, followed by an itimized list of the SHA's and commit header for those commits

