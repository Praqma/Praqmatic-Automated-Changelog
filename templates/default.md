# PAC Changelog
{% for task in tasks.referenced %}
## {{task.task_id}}
{% for commit in task.commits %}
- {{commit.shortsha}}: {{commit.header}}
{% endfor %}
{% endfor %}
{% if tasks.unreferenced.size > 0 %}
## Unspecified
{% for commit in tasks.unreferenced %}
- {{commit.shortsha}}: {{commit.header}} 
{% endfor %}
{% endif %}
