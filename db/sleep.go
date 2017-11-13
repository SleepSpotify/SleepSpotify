package db

import (
	"encoding/json"
	"errors"
	"sync"

	"golang.org/x/oauth2"
)

// Sleep Object that represent a row in the database
type Sleep struct {
	ID      int64 `json:"-" xml:"-"`
	Uts     int64
	token   string
	refresh string
}

// NewSleep Constructor for the Sleep Object
func NewSleep(id int64, tok *oauth2.Token, uts int64) (*Sleep, error) {
	s := new(Sleep)
	s.ID = id
	s.Uts = uts
	s.refresh = tok.RefreshToken
	err := s.SetToken(tok)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// GetToken function to get a token from a object sleep
func (s *Sleep) GetToken() (*oauth2.Token, error) {
	var tok *oauth2.Token
	err := json.Unmarshal([]byte(s.token), &tok)
	return tok, err
}

// SetToken function to set the oauth token to an object sleep
func (s *Sleep) SetToken(token *oauth2.Token) error {
	tokByte, err := json.Marshal(token)
	if err != nil {
		return err
	}
	s.token = string(tokByte)
	return nil
}

// GetFromRefreshToken function to get the sleep timer from the refresh token
func GetFromRefreshToken(refresh string) (*Sleep, error) {
	rows, errSel := DB.Query("SELECT * FROM pause WHERE refresh=?", refresh)
	if errSel != nil {
		return new(Sleep), errSel
	}

	if !rows.Next() {
		return nil, nil
	}

	ret := new(Sleep)
	errScan := rows.Scan(&ret.ID, &ret.token, &ret.Uts, &ret.refresh)
	if errScan != nil {
		return new(Sleep), errScan
	}

	return ret, nil
}

// GetFromUts Function to retrieve all the pause object to pause at an uts
func GetFromUts(uts int64) ([]Sleep, error) {
	rows, err := DB.Query("SELECT * FROM pause WHERE uts<?", uts)
	if err != nil {
		return nil, err
	}

	var ret []Sleep
	for rows.Next() {
		var tmp Sleep
		errScan := rows.Scan(&tmp.ID, &tmp.token, &tmp.Uts, &tmp.refresh)
		if errScan != nil {
			return nil, errScan
		}
		ret = append(ret, tmp)
	}
	return ret, nil
}

var mutexInsert = &sync.Mutex{}

// Insert function to insert a sleep obejct in the database
func (s Sleep) Insert() error {
	mutexInsert.Lock()
	defer mutexInsert.Unlock()
	res, errQry := DB.Exec("INSERT INTO pause(token, uts, refresh) VALUES (?, ?, ?)", s.token, s.Uts, s.refresh)
	if errQry != nil {
		return errQry
	}
	i, errID := res.LastInsertId()
	if errID != nil {
		return errID
	}
	s.ID = i
	return nil
}

// Update function to Update the object
// call insert if the object hasn't an ID
func (s Sleep) Update() error {
	if s.ID == 0 {
		return s.Insert()
	}
	_, errQry := DB.Exec("UPDATE pause SET token=?, uts=?, refresh=? WHERE ID=?;", s.token, s.Uts, s.ID, s.refresh)
	return errQry
}

// Delete function to delete an object
func (s Sleep) Delete() error {
	if s.ID == 0 {
		return errors.New("Can't delete an uninserted object")
	}
	_, errQry := DB.Exec("DELETE FROM pause WHERE ID=?;", s.ID)
	return errQry
}
