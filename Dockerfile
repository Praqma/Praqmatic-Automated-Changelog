FROM ruby:2.5.5

RUN apt-get update && \
    apt-get install -y cmake && \
    apt-get install -y locales && \
    rm -rf /var/lib/apt/lists/* 

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY Gemfile /usr/src/app/

COPY . /usr/src/app

#Ruby knows best how to install this particular version of PAC CLI
RUN rake install

VOLUME ["/data"]
WORKDIR /data

CMD ["pac"] 
