# stackgo

Configuration management you'll actually use

`stackgo` takes a different approach than other configuration management
systems. When configuring a server, stackgo uses a predetermined build order
for all resources.

## Configuration Order

1. Users
2. Groups
3. Software Repositories (ppas, yum repos, apt repos)
4. Packages
5. Source Repositories (git, svn)
5. Templates
6. Services
7. Databases
