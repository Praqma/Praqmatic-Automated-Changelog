FROM ruby:2.4

RUN apt-get update && \
    apt-get install -y cmake && \
    apt-get install -y libxslt-dev && \
    apt-get install -y libxml2-dev && \
    rm -rf /var/lib/apt/lists/* 

RUN apt-get update && \
    apt-get install -y locales && \
    rm -rf /var/lib/apt/lists/*

#Set the locale
RUN sed -i -e 's/# en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8/' /etc/locale.gen && \
    echo 'LANG="en_US.UTF-8"'>/etc/default/locale && \
    dpkg-reconfigure --frontend=noninteractive locales && \
    update-locale LANG=en_US.UTF-8

ENV LANG en_US.UTF-8  
ENV LANGUAGE en_US:en  
ENV LC_ALL en_US.UTF-8  

#This istalls a patched version of wkhtmltopdf that allows it to run headless without configuration
RUN mkdir -p /var/lib/wkhtml

WORKDIR /var/lib/wkhtml
RUN wget https://github.com/wkhtmltopdf/wkhtmltopdf/releases/download/0.12.3/wkhtmltox-0.12.3_linux-generic-amd64.tar.xz && \
		tar -xvf wkhtmltox-0.12.3_linux-generic-amd64.tar.xz && \
		rm wkhtmltox-0.12.3_linux-generic-amd64.tar.xz && \ 
		ln -sf /var/lib/wkhtml/wkhtmltox/bin/wkhtmltopdf /usr/bin/wkhtmltopdf

RUN bundle config --global frozen 1

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY Gemfile /usr/src/app/
COPY Gemfile.lock /usr/src/app/

#Ruby knows best how to install this particular version of PAC
#This means that this dockerfile can build any version of PAC.
RUN bundle install --without=test_gems

COPY . /usr/src/app

RUN ln -s /usr/src/app/pac.rb /usr/bin/pac

VOLUME ["/data"]
WORKDIR /data

COPY entrypoint.sh /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

CMD ["pac"] 
