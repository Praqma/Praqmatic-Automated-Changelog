#!/bin/bash
#Very simple demonstration script that starts off 
rm -rf demorepo/
mkdir demorepo
unzip test/resources/idReportTestRepository.zip -d demorepo/

#Write a series of commits that we can match. We want to show that data ends up in the 
cd demorepo/idReportTestRepository
echo "FAS-1" > fas1.txt
git add fas1.txt
git commit -m"Fixed FAS-1 Added fas1"

echo "FAS-2" > fas2.txt
git add fas2.txt
git commit -m"Fixed FAS-2 Added fas2"

echo "FAS-3" > fas3.txt
git add fas3.txt
git commit -m"Fixed FAS-3 Added fas3"

cd ..
cd ..

#Now we have this repository:

#* f818a48 Fixed FAS-3 Added fas3
#* 287be34 Fixed FAS-2 Added fas2
#* 690a6aa Fixed FAS-1 Added fas1
#* fb49307 Test for multiple
#* 55857d4 Test for empty
#* a789b47 Test for none reference
#* cd32697 Updated readme file again - third commit
#* a7b63f1 Revert "Updated readme file"
#* 881b321 Updated readme file
#* f9a66ca Initial commit - added README


cat << EOF > demorepo/idReportTestRepository/default_settings.yml 
:general:
  :strict: true

:templates:
  - { location: /data/jira_template.md, output: jira.md }

:task_systems:  
  -
    :name: jira
    :debug: true
    :regex:
      - { pattern: '/(FAS-\d+)/', label: jra }        
    :query_string: "http://localhost:28080/rest/api/latest/issue/#{task_id}"
    :usr: 'admin'
    :pw: 'admin'     

:vcs:
  :type: git
  :repo_location: '.'
EOF

#We write a very simple jira template inside the repository we're mounting and adding to our container.
cat << EOF > demorepo/idReportTestRepository/jira_template.md
# {{title}}
{% for task in tasks.jra %}
## {{task.task_id}}

### Summary 

{{task.attributes.data.fields.summary}}

### Description

{{task.attributes.data.fields.description}}

### Associated commits
{% for commit in task.commits %}
- {{commit.shortsha}}: {{commit.header}}
{% endfor %}{% endfor %}
EOF

#Start an instance of jira. We need this to get data from Jira. For now we just use the container we built for test
#This can be adapted or removed if you already have a running instance of Jira configured.
./test/resources/start_task_system.sh "jira"

#Run PAC. We do it by mounting our repository inside the container.
docker build -t praqma/pac:snapshot .
docker run --rm --net=host -v $(pwd)/demorepo/idReportTestRepository:/data praqma/pac:snapshot -s f9a66ca6d2e6 
#Stop it again, if needed for this we just use the one that came with the test
./test/resources/stop_task_system-jira-0000.sh "jira"