# SleepSpotify v0.1.2
This project allows you to pause spotify at a specific UTS

## What do you need
- GoLang (^1.9)
- MariaDB/MySQL (^10.1.23/^5.7.19)

## How to install
### MySQL Part

Create a database and a user on your SQL server. Then, use this script to initialize the database

```sql
CREATE TABLE `pause` (
  `ID` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `token` text NOT NULL,
  `Uts` int(10) unsigned NOT NULL,
  `refresh` text NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8
```

### Configuration
To configure SleepSpotify you will need to [register an application on the Spotify Website](https://developer.spotify.com/my-applications/#!/applications)

In the `config/config.yaml` file you will find the following config
- `Spotify`
  - `ClientID` => the client ID found when you register the application
  - `ClientSecret` => the Client Secret found when you register the application
- `SessionSecret` => something secret to protect the data stored in the cookie
- `DomainName` => the Domain Name where you can found the SleepSpotify API (e.g.: localhost:8000)
- `DB`
  - `Host` => the domain name or IP address where you can found your database (caution the communication between SleepSpotify and your DB is unsecured)
  - `Port` => the port where the SQL server listens
  - `Name` => the database name
  - `Username` => the username to access the database
  - `Password` => the password to access the database
- `Angular` => the URL to access to the Angular front-end (e.g.: http://localhost:4200 for dev)

### Golang Part
```bash
go get -u github.com/SleepSpotify/SleepSpotify
cd $GOPATH/src/github.com/SleepSpotify/SleepSpotify
cp config/config.yaml.sample config/config.yaml
```

edit the `config/config.yaml` file

```bash
go build
./SleepSpotify
```

## API
this project is an API:
- Account
  - GET `/login` the call to redirect to the Spotify API to connect with the Spotify OAuth2 service
  - GET `/callback` the call that the Spotify OAuth2 service will call after a successful authentication
  - GET `/account/isConnected` to check if you are connected.
    - response
      - JSONConnected
      - JSONError
- Player
  - PUT `/spotify/pause` to pause the music on Spotify
    - response
      - JSONActionDone
      - JSONError
  - GET `/spotify/sleep` to check the status of the timer
    - response
      - JSONSleepFound
      - JSONError
  - POST `/spotify/sleep` to set a timer
    -  body params
      - uts => an uts to stop the timer
    - response
      - JSONSleepFound
      - JSONError
  - PUT `/spotify/sleep` to update a timer.
    -  body params
      - uts => an uts to stop the timer
    - response
      - JSONSleepFound
      - JSONError
  - DELETE `/spotify/sleep` to delete a timer.
    - response
      - JSONActionDone
      - JSONError

All the response object are [defined here](https://github.com/SleepSpotify/SleepSpotify/blob/master/controler/json.go)

## How it works
there is a cron defined in the [cron.go](https://github.com/SleepSpotify/SleepSpotify/blob/master/cron/cron.go) that will be executed every second. If there is a Pause object in the database with a `UTS > current UTS`: it will call the Spotify API to pause the music

## Dep
This project use [dep](https://github.com/golang/dep) to manage its dependencies.
And the dependencies are in this repo to avoid the risk of dead dependencies.

## License
this project is under the [MIT License](https://github.com/SleepSpotify/SleepSpotify/blob/master/LICENSE)
