# Configuration 

Below is an example of a full PAC configuration file which uses all features from PAC.
Each configuration part is explained below so you can pick and choose for your own project.

Configuration file is YAML, so the : (colons), - (dash) and indentation matters.

	:general:
	  :strict: false

	:properties:
		title: 'Changelog name'
		product: 'Awesome product'

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
	    :debug: false
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

	:vcs:
	  :type: git
	  :repo_location: '.'

## General

### strict

**`strict`** If set to true PAC returns a non-zero exit code when a referenced task cannot be looked up your task system. 

_Defaults to `false`_.   

## Properties (_optional section_)

This section specifies properties that you want to use in your template. You can specify any arbitrary number of properties in this section. In the example shown above, the following variables can be referenced in Liquid: 

- `{{properties.title}}`
- `{{properties.product}}`

These values can be overridden at runtime by adding the `--properties` option when running PAC. Run PAC with the `-h` switch for an explanation on how to set a correct value for the `--properties` option.

## Templates

One or more [template configurations](templates.md). Each will result in a changelog report.
Put one template configuration pr. line

**`- { location: <path_to_template_file>, output: <path_to_output_file> }`**

Location and output values can be either relative (to PAC) or absolute. 
Each item should point to a Liquid template crafted for your changelog report and the destination for output.

You can specify a boolean **`pdf: true`** only together with HTML templates and HTML output. We require valid html in order to render pdf documents since we use a library called `pdfkit` to convert html files to pdf. As such there is no way for us to control the way the pdf is rendered, so if you do not like the way pdfkit renders your html you're free to use a different tool to create your pdf, an example could be the built in pdf printer in modern browsers.  

_`pdf` defaults to false._
 
## Task Systems

One or more task system configurations. Note the - (dash) before each.

A task system configuration must specify:

* **`name`** (_required_) one of `trac`, `jira`, `none`. Selects task system to extract data for collected tasks in the SCM commits. The `none` is special as it do not extract data from any task system. You only have the collected task references from the SCM commit messages.
* **`regex`** (_section is required_) is a list of regular expressions used to find the tasks in the SCM commits. Each entry is in the form: `{ pattern: <pattern>, label: <label> }`:
 * **`pattern`** (_one regexp is required_) is the reg exp used for matching tasks
 * **`label`** (_required_) is used to group the results, and be used for selecting, grouping and iteration in the templates. See [How to use labels](label_configuration.md)
* **`debug`** (_optional_) can be set to true (`:debug: true`) to print out to standard out the raw data returned from the task system. Useful information when writing themplates, so you can see what raw data is available with `dot`s in the templates. See [Using debug to inspect raw task system data](templates#using-debug-to-inspect-raw-task-system-data)
* **`delimiter`** (_optional_) an regex used to split commits further after the first match. PAC 1.x didn't support greedy matching, in order to match e.g. `#1,#2,#3`, one would have to specify a regex as the split delimiter. We generally discourage the use of this flag. 

Help writing regexp using Ruby IRB see this litle howto: [Howto write regexp using IRB](howto_write_regexp_using_irb.md)

### JIRA specific configuration

For JIRA task system (`:name: 'jira'`) the following is _required_ configuration:

* **`query_string`** as the location of your JIRA instance. _Always_ have the `task_id` in the string as in the example. Is is replaced by PAC with the task id captured by the regular expressions.
* **`usr`** is the JIRA user ID. The user must have read permission to the issues.
* **`pw`** is the password of the above JIRA user. It needs to be plain text.

_There is usually no required configuration to do in your JIRA_.

### Trac specific configuration

For Trac task system (`:name: 'trac'`) the following is _required_ configuration:

* **`trac_url`** is the base Trac location
* **`trac_usr`** is the Trac user ID, and the user must have `XML_RPC` permissions
* **`trac_pwd`** is the password of the above Trac user. It needs to be plain text.

Aside from granting the correct permissions to the user configured for PAC, another common issue you might run into is the way the XMLRPC interacts with the `AccountManager` plugin. When the AccountManager plugin is installed, every request will look like anonymous access. How to you fix that is explained [here](https://trac-hacks.org/wiki/XmlRpcPlugin#ProblemswhenAccountManagerPluginisenabled).

## VCS

Used for configuring the VCS to use. You can chose either Git or Mercurial (hg).

* **`type`** The VCS to use, we currently support `git` or `hg`
* **`repo_location`** The location of the repository to use. _Defaults to `.`_ for current working directory as PAC assumes to be called from within the repository for which to create a changelog.