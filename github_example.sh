#!/bin/bash
start="${1:-f34ad72}"
docker build . -t pac
docker run -v $(pwd):/data pac:latest pac from $start --settings=/data/settings/default_settings.yml