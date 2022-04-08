FROM ruby:3.0.1

WORKDIR /usr/src/app

RUN adduser --disabled-password --gecos --quiet pyroscope
USER pyroscope

COPY --from=pyroscope/pyroscope:latest /usr/bin/pyroscope /usr/bin/pyroscope
COPY Gemfile ./Gemfile
COPY Gemfile.lock ./Gemfile.lock

COPY main.rb ./main.rb
RUN bundle install

ENV PYROSCOPE_APPLICATION_NAME=simple.ruby.app
ENV PYROSCOPE_SERVER_ADDRESS=http://pyroscope:4040/
ENV PYROSCOPE_LOG_LEVEL=debug

CMD ["ruby", "main.rb"]
