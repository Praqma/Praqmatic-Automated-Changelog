#!/bin/bash

# setting names and stuff
if [ -z "$1" ]; then
	VERSION=""
else
	VERSION="_$1"
fi
NAME=GetCommitMessagesTestRepository
REPO_NAME=$NAME$VERSION # used for manual testing of script and re-runs
WORK_DIR=`pwd`

LOG=$WORK_DIR/$REPO_NAME-repo_description.log
echo "# Repository view and commits" >> $LOG
echo "" >> $LOG
echo "Git version:" >> $LOG
git --version >> $LOG
echo "" >> $LOG


mkdir -v $REPO_NAME
cd $REPO_NAME
git init # not using bare, pac expect checked out repositories

git config user.name "Praqma Support"
git config user.email "support@praqma.net"

# except from the initial commit to readme, we use commit log files
# to avoid merge conflicts (not the purpose to get those..)

touch README.md
echo "# README of repository $REPO_NAME" >> README.md
echo "" >> README.md
echo "This is a test repository for functional tests." >> README.md
git add README.md
git commit -m "Initial commit - added README"

MCL="master-commit.log"
date >> $MCL
git add $MCL
git commit -m "Commit on master

Issue: 3
"
echo "The test repository with commits:" >> $LOG
echo "---------------------------------" >> $LOG
echo "git log --graph --source --all" >> $LOG
git log --graph --source --all >> $LOG
echo "" >> $LOG
echo "" >> $LOG
echo "---------------------------------" >> $LOG
echo "git branch -a" >> $LOG
git branch -a >> $LOG
echo "" >> $LOG
echo "" >> $LOG

# Now create parallel development on several dev branches and mimics the features of the pretested integration plugin and the automated git flow:
# https://wiki.jenkins-ci.org/display/JENKINS/Pretested+Integration+Plugin
# http://www.josra.org/blog/An-automated-git-branching-strategy.html

git checkout -b dev1
d1CL="dev1-commit.log"
date >> $d1CL
git add $d1CL
git commit -m "Commit on dev1 branch

Issue: 100
"

git checkout master
git checkout -b dev2
d2CL="dev2-commit.log"
date >> $d2CL
git add $d2CL
git commit -m "Commit on dev2 branch

Issue: 200
"

git checkout master
date >> $MCL
git commit -am "Commit on master again

Issue: 4
"

git checkout -b dev3
d3CL="dev3-commit.log"
date >> $d3CL
git add $d3CL
git commit -m "Commit on dev3 branch

Issue: 300
"

git checkout master
date >> $MCL
git commit -am "Commit on master again again

Issue: 5
"

git checkout dev3
date >> $d3CL
git commit -am "Commit on dev3 branch again

Issue: 301
"

# deliver branch dev3 in a accumulated commit ish fashion
git checkout master
git merge -m "Accumulated commit of branch 'origin/dev3':

To master in pretested integration plugin fashion,
deleting branch dev3.

This message done manually, but issue references added.

    Issue: 300


    Issue: 301

" dev3 --no-ff
git branch -D dev3

# making a few more commit on dev1 and dev2 branch, including update from master on dev1

git checkout dev1
git merge master --no-edit

date >> $d1CL
git commit -am "Another commit on dev1, after merging master in to be updated"
date >> $d1CL
git commit -am "Last commit on dev1 before deliver"

# deliver branch dev3 in a accumulated commit ish fashion
git checkout master
git merge -m "Accumulated commit of branch 'origin/dev1':

To master in pretested integration plugin fashion,
deleting branch dev1.

This message done manually, but issue references added.

    Issue: 100

" dev1 --no-ff
git branch -D dev1


# commit on dev2 now

git checkout dev2
date >> $d2CL
git commit -am "Second commit on dev2"


# final commit will be on dev4, to check if head is not master in the test we write
git checkout master
git checkout -b dev4
d4CL="dev4-commit.log"
date >> $d4CL
git add $d4CL
git commit -m "Commit on dev4 branch

Issue: 400
"

# leaving branch master checked out

git checkout master



echo "The test repository with commits:" >> $LOG
echo "---------------------------------" >> $LOG
echo "git log --graph --source --all" >> $LOG
git log --graph --source --all >> $LOG
echo "" >> $LOG
echo "" >> $LOG
echo "---------------------------------" >> $LOG
echo "git branch -a" >> $LOG
git branch -a >> $LOG
echo "" >> $LOG
echo "" >> $LOG
echo "---------------------------------" >> $LOG
echo "git show-branch --all" >> $LOG
git show-branch --all >> $LOG
echo "" >> $LOG
echo "" >> $LOG
echo "---------------------------------" >> $LOG
echo "git show-branch --all --sha1-name" >> $LOG
git show-branch --all --sha1-name >> $LOG
echo "" >> $LOG
echo "" >> $LOG

# Post process

cd $WORK_DIR
zip -r $NAME$VERSION.zip $REPO_NAME
rm -rf $REPO_NAME

