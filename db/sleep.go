package db

import (
	"encoding/json"
	"errors"
	"sync"

	"golang.org/x/oauth2"
)

// Sleep Object that represent a row in the database
type Sleep struct {
	ID    int64 `json:"-"`
	Uts   int64 `json:"uts"`
	token string
}

// NewSleep Constructor for the Sleep Object
func NewSleep(ID int64, tok *oauth2.Token, Uts int64) (Sleep, error) {
	s := Sleep{ID, Uts, ""}
	err := s.SetToken(tok)
	if err != nil {
		return Sleep{}, err
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

// GetFromID function to get the sleep object from the oauth2 token
func GetFromID(id int64) (Sleep, error) {
	rows, errSel := DB.Query("SELECT * FROM pause WHERE ID=?", id)
	if errSel != nil {
		return Sleep{}, errSel
	}

	if !rows.Next() {
		return Sleep{}, errors.New("Not Found")
	}

	var ret Sleep
	errScan := rows.Scan(&ret.ID, &ret.token, &ret.Uts)
	if errScan != nil {
		return Sleep{}, errScan
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
		errScan := rows.Scan(&tmp.ID, &tmp.token, &tmp.Uts)
		if errScan != nil {
			return nil, errScan
		}
		ret = append(ret, tmp)
	}
	return ret, nil
}

var mutexInsert = &sync.Mutex{}

// Insert function to insert a sleep obejct in the database
func (s *Sleep) Insert() error {
	mutexInsert.Lock()
	res, errQry := DB.Exec("INSERT INTO pause(token,uts) VALUES (?, ?)", s.token, s.Uts)
	if errQry != nil {
		mutexInsert.Unlock()
		return errQry
	}
	i, errID := res.LastInsertId()
	if errID != nil {
		mutexInsert.Unlock()
		return errID
	}
	s.ID = i
	mutexInsert.Unlock()
	return nil
}

// Update functionto Update the object
// call insert if the object hasn't an ID
func (s *Sleep) Update() error {
	if s.ID == 0 {
		return s.Insert()
	}
	_, errQry := DB.Exec("UPDATE pause SET token=?, uts=? WHERE ID=?;", s.token, s.Uts, s.ID)
	return errQry
}

// Delete function to delete an object
func (s *Sleep) Delete() error {
	if s.ID == 0 {
		return errors.New("Can't delete an uninserted object")
	}
	_, errQry := DB.Exec("DELETE FROM pause WHERE ID=?;", s.ID)
	return errQry
}