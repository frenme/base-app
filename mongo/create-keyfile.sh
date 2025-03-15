#!/bin/bash
openssl rand -base64 756 > mongodb-keyfile
chmod 600 mongodb-keyfile
echo "Created"