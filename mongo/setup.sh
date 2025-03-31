#!/bin/bash

mkdir -p /etc/mongo
cp /tmp/keyfile /etc/mongo/keyfile
chown mongodb:mongodb /etc/mongo/keyfile
chmod 400 /etc/mongo/keyfile

(

until mongosh --host localhost:27017 --eval "db.adminCommand('ping')" > /dev/null 2>&1; do
  echo "Waiting for Mongo..."
  sleep 2
done

mongosh --host localhost:27017 --eval "
rs.initiate({
  _id: 'rs0',
  members: [
    { _id: 0, host: 'mongo1:27017' },
    { _id: 1, host: 'mongo2:27017' },
    { _id: 2, host: 'mongo3:27017' }
  ]
})"

sleep 20

mongosh --host localhost:27017 --eval "
  db.getSiblingDB('admin').createUser({
    user: 'root',
    pwd: 'example',
    roles: [{ role: 'root', db: 'admin' }]
  })
"

sleep 10

mongosh --host localhost:27017 -u admin -p password --authenticationDatabase admin --eval "
  db.adminCommand({ setParameter: 1, logLevel: 1 });
"

echo "Mongo cluster created"

) &

mongod --replSet rs0 --keyFile /etc/mongo/keyfile --port 27017 --bind_ip_all
