#!/bin/bash
docker build . -t pac
docker run -v $(pwd):/data pac:latest pac from HEAD~20 --settings=/data/settings/default_settings.yml