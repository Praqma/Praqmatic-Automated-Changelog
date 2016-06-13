#!/bin/bash

# Script used to start the containers with a tast system to use for functional testing.
# Each task system is started once for a series of functional tests.
# - start and stop is handled by this script, called by the functional test suites
# - this start script writes the stop script

# For error reporting:
SCRIPTNAME=`basename "$0"`
# To make sure stop script is placed correctly relative to this script
SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

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
	echo "BUILD_NUMBER env. var. is NOT set - using default 0000"
	BN=0000
else
	BN=$BUILD_NUMBER
fi 

# For each time we start container, we create a simple stop with stop command
# matching the unique container names.
STOP_SCRIPT=$SCRIPTDIR/stop_task_system-$TASK_SYSTEM-$BN.sh
# remove old script
rm -vf $STOP_SCRIPT

# Little function that write the stop script
# to stop the started uniquely named containers
# Arguments:
# 1: unique container name
write_stop_script ()
{
	NAME=$1
	echo "Writing stop script to stop this container again: use ./$STOP_SCRIPT"
 	echo "docker stop $NAME" >> $STOP_SCRIPT # always append to script as it is used for two containers and is deleted above initially
 	echo "docker rm $NAME" >> $STOP_SCRIPT
	chmod +x $STOP_SCRIPT
	echo "Current stop script contains:"
	cat $STOP_SCRIPT
}

if [ "$TASK_SYSTEM" = "jira" ]
then
	
	# Binary created and tagged images as a temporary solution to avoid the problem with Jira not supporting fully un-attended install, initial setup and configuration.
	DB=pac_test_postgres-for-jira-$BN
	JIRA=pac_test_jira-$BN
	# Using jira images maintained by the guy called blacklabelops - they seem to fit our purpose.
	docker run --name $DB -d -e 'DB_USER=jiradb' -e 'DB_PASS=jellyfish' -e 'DB_NAME=jiradb' praqma/pac_test_postgres-for-jira:v2
	docker run -d --name $JIRA -e "JIRA_DATABASE_URL=postgresql://jiradb@postgres/jiradb" -e "JIRA_DB_PASSWORD=jellyfish" --link $DB:postgres -p $HOSTPORT:8080 praqma/pac_test_jira:v2
	docker ps -a | grep pac_test
  write_stop_script $JIRA
  write_stop_script $DB
	# while jira is starting up the following url is shown: http://localhost:28080/startup.jsp?returnTo=%2Fsecure%2FDashboard.jspa
	# so checking if get the real url back as url_effective
	URL_RETRIES=0
	URL_MAX_RETRIES=30
	URL="http://localhost:$HOSTPORT/secure/Dashboard.jspa"
	sleep 15 # takes a least 15 secs
	URL_RETRIES=0
	URL_MAX_RETRIES=30
	URL="http://localhost:$HOSTPORT/secure/Dashboard.jspa"
  echo "Checking is Jira is up and running"	
	while ! curl -Ls -o /dev/null -w %{url_effective} $URL | grep $URL
	do
		echo "$(date) URL is still not up"
		sleep 10    
		URL_RETRIES=$((URL_RETRIES + 1))
		if [ "$URL_RETRIES" -gt "$URL_MAX_RETRIES" ]; then
			echo "URL timed out"
			break
		fi    
	done
	echo "Jira is ready"

elif [ "$TASK_SYSTEM" = "trac" ]
then
	# Using home-build trac images, as we had to do minor adjustments from one of the public ones.
	echo "Building docker image for Trac"
	docker build -t praqma/pac_test_trac $SCRIPT_DIR/trac-env	# build without tag or version - alway use the latest one
	echo "Starting docker image for Trac"
	CONNAME=pac_test_trac-$BN
	docker run -p $HOSTPORT:80 -d --name $CONNAME praqma/pac_test_trac
	docker ps -a | grep pac_test
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
