#!/bin/bash
pushd /vagrant

./stack software/postgresql/manifest.json software/redis/manifest.json \
	software/python/manifest.json

# DEVELOPMENT
sh scripts/installgo.sh
