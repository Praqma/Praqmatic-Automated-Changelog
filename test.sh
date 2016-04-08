#!/bin/bash
docker build -t praqma/pac:snapshot .
#This should output the help for PAC if build succeeded.
docker run --rm -t praqma/pac:snapshot pac 

