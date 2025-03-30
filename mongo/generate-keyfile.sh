#!/bin/bash

set -e

if [ -f keyfile ]; then
  echo "Mongo keyfile already exists"
else
  openssl rand -base64 756 > keyfile
  chmod 400 keyfile
fi