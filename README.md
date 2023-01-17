# nextbus

This project aims to provide a terminal-based realtime dashboard for monitoring
buses and tram in the city of Bordeaux, France, from the TBM transport company.

/!\ This project is not an official TBM project /!\

## How to install

```bash
$ git clone https://github/drawbu/nextbus.git
$ cd nextbus
$ go build -o nextbus main.go
```

## How to use

```bash
$ ./nextbus
Usage: nextbus [transport] [line] [stop]
Options and arguments (and corresponding environment variables):
-h, --help : print help (this message) and exit (also print this message if no
             argument is provided)
transport  : type of transport (bus, car, tram...)
line       : line number
stop       : stop name, optional, will print all stop in the line if missing
```

## Examples

```bash
$ ./nextbus bus 10 Peixotto
Bus 10, Peixotto (Tram B), direction Beausoleil
- 3 minutes
- 14 minutes
- 21 minutes
Bus 10, Peixotto, direction Jardin Botanique Terminus
- 9 minutes
- 14 minutes
- 24 minutes
```


```bash
$ ./nextbus bus 10
Bus Lianes 10

GRADIGNAN  BEAUSOLEIL
- Jardin Botanique Terminus
- Stalingrad (Tram A)
- Pôle emploi Bastide
etc...
```
