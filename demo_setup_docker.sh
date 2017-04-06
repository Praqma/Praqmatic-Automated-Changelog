#!/bin/bash
rm -rf demorepo/
mkdir demorepo
unzip test/resources/idReportTestRepository.zip -d demorepo/
cat << EOF > demorepo/idReportTestRepository/default_settings.yml 
:general:
  :strict: true

:templates:
  - { location: /usr/src/app/templates/default_id_report.md, output: ids.md }
  - { location: /usr/src/app/templates/default.md, output: default.md }
  - { location: /usr/src/app/templates/default_html.html, pdf: true, output: default.html }

:task_systems:
  -
    :name: none
    :regex:
      - { pattern: '/.*Issue:\s*(?<id>[\d+|[,|\s]]+).*?\n/im', label: none }
      - { pattern: '/.*Issue:\s*?(none).*?\n/im', label: none}
    :delimiter: '/,|\s/'

:vcs:
  :type: git
  :usr:
  :pwd:
  :repo_location: '.'
  :release_regex: 'tags'
EOF
docker build -t praqma/pac:snapshot .
docker run --rm -v $(pwd)/demorepo/idReportTestRepository:/data praqma/pac:snapshot from f9a66ca6d2e6
