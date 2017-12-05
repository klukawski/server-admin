# Service Management Panel: #
Let's say that a few people bought a dedicated server to set up a few services. 
Every one of them wants root and SSH access, which may introduce problems like occupied ports, 
or someone can break someone else's service. Because of virtualization having quite an overhead, 
in this case we're using LXC (Linux Containers). 
It gives us an ability to set up multiple containers with their own IPs and environments.
Those containers need to be administered, and we don't want to call someone with main container's root access just to
restart ours, so we need some kind of panel. Also, there are people who don't own a container, but have a certain 
service set up (i.e. nginx to host a single site). Those people must be able to perform some operations on their 
services, like reload configs.  

### Roles: ###
* Head Admin: Has root access to main container.
  * Admin: Can manage all containers on the server.
    * Container Owner: Can manage their container(s) and services.
    * Service Owner: Can manage their service(s).

### Custom objects: ###

* Machine: Physical machine with LXC set up on it.
* Container: Single LXC container.
* Service: Single service set up by user on their container or by an Admin for a service owner on a separate container.

For testing, instead of operating on a live LXC machine, we will set up a sample microservice using a middleware 
written in go and set up on heroku. It will be imitating [LXD's](https://linuxcontainers.org/lxd/rest-api/) REST API,
and taking care of request authorization with JWT. 

### Simulated endpoints: ###
* /token
* /container/\<id>
  * /container/\<id>/start
  * /container/\<id>/stop
  * /container/\<id>/restart
  * /container/\<id>/exec
  * /container/\<id>/state
  * /container/\<id>/logs
