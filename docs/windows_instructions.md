# PAC on Windows

In all the examples below. Please exchange paths and templates as appropriate. These are just examples. 

## Using docker for windows

### Windows CMD

1. `docker run -v %cd%:/data --net=host praqma/pac pac from 96025f592 --settings=/data/settings_pac-jira.yml`

### Git bash

1. `docker run -v $(cygpath -w $(pwd)):/data --net=host praqma/pac bash -c "pac from 96025f592 --settings=/data/settings_pac-jira.yml"`


## Using wsl

1. https://nickjanetakis.com/blog/setting-up-docker-for-windows-and-wsl-to-work-flawlessly
2. docker run -v $(pwd | sed -e 's/mnt//'):/data --net=host praqma/pac pac from 96025f592 --settings=/data/settings_pac-jira.yml



