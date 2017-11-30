# Service Management Panel: #
Rozpatrzmy przypadek gdzie community wspólnie zakupiło serwer dedykowany aby postawić różne rzeczy. Każdy chce mieć dostęp do SSH i prawa roota, co może doprowadzić do problemów takich jak porty będące w użyciu lub ktoś może zepsuć coś komuś innemu. Z powodu dużego overhead'u wirtualizacji używamy wtedy LXC (Linux Containers). Dzięki temu na jednej maszynie każdy ma dostęp do własnego roota oraz swój adres IP. Kontenerami LXC trzeba zarządać, a restart takiego kontenera wymaga interwencji osoby która ma dostęp do roota maszyny. Panel powinien więc posiadać przycisk restartu kontenera dostępny dla jego ownera. Dodatkowo potrzebujemy możliwości dodawania kontenerów, ponieważ użytkownicy mogą dodawać swoje subkontenery. Chcemy też mieć możliwość zarządzania poszczególnymi service'ami ustawionymi na kontenerach. Mamy także osoby które nie potrzebują kontenera i płacą za posiadanie samego service'u.
Let's say that a few people bought a dedicated server to set up a few services. Every one of them wants root and SSH access, which may introduce problems like occupied ports, 

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
