#!/usr/bin/env bash

MYSQL_PASS="7204f5e2-8b4f-11e8-8065-48bf6bd8fe98"
MYSQL_LOCAL_DIR=$GOPATH/src/github.com/mfioravanti2/entropy-api/test/mysql/data

# https://docs.docker.com/samples/library/mysql/
docker run -p 3306:3306 --detach --name=entropy-mysql -v ${MYSQL_LOCAL_DIR}:/var/lib/mysql --env="MYSQL_ROOT_PASSWORD=$MYSQL_PASS" mysql

# mysql -u root -p -h 127.0.0.1 -P 3306
