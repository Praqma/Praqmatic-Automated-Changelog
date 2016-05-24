# {{title}}
{% for task in tasks.real %}
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
- Total numerber of commits: {{pac_c_count}}
- Referenced commits: {{pac_c_referenced}}
- Health: {{pac_health}}