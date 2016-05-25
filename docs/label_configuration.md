# Labels

Labels are useful to group collected tasks, assign traits to them and and categorize them.

Labeled tasks are accessible in the template by using `tasks.<label>`, but the PAC configuration file need to assign the labels based on the regexp matching.

An use-case could be one repository, where several patterns for task references (e.g. Microsoft Office programs) are grouped in different sections in the changelog.

## Configuration file assigning labels 


Relevant configuration snippet:

	:task_systems:
	  - 
	    :name: none	    
	    :regex:
	      - { pattern: '/EXCEL-(\d+)/i', label: excel }
	      - { pattern: '/WORD-(\d+)/i ', label: word }
	      - { pattern: '/POWERPOINT-(\d+)/i ', label: powerpoint }

Collected tasks are grouped with labels using the above three regexps.

## Template using labels

An example template using labels:

	# PAC Changelog
	## Word 
	{% for task in tasks.word %}
	- {{task.task_id}}
	{% endfor %}
	## Excel
	{% for task in tasks.excel %}
	- {{task.task_id}}
	{% endfor %}	
	## Powerpoint
	{% for task in tasks.powerpoint %}
	- {{task.task_id}}
	{% endfor %}

The layout of the changelog will be three section with task references for each of the three Microsoft Office programs.

