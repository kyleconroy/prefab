#!/bin/bash
pushd /vagrant
./stack postgresql/manifest.json redis/manifest.json python/manifest.json

# DEVELOPMENT
sh installgo.sh
