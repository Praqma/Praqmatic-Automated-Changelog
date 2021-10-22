FROM ruby:3.0.2 as build

RUN apt-get update && \
    apt-get install -y cmake && \
    apt-get install -y libxslt-dev && \
    apt-get install -y libxml2-dev && \
    rm -rf /var/lib/apt/lists/*

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY Gemfile /usr/src/app/
COPY Gemfile.lock /usr/src/app/

#Ruby knows best how to install this particular version of PAC
#This means that this dockerfile can build any version of PAC.
RUN bundle install --without=test_gems

COPY . /usr/src/app

VOLUME ["/repo"]
VOLUME ["/templates"]

FROM ruby:3.0.2-slim as executor

RUN apt-get update && \
    apt-get install -y git && \
    rm -rf /var/lib/apt/lists/*

COPY --from=build /usr/src/app /usr/src/app
COPY --from=build /usr/local/bundle/ /usr/local/bundle/
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN ln -s /usr/src/app/pac.rb /usr/bin/pac

WORKDIR /repo
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["pac"]