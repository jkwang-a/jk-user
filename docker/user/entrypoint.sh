#!/bin/sh

if [[ -z $DB_HOST ]]; then DB_HOST=mysql; fi
if [[ -z $DB_PORT ]]; then DB_PORT=3306; fi
if [[ -z $DB_NAME ]]; then DB_NAME=bmc_user; fi
if [[ -z $DB_USER ]]; then DB_USER=root; fi
if [[ -z $DB_PASS ]]; then DB_PASS=root; fi

sed -i -e 's/DB_HOST/'$DB_HOST'/g' \
       -e 's/DB_PORT/'$DB_PORT'/g' \
       -e 's/DB_NAME/'$DB_NAME'/g' \
       -e 's/DB_USER/'$DB_USER'/g' \
       -e 's/DB_PASS/'$DB_PASS'/g' \
     /opt/user/conf/app.conf


if [[ -z $CSE_REGISTRY_ADDRESS ]] ; then CSE_REGISTRY_ADDRESS="service-center:31100"; fi
if [[ -z $CSE_GRPC_LISTEN ]] ;      then CSE_GRPC_LISTEN="0.0.0.0:5001"; fi
if [[ -z $CSE_GRPC_ADVERTISE ]] ;   then CSE_GRPC_ADVERTISE="0.0.0.0:5001"; fi

sed -i -e 's/CSE_REGISTRY_ADDRESS/'$CSE_REGISTRY_ADDRESS'/g' \
       -e 's/CSE_GRPC_LISTEN/'$CSE_GRPC_LISTEN'/g' \
       -e 's/CSE_GRPC_ADVERTISE/'$CSE_GRPC_ADVERTISE'/g' \
     /opt/user/conf/chassis.yaml


[ -d /logs ] || mkdir -p /logs
cd /opt/user; ./bmc-user | rotatelogs /logs/bmc-user.%Y-%m-%d.log 86400 2>&1

