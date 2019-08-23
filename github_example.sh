#!/bin/bash
docker build . -t pac
docker run -v $(pwd):/data pac:latest pac from f34ad72 --settings=/data/settings/default_settings.yml