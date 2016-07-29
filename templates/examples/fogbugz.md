{% for task in tasks.referenced %}
## Case {{task.task_id}}

### {{task.attributes.data.cases[0].case[0].sTitle[0]}} ({{task.attributes.data.cases[0].case[0].sStatus[0]}})
{% for commit in task.commits %}
- {{commit.shortsha}}: {{commit.header}}
{% endfor %}
{% endfor %}
## Unspecified
{% for commit in tasks.unreferenced %}
- {{commit.shortsha}}: {{commit.header}} 
{% endfor %}
