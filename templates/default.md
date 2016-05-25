# PAC Changelog
{% for task in tasks.referenced %}
## {{task.task_id}}: [{{task.attributes.data.html_url}}]({{task.attributes.data.title}})
State: _{{task.attributes.data.state}}_
Milestone: _{{task.attributes.data.milestone}}_

{% for commit in task.commits %}
- {{commit.shortsha}}: {{commit.header}}
{% endfor %}
{% endfor %}
## Unspecified
{% for commit in tasks.unreferenced %}
- {{commit.shortsha}}: {{commit.header}} 
{% endfor %}
