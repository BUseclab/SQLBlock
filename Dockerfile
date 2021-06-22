FROM ubuntu:latest

MAINTAINER Rasoul <rasoulj@bu.edu>

ENV DEBIAN_FRONTEND noninteractive
ENV INITRD No
ENV LANG en_US.UTF-8
ENV GOVERSION 1.12.2
ENV GOROOT /opt/go
ENV GOPATH /root/.go

# install required 
RUN apt update 
#RUN apt-get install -y software-properties-common
#RUN add-apt-repository -y ppa:ondrej/php

RUN apt-get install -y  php-fpm php-cli php-curl php-gd php-imap php-intl php-json php-ldap php-mbstring php-mysql php-opcache php-zip php-pgsql php-xml php-common php-dev cmake gcc make

RUN apt install -y nginx

RUN apt install -y mysql-server-5.7

## copy nginx and php config
COPY /storage/config/site.conf /etc/nginx/sites-enabled/default
COPY /storage/config/php.ini /etc/php/7.2/fpm/php.ini

## copy webapp
COPY /storage/webapp/ /var/www/html/
RUN chown -R www-data:www-data /var/www/html

# stupid php-fpm
RUN mkdir -p /var/run/php

RUN usermod -d /var/lib/mysql/ mysql &\
    mkdir /var/run/mysqld & \
    chown -R mysql:mysql /var/lib/mysql /var/run/mysqld

ADD /storage/config/createdb /home/
ADD /storage/startservice /home

RUN chmod +x /home/createdb
RUN chmod +x /home/startservice
RUN /home/createdb



WORKDIR /home/
ENTRYPOINT ["./startservice"]

# share modified php dbi
COPY storage/mysql-server/ /home/mysql-server/
COPY storage/php/ /home/php/
COPY storage/boost_1_59_0 /home/boost_1_59_0/

# install modified php dbi
RUN cd /home/php/php-mysqli &&\
    phpize &&\
    ./configure &&\
    make install

RUN cd /home/php/pdo-mysql &&\
    phpize &&\
    ./configure &&\
    make install

# installed modified php dbi

# install boost
RUN apt install -y libbz2-dev libncurses-dev bison wget

#RUN wget http://sourceforge.net/projects/boost/files/boost/1.59.0/boost_1_59_0.tar.gz &&\
#    tar xzf boost_1_59_0.tar.gz &&\

RUN cd /home/boost_1_59_0 &&\
    ./b2 install &&\
    cd /home/ &&\
    rm -r /home/boost_1_59_0

# installing mysql-server and SQLBlock module
RUN cd /home/mysql-server &&\
    make &&\
    make install

# copy modified mysql-server and SQLBLock plugin
RUN cp /usr/local/mysql/bin/* /usr/bin
RUN cp /home/mysql-server/plugin/sqlblock/sqlblock.so /usr/lib/mysql/plugin


RUN apt install -y wget git make gcc python3

# install go 1.12
RUN cd /opt && wget https://golang.org/dl/go${GOVERSION}.linux-amd64.tar.gz &&\
    tar zxf go${GOVERSION}.linux-amd64.tar.gz && rm go${GOVERSION}.linux-amd64.tar.gz &&\
    ln -s /opt/go/bin/go /usr/bin/ &&\
    mkdir $GOPATH

#copy SQLBlock SA to $GOPATH directory
COPY /storage/static_analysis/parse-function $GOPATH/src/parse-function
COPY /storage/static_analysis/createProfile.py /home/

# copy sqlblock module scripts to tmp
COPY /storage/sqlblock/record /home/record_mysql_query
COPY /storage/sqlblock/enforce /home/enforce_profile

RUN chmod +x /home/record_mysql_query
RUN chmod +x /home/enforce_profile

#copy runme.sh
COPY /storage/runme.sh /home/runme
RUN chmod +x /home/runme
