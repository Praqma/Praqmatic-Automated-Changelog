#!/bin/bash
git init
touch readme.md
echo "#Test" > readme.md
git add .
git commit -am"Initial commit on master"
git checkout -b "branch1"
echo "#Branch1 commit 1" >> readme.md
git add .
git commit -am"Branch1 commit 1"
git checkout master
git merge branch1 --no-ff --no-edit
echo "#Master commit 2" >> readme.md
git commit -am"Second commit on master"
git checkout -b Branch2
echo 'Branch 2 commit 1' >> readme.md
git commit -am"Issue 1"
git checkout master
git merge Branch2 --no-ff --no-edit
git checkout -b branch3
echo 'Branch 3 commit 1' >> readme.md
git commit -am"Branch 3 commit 1"
git checkout master
git checkout -b Branch4
echo 'Branch 4 commit 1' >> readme.md
git commit -am"Branch 4 commit 1"
git checkout master
git merge branch3 --no-ff --no-edit
git merge branch4 --no-ff --no-edit
git merge Branch4 --no-ff --no-edit
git add . 
git commit -am"Commited merge conflict"