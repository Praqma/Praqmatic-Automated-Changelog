setup_components() {
    trac-admin /trac config set components tracext.git.* enabled 
    trac-admin /trac config set components trac_gitolite.* enabled 
    trac-admin /trac config set components trac.web.auth.LoginModule disabled 
    trac-admin /trac config set components acct_mgr.web_ui.* enabled 
    trac-admin /trac config set components acct_mgr.admin.* enabled 
    trac-admin /trac config set components acct_mgr.register.* enabled 
    trac-admin /trac config set components acct_mgr.notification.* enabled 
    trac-admin /trac config set components acct_mgr.htfile.* disabled 
    trac-admin /trac config set components acct_mgr.db.sessionstore enabled 
    trac-admin /trac config set components acct_mgr.pwhash.htdigesthashmethod enabled 
    trac-admin /trac config set components acct_mgr.pwhash.htpasswdhashmethod disabled 
    trac-admin /trac config set components acct_mgr.svnserve.svnservepasswordstore disabled
    trac-admin /trac config set components tracrpc.* enabled
}

setup_accountmanager() {
    trac-admin /trac config set account-manager hash_method HtDigestHashMethod 
    trac-admin /trac config set account-manager db_htdigest_realm TracDB 
    trac-admin /trac config set account-manager password_store SessionStore 
    trac-admin /trac config set account-manager reset_password false
}

setup_admin_user() {
    trac-admin /trac session add admin admin root@localhost
    trac-admin /trac permission add admin TRAC_ADMIN
    trac-admin /trac permission add anonymous XML_RPC
    trac-admin /trac permission add anonymous TICKET_CREATE

    if [ ! -f /.trac_admin_password ]
    then
        TRAC_PASS=UqmvS76r7D
        set_trac_user_password.py /trac admin ${TRAC_PASS}
        echo $_S
        echo Trac admin login
            echo -e "\tUser: admin"
            echo -e "\tPassword: ${TRAC_PASS}"
        echo $_S
        echo ${TRAC_PASS} >/.trac_admin_password
    fi

}
