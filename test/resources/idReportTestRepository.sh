#!/bin/bash

# setting names and stuff
if [ -z "$1" ]; then
	VERSION=""
else
	VERSION="_$1"
fi
NAME=idReportTestRepository
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
git init

git config user.name "Praqma Support"
git config user.email "support@praqma.net"

touch README.md
echo "# README of repository $REPO_NAME" >> README.md
echo "" >> README.md
echo "This is a test repository for functional tests." >> README.md
git add README.md
git commit -m "Initial commit - added README"




echo "Second commit to README of repository $REPO_NAME" >> README.md
echo "" >> README.md
git add README.md
git commit -m "Updated readme file

Issue: 3
"

git revert --no-commit `git rev-parse HEAD`
echo "Issue: 1" >> .git/MERGE_MSG 
git commit --no-edit


echo "Third commit to README of repository $REPO_NAME" >> README.md
echo "" >> README.md
git add README.md
git commit -m "Updated readme file again - third commit

Issue: 1
"
echo "Fourth commit to README of repository $REPO_NAME" >> README.md
echo "" >> README.md
git add README.md
git commit -m "Test for none reference

Issue: none
"
echo "Fifth commit to README of repository $REPO_NAME" >> README.md
echo "" >> README.md
git add README.md
git commit -m "Test for empty"

echo "Sixth commit to README of repository $REPO_NAME" >> README.md
echo "" >> README.md
git add README.md
git commit -m "Test for multiple

Issue: 1,2
"
echo "The final repository looks like this:" >> $LOG
echo "-------------------------------------" >> $LOG
git log --graph --pretty=format:'%Cred%h%Creset -%C(yellow)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset' --abbrev-commit --date=relative >> $LOG
echo "" >> $LOG
echo "" >> $LOG
git log >> $LOG
echo "" >> $LOG
echo "" >> $LOG


# Post process

cd $WORK_DIR
zip -r $NAME$VERSION.zip $REPO_NAME
rm -rf $REPO_NAME

