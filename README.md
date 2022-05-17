# SilphTelescope

![Silph Telescope logo, a friendly spider who's waving at you](doc/logo.png)

This is a Matrix bot that users can set up to watch for events in Pokemon Go.

Imagine you're playing Pokemon Go and want to be notified whenever an interesting Pokemon (Unown, Axew, Noibat, ...) spawns in your area so you can go catch it. Or you want to be notified when a Raid you're looking for pops up near you.

## Admins

### Requirements

* a working [Map-A-Droid](https://github.com/Map-A-Droid/MAD) setup with workers who scrape Pokemon Go events
* a [Matrix](https://matrix.org/) server (doesn't need to be your own as you don't need special privileges)
* Docker

### Initialize Pokedex

To associate Pokemon names with their numbers we need a Pokedex.

By default `pokedex.json` resides in the `data` directory. It isn't included here, so you need to generate it from current PokeAPI data using `pokedexgen` from inside the container:

```console
% docker-compose exec app sh
/app # ./pokedexgen
```

This creates a new file named `/app/pokedex.json`. To use it in the bot copy it out of the container and move it to `./data/pokedex.json`. Then rebuild the docker image and restart your container.

### Initialize GeoDex

MAD builds a database of Forts (i.e. Gyms and Pokestops) its workers see. The internal representation is a mapping of a GUID to a location.
You can also get names for Pokestops when using the Quest scanner feature of MAD. Or you can scrape them from Ingress Intel. Or you can scrape them from 3rd party services.

To fill the fort databases with data from MAD, run `geodexgen` inside the container.

You need to supply:

* GeoDex location: `/data/geodex` if you're using the included `docker-compose.yaml`
* SQL hostname: address of MAD's MySQL server. Your container must be able to reach it, of course (e.g. be in MAD's network).
* Tile38 hostname: host and port of Tile38 server. `tile38:9851` if you're using the included `docker-compose.yaml`
* Zero or more `-b <boq.json>` flags: BookOfQuests `stops` data to import gym names from. See wiki.

```console
% docker-compose exec app sh
/app # ./geodexgen --geodex /data/geodex --sql-hostname mariadb --t-hostname tile38:9851 -b ./data/boq_a.json -b ./data/boq_b.json 
INFO[0000] Connected to sqldb mariadb running 10.3.27-MariaDB-1:10.3.27+maria~focal 
INFO[0000] > setup took 10.270889ms                     
INFO[0016] Pokestops read from MAD: 7926                
INFO[0016] > mad pokestop import took 16.773985772s     
INFO[0021] Gyms read from MAD: 2003                     
INFO[0021] > mad gym import took 4.2817277s             
INFO[0036] processed BOQ data: 513 cells containing 110512 POIs with 11461 gyms 
INFO[0036] added names to 1343 gyms, got 16 gyms which already had a name 
INFO[0036] > boq import took 15.300342142s              
INFO[0036] Fort nearest to (52.5395,13.4161): GUID=2342cafef00d0101010101010101010.16 Type=Stop (52.5395365,13.4161123) Name: Relief
INFO[0036] Fort nearest to (52.5399,13.4208): GUID=42cafef00d010101010101010101023.16 Type=Gym (52.5399245,13.4208453) Name: Women Graffiti
INFO[0036] > example lookup took 4.748235ms
```

## Developers

### Set up pre-commit Git hook

Install dependencies:

* Go, of course
* `go install golang.org/x/lint/golint@latest`

Use `./run_tests.sh` before committing or use the included hook to do it automatically:

```shell
ln -s ../../pre-commit.sh .git/hooks/pre-commit
```

For all tests to run you also need the pokedex, see section *Initialize Pokedex* above.

## Get Access Token

The bot needs an access token from the Matrix server.

### Register

Register a new user for the bot:

```console
./multitool.py register --url https://matrix.example.com silphtelescope01 secret
2021-02-08 04:42:25.654 | INFO     | __main__:register:41 - registered user silphtelescope01 successfully
response: {
  "user_id": "@silphtelescope01:matrix.example.com",
  "home_server": "matrix.example.com",
  "access_token": "redacted",
  "device_id": "redacted"
}
```

The response already contains a token. No need to login.

### Login

Login with an existing user to get the token:

```console
./multitool.py login --url https://matrix.example.com silphtelescope01 secret
2021-02-08 04:42:57.221 | INFO     | __main__:login:54 - login successful
response: {
  "user_id": "@silphtelescope01:matrix.example.com",
  "access_token": "redacted"
  "home_server": "matrix.example.com",
  "device_id": "redacted"
}
```
