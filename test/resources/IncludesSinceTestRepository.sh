#!/bin/bash

# setting names and stuff
if [ -z "$1" ]; then
	VERSION=""
else
	VERSION="_$1"
fi
NAME=IncludesSinceTestRepository
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
echo "# $REPO_NAME" >> README.md
echo "" >> README.md
echo "This is a test repository." >> README.md
git add README.md
git commit -m "TASK 34 - Added a README"
git tag -a 'v1.0' -m"First release"

echo "Warning: Never cross the streams!" >> README.md
echo "" >> README.md
git add README.md
git commit -m "ISSUE 24 - Updated README to inform users about a bug."

echo "Comes with anti-marshmellow support." >> README.md
echo "" >> README.md
git add README.md
git commit -m "TASK 37 - Updated README to inform users about a new feature."
git tag -a 'v2.0' -m"Second release"

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

