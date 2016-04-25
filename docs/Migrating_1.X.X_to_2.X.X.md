##Migrating from 1.X.X to 2.X.X

The format of the settings file has changed in 2.X.X so if you're updating to the latest version you'll need to update it.

Here is an example of how it has changed for a simple project

Old format: 

	:general:
	  date_template: "%Y-%m-%d"
	  changelog_name: "Changelog"
	  changelog_formats:
	    - "html"
	    - "pdf"
	  changelog_css: 
	  verbose: false

	:none:
	  regex:	    
	    - '/(JENKINS-[0-9]+)/i'

	:vcs:
	  type: git
	  usr:
	  pwd:
	  repo_location: "."

New format:

	:general:
	  date_template: '%Y-%m-%d'
	  :strict: false
	 
	:templates:
	  - { location: templates/default.md, output: Changelog.md }
	  - { location: templates/default_html.html, pdf: true, output: Changelog.html }

	:task_systems:
	  - 
	    :name: none
	    :regex:
	      - { pattern: '/(JENKINS-[0-9]+)', label: none }
	    :delimiter: '/,|\s/'
	  
	:vcs:
	  :type: git
	  :repo_location: '.'

##Steps to migrating

This section covers what you need to change in order to migrate. Section by section

### :general: section

1. You need to remove the `changelog_name` and `changelog_formats` from the `general` section. These setting are now part of the `:templates:` section. 
2. `changelog_css:` and `verbose:` has been deprecated. Delete those.

### :templates:

The previous attributes `changelog_name`, `changelog_css` and `changelog_formats` are all deprecated. Instead this is fully customized and you need to have atleast one template specified to produce an output. 

So if you previously had specified the following 
	
  changelog_name: "Changelog"
  changelog_formats:
    - "html"
    - "pdf" 

In 1.0.0 This would produce three files `Changelog.md`, `Changelog.html` and `Changelog.pdf` in order to produce the same output files in 2.X.X you need to add this to the `:templates:` section

	:templates:
	  - { location: templates/default.md, output: Changelog.md }
	  - { location: templates/default_html.html, pdf: true, output: Changelog.html }

That covers the output part. For the input we now requrie you to actally write a liquid template to base your output off this is `location:` attribute on the items in the template list. This adds flexibility to the output as you could decide to write a csv file. As for the templates and how they look like, take a look at the examples included in the `templates` folder of this project. Templates are also explained with examples in the Readme file of this project.

### :task_systems: 

If you used the task systems `none` or `trac` you need to move these under the `:task_system:` section, use the list notation demonstrated above:

	:task_systems:
		- 
			:name: '[trac or none]'

Next. add the `:regex:` section under your task system. Copy your regex into the section like so and asssign a label, just use the name of your 

	:regex:
	  - { pattern: '/[your pattern]', label: [none or trac] }

**NOTE:** Previously it was required that you had a named capture group `<id>` in your regex. This is no longer the case. We use the first capture group now and we're greedy, so one regex can have multiple matches on the same commit.  

### :vcs: 

Not many changes, usr and pw has been removed as these are not used at all in code.



