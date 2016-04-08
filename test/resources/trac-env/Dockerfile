FROM ubuntu:14.04
MAINTAINER = Jason M. Mills <jmmills@cpan.org>
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update
RUN apt-get install -y pwgen
RUN apt-get install -y git-core
RUN apt-get install -y trac
RUN apt-get install -y trac-git
RUN apt-get install -y trac-accountmanager
RUN apt-get install -y trac-customfieldadmin
RUN apt-get install -y trac-xmlrpc
RUN apt-get install -y trac-wikiprint
RUN apt-get install -y trac-wysiwyg
RUN apt-get install -y trac-mastertickets
RUN apt-get install -y trac-tags
RUN apt-get install -y trac-diavisview
RUN apt-get install -y trac-announcer
RUN apt-get install -y trac-graphviz
RUN apt-get install -y python-flup

ADD setup_trac_config.sh /.setup_trac_config.sh
ADD setup_trac.sh /.setup_trac.sh
ADD run.sh /run.sh
ADD trac_logo.png /var/www/trac_logo.png

ADD set_trac_user_password.py /usr/local/bin/
RUN chmod 755 /usr/local/bin/set_trac_user_password.py

EXPOSE 80
CMD ["/run.sh"]

#SHA: 6c751b130fe43e436d4a7515106a6f07c63b94a1
#REPO: https://github.com/jmmills/docker-trac/