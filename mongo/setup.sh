#!/bin/bash

until mongosh --host mongo1:27017 -u ${MONGO_USERNAME} -p ${MONGO_PASSWORD} --authenticationDatabase admin --eval "db.runCommand({ ping: 1 })" > /dev/null 2>&1; do
  echo "Waiting for Mongo..."
  sleep 5
done

echo "Mongo is available"
mongosh --host mongo1:27017 -u ${MONGO_USERNAME} -p ${MONGO_PASSWORD} --authenticationDatabase admin --eval "rs.initiate({
  _id: 'rs0',
  members: [
    { _id: 0, host: 'mongo1:27017' },
    { _id: 1, host: 'mongo2:27017' },
    { _id: 2, host: 'mongo3:27017' }
  ]
})"
echo "Mongo cluster created"
sleep 5