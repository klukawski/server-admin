# Service Management Panel: #
Let's say that a few people bought a dedicated server to set up a few services. Every one of them wants root and SSH access, which may introduce problems like occupied ports, or someone can break someone else's service. Because of virtualization having quite an overhead, in this case we're using LXC (Linux Containers). It gives us an ability to set up multiple containers with their own IPs and environments. 

### Role: ###
* Head Admin: Posiada bezpośredni dostęp do maszyny.
  * Admin: Może zarządzać wszystkimi kontenerami w razie problemów.
    * Container Owner: Osoba posiadająca kontener.
    * Service Owner: Osoba posiadająca wyłącznie service.

### Custom object'y: ###

* Machine: Fizyczna maszyna na której stoi LXC.
* Container: Pojedynczy kontener LXC stojący na maszynie.
* Service: Pojedynczy Service stworzony przez użytkownika na jego kontenerze lub stworzony przez Admina na specjalnym kontenerze dla osób posiadających tylko Service.

Dla testów zamiast operować na maszynie z LXC ustawimy tylko przykładowy microservice za pomocą middleware'u napisanego w go i ustawionego na heroku. Będzie on końcowo nakładką na REST API udostępniane przez [LXD](https://linuxcontainers.org/lxd/rest-api/). Do zabezpieczenia API użyte zostanie JWT na zasadzie współdzielonego pre-share'owanego klucza i jednorazowych tokenów operacji.

### Symulowane endpointy: ###
* /token
* /container/\<id>
  * /container/\<id>/start
  * /container/\<id>/stop
  * /container/\<id>/restart
  * /container/\<id>/exec
  * /container/\<id>/state
  * /container/\<id>/logs
