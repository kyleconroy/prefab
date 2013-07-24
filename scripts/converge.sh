#!/bin/bash
set -x

pushd /vagrant/software

./../bin/prefab postgresql/manifest.json redis/manifest.json \
  python/manifest.json java6/manifest.json rabbitmq/manifest.json \
  nginx/manifest.json mongodb/manifest.json memcached/manifest.json \
  ruby/manifest.json bundle/manifest.json

# DEVELOPMENT
popd

pushd /vagrant
sh scripts/installgo.sh
popd
