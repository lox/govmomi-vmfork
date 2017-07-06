Govmomi VMFork
==============

VMFork is an experimental VMWare feature for forking VM's with copy-on-write disk and memory. This exposes those APIs to golang.

See http://www.yellow-bricks.com/2014/10/07/project-fargo-aka-vmfork-what-is-it/


Status
------

So far, the APIs work. I'm stuck on getting the quiesce process working on macOS though.  


Getting Started
---------------

```
go get github.com/lox/govmomi-vmfork/cli/vmfork
vmfork --help
```
