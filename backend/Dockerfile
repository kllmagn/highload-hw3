FROM ubuntu:19.04
ENV DEBIAN_FRONTEND noninteractive

# Installing packages
USER root
RUN apt-get update && apt-get upgrade -y && apt-get install -y gnupg git curl postgresql-11 postgresql-contrib

# Installing golang
ENV GOVERSION 1.13.1
RUN curl -s -O https://dl.google.com/go/go$GOVERSION.linux-amd64.tar.gz
RUN tar -xzf go$GOVERSION.linux-amd64.tar.gz -C /usr/local
RUN chown -R root:root /usr/local/go
ENV GOPATH $HOME/work
ENV PATH $PATH:/usr/local/go/bin
ENV GOBIN $GOPATH/bin
RUN mkdir -p "$GOPATH/bin" "$GOPATH/src"
RUN GO11MODULE=on

ARG CACHE_DATE=2016-01-01
WORKDIR /tech-db-forum
COPY . .

# Starting PostgreSQL and creating a database
USER postgres
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER forum WITH SUPERUSER PASSWORD 'forum';" &&\
    createdb -O forum forum &&\
    psql forum -a -f database/init.sql &&\
    /etc/init.d/postgresql stop

USER root
RUN cat database/postgresql.conf >> /etc/postgresql/11/main/postgresql.conf
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]
EXPOSE 5432

USER root
RUN go get
RUN go build main.go
CMD service postgresql start && ./main
EXPOSE 5000
