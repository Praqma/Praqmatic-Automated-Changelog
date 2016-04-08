#!/bin/bash

source /.setup_trac_config.sh

setup_trac() {
    [ ! -d /trac ] && mkdir /trac
    if [ ! -f /trac/VERSION ]
    then
        trac-admin /trac initenv "My New Project" sqlite:db/trac.db git /repo.git
        setup_components
        setup_accountmanager
        setup_admin_user
        trac-admin /trac config set logging log_type stderr
        [ -f /var/www/trac_logo.png ] && cp -v /var/www/trac_logo.png /trac/htdocs/your_project_logo.png
    fi
}

setup_repo() {
    if [ ! -d /repo.git ]
    then 
        git config --global user.name "Trac Admin"
        git config --global user.email trac@localhost

        mkdir /repo.git
        pushd /repo.git
            git init --bare
        popd

        pushd /tmp
            git clone --no-hardlinks /repo.git repo
            pushd repo
                echo repository init >README
                git add README
                git commit README -m "initial commit"
                git push origin master
            popd
            rm -rf repo
        popd
    fi
}

clean_house() {
    if [ -d /.setup_trac.sh ] && [ -d /.setup_trac_config.sh ]
    then
        rm -v /.setup_trac.sh
        rm -v /.setup_trac_config.sh
    fi
}

setup_repo
setup_trac
clean_house
