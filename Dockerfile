FROM ruby:2.5.5

RUN apt-get update && \
    apt-get install -y cmake && \
    rm -rf /var/lib/apt/lists/* 

RUN apt-get update && \
    apt-get install -y locales && \
    rm -rf /var/lib/apt/lists/*  

RUN bundle config --global frozen 1

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY Gemfile /usr/src/app/
COPY Gemfile.lock /usr/src/app/

#Ruby knows best how to install this particular version of PAC
COPY . /usr/src/app
RUN rake install

VOLUME ["/data"]
WORKDIR /data

CMD ["pac"] 
