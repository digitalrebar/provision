#!/usr/bin/env bash
# RackN Copyright 2021
# Turn DRP server into a bridge gateway

# get DRP kvm-test base network for $1 (kvm-test by default if not defined)
# saves configurations so they survive reboots
PATH=$PATH:/usr/local/bin

echo "RackN Digital Rebar Gateway Create"
echo "=============================================================="

if ! which drpcli > /dev/null ; then
  echo "MISSING: drpcli required!"
  exit 1
else
  echo "  verified drpcli is installed"
fi

if ! which jq > /dev/null ; then
  echo "  adding jq via drpcli as soft link in current directory"
  ln -s $(which drpcli) jq >/dev/null
else
  echo "  verified jq is installed"
fi

if ! drpcli info get > /dev/null ; then
  echo "MISSING: set RS_ENDPOINT and RS_KEY!"
  exit 1
else
  drpid=$(drpcli info get | jq -r .id | sed 's/:/-/g')
  echo "  verified drp $drpid access and credentials"
fi

NW=$(drpcli subnets list | jq -r .[0].Name)
NET=${1:-$NW}
echo "using network $NET"
JSON=/tmp/$NET-network.json
drpcli subnets show $NET > $JSON
NETWORK=$(cat $JSON | jq -r '.Subnet' | sed 's/\(.*\)\.\(.*\)\.\(.*\)\..*$/\1.\2.\3.0/')
HOST=$(cat $JSON | jq -r ".Options | .[] | select(.Code==3) | .Value")
MASK=$(cat $JSON | jq -r ".Options | .[] | select(.Code==1) | .Value")

echo "
NETWORK  :: $NETWORK
HOST     :: $HOST
MASK     :: $MASK
"

systemctl start iptables
systemctl enable iptables

iptables -t nat -A POSTROUTING -s "$NETWORK/$MASK" ! -d "$NETWORK/$MASK" -j MASQUERADE
iptables -I FORWARD 1 -i $NET -j ACCEPT
iptables -I FORWARD 1 -o $NET -m state --state RELATED,ESTABLISHED -j ACCEPT

service iptables save

sysctl net.ipv4.ip_forward=1
echo "net.ipv4.ip_forward=1" > /etc/sysctl.d/50-ipv4_ip_forward.conf
exit 0