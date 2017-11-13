package controler

import "github.com/SleepSpotify/SleepSpotify/db"

// JSONSleepFound Object to return if a sleep is found
type JSONSleepFound struct {
	Found bool
	Sleep *db.Sleep
}

// JSONConnected Object to return is the user is connected
type JSONConnected struct {
	IsConnected bool
	Username    string
}

// JSONError Object to return an error
type JSONError struct {
	Message string
}

// JSONActionDone Object to return an action done
type JSONActionDone struct {
	Message string
}
