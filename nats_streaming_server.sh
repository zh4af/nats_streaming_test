# start nats-streaming-server, specife conf file, http monitor port, log file
nats-streaming-server -sc /etc/nats.d/nats.conf -m 8841 --log /var/log/nats.log &