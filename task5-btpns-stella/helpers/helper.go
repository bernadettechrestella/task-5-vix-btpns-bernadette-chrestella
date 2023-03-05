package helpers

import (
	errLog "errors"
	app "task5-btpns-stella/app"
)

func LogError(hint string, err error, callback ErrorFunc) {
	if err != nil {
		callback(hint + " " + err.Error())
	}
}

func ValidateUser(req app.Login) error {
	if req.Username == "" {
		return errLog.New("Username Cannot Empty")
	}
	if req.Password == "" {
		return errLog.New("Password Cannot Empty")
	}
	return nil
}
