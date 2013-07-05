#!/bin/bash
pushd /vagrant
./stack manifests/postgresql.json manifests/redis.json manifests/python3.json
