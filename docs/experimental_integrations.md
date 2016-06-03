# Experimental Integrations

The current version of PAC is actually quite capable of integrating with other tools than those officially supported.

This document roughtly documents some of the attemps, whether succesful or WIP.

## GitHub - Confirmed
Write up will follow soon.

## TFS (Visual Studio online) - Confirmed
I got it working against VSO without any code changes. I used my user name and a "personal access token" for authentication.

### default_settings    
```
    :name: jira
    :query_string: "https://reducto.visualstudio.com/DefaultCollection/_apis/wit/workitems?api-version=1.0&ids=#{task_id}"
    :usr: jak@praqma.net
    :pw: <my access token>
    :debug: false
    :regex:
    - { pattern: '/#(\d+)/i', label: tfs }      
```

### Template
```
# PAC Changelog from TFS
{% for task in tasks.tfs %}
## {{task.task_id}} {{task.data.value[0].fields["System.Title"]}}
{% for commit in task.commits %}
- {{commit.shortsha}}: {{commit.header}}
{% endfor %}
{% endfor %}
## Unspecified
{% for commit in tasks.unreferenced %}
- {{commit.shortsha}}: {{commit.header}} 
{% endfor %}
```

This was the tricky part, as it took some time figuring out how to correctly address the System.Title field from Liquid.

### Git log used
```
fc48c07 2016-06-03 Jan Krag (HEAD -> master) #6 commit number six
e7a4cc9 2016-06-03 Jan Krag #5 fifth commit
f3a0197 2016-06-03 Jan Krag #4 fourth commit
0b6aa4c 2016-06-03 Jan Krag #3 third commit
c36180f 2016-06-03 Jan Krag #2 second commit
408bacf 2016-06-03 Jan Krag #1 first commit
a8c6a55 2016-06-03 Jan Krag initial commit
```
