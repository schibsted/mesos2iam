#!/usr/bin/env bash

set -a
eval $(/usr/local/bin/metadatavars)
eval $(/usr/local/bin/tags $EC2_INSTANCE_ID)
eval $(/usr/local/bin/dynamodbdata mesos-clusters-config-$realm)
set +a

/opt/mesos2iam/sbin/mesos2iam -iptables -host-ip $(ip route get 1 | awk '{print $NF;exit}') -smaug-url http://${SMAUG_URL} > /var/log/mesos2iam.log
