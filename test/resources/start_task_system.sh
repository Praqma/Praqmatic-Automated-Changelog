#!/bin/bash

# Script used to start the containers with a tast system to use for functional testing.
# Each task system is started once for a series of functional tests.
# - start and stop is handled by this script, called by the functional test suites
# - this start script writes the stop script

# For error reporting:
SCRIPTNAME=`basename "$0"`

USAGE="Usage: $SCRIPTNAME jira|trac [docker host port number to bind to]"

# Chose task system
if [ -z "$1" ]
	  then
			    echo "$SCRIPTNAME - ERROR: No task system supplied"
					echo $USAGE
				fi
TASK_SYSTEM=$1
# Chose port number for the task system to be available on on the host
HOSTPORT=${2:-28080} # if ARG 1 is nul or empty - use default

# Using environment variable BUILD_NUMBER (from build system)
# to ensure we get unique containers started each item and in principle
# could run several jobs in parallel
if [ -z "$BUILD_NUMBER" ]; then
	echo "BUILD_NUMBER env. var. set - using default 0000"
	BN=0000
else
	BN=$BUILD_NUMBER
fi 

# For each time we start container, we create a simple stop with stop command
# matching the unique container names.
STOP_SCRIPT=test/resources/stop_task_system-$TASK_SYSTEM-$BN

# Little function that write the stop script
# to stop the started uniquely named containers
# Arguments:
# 1: unique container name
write_stop_script ()
{
	NAME=$1
	echo "Writing stop script to stop this container again: use ./$STOP_SCRIPT"
 	echo "docker stop $NAME" > $STOP_SCRIPT
 	echo "docker rm $NAME" >> $STOP_SCRIPT
	chmod +x $STOP_SCRIPT
}

if [ "$TASK_SYSTEM" = "jira" ]
then
	DB=postgres-jira-$BN
	JIRA=jira-$BN
	# Using jira images maintained by the guy called blacklabelops - they seem to fit our purpose.
	docker run --name $DB -d -e 'DB_USER=jiradb' -e 'DB_PASS=jellyfish' -e 'DB_NAME=jiradb' sameersbn/postgresql:9.4-12
	docker run -d --name $JIRA -e "JIRA_DATABASE_URL=postgresql://jiradb@postgres/jiradb" -e "JIRA_DB_PASSWORD=jellyfish" --link $DB:postgres -p $HOSTPORT:8080 blacklabelops/jira:7.1.4 #Jira Software server 
	# FIXME - must wait for container to start
  write_stop_script $JIRA
  write_stop_script $DB
elif [ "$TASK_SYSTEM" = "trac" ]
then
	# Using home-build trac images, as we had to do minor adjustments from one of the public ones.
	echo "Building docker image for Trac"
	docker build -t praqma/pac_test_trac test/resources/trac-env	# build without tag or version - alway use the latest one
	echo "Starting docker image for Trac"
	CONNAME=pac_test_trac-$BN
	docker run -p $HOSTPORT:80 -d --name $CONNAME praqma/pac_test_trac
	echo "Docker image for Trac"
	TRAC_RETRIES=0
	TRAC_MAX_RETRIES=30
	while ! curl -s http://localhost:$HOSTPORT/trac
	do
		echo "$(date) Trac is still not up"
	  sleep 10	  
		TRAC_RETRIES=$((TRAC_RETRIES + 1))
		if [ "$TRAC_RETRIES" -gt "$TRAC_MAX_RETRIES" ]; then
			echo "Trac timed out"
			break
		fi	  
  done
  echo "Trac is up and running" 
 	write_stop_script $CONNAME
else 
	echo "$SCRIPTNAME - ERROR: unknown task system"
	echo $USAGE
fi 

