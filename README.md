# SleepSpotify v0.1
This project allows you to pause spotify at the end of a timer or after a certain timer

## What you need
- GoLang (^1.9)
- MariaDB/MySQL (^10.1.23/^5.7.19)

## How to install
### MySQL Part

Create a database and an user as you wish. After that, use this script to init the database

```sql
CREATE TABLE `pause` (
  `ID` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `token` text NOT NULL,
  `Uts` int(10) unsigned NOT NULL,
  `refresh` text NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8
```

### Golang Part
```bash
go get -u github.com/SleepSpotify/SleepSpotify
cd $GOPATH/src/github.com/SleepSpotify/SleepSpotify
cp config/config.yaml.sample config/config.yaml
```

edit the config/config.yaml file

```bash
go build
./SleepSpotify
```

## API
this project is an API.:
- Account
  - GET `/login` the call to redirect to the spotify API to connect with the Spotify OAuth2 sevice
  - GET `/callback` the call that the Spotify OAuth2 service will call after a successfull authentification
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
  - POST `/spotify/sleep` to set a timer.
    -  body params
      - uts => an uts to stop the timer  
    - response
      - JSONSleepFound
      - JSONError
  - PUT `/spotify/sleep` to update a timer.
    -  body params
      - uts => an uts to stop the timer (int64)
    - response
      - JSONSleepFound
      - JSONError
  - DELETE `/spotify/sleep` to delete a timer.
    - response
      - JSONActionDone
      - JSONError

All the response object are [defined here](https://github.com/SleepSpotify/SleepSpotify/blob/master/controler/json.go)

## How it works
there is a cron defined in the [cron.go](https://github.com/SleepSpotify/SleepSpotify/blob/master/cron.go) that will be executed every seconde. If there is a Pause object in the database to pause it will execute the call to the Spotify API

## Dep
This project use [dep](https://github.com/golang/dep) to manage its depandancies.
And the depandancies are in this repo because this way no risk of dead depandancies

## License
this project is under the [MIT Licence](https://github.com/SleepSpotify/SleepSpotify/blob/master/LICENSE)
