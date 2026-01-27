# PAC Changelog

**Generated:** {{ 'now' | date: '%Y-%m-%d' }}

## Referenced Tasks

{% for task in tasks.referenced %}
### {{ task.task_id }} - {{ task.attributes.data.title }}

{% for commit in task.commits %}
- {{ commit.header }}
{% endfor %}

{% endfor %}

## Unreferenced Commits

{% for commit in tasks.unreferenced %}
- {{ commit.header }}
{% endfor %}
