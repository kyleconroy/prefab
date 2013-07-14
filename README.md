# prefab - a basic configuration manager

[![Build Status](https://travis-ci.org/stackmachine/prefab.png?branch=master)](https://travis-ci.org/stackmachine/prefab)

`prefab` takes a different approach than other configuration management
systems. When configuring a server, prefab uses a predetermined build order
for all resources.

## Configuration Order

1. Users
2. Groups
3. Software Repositories (ppas, yum repos, apt repos)
4. Packages
5. Source Repositories (git, svn)
5. Directories
5. Templates
6. Tarballs
6. Services
7. Databases
