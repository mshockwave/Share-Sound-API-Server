package common

import (
	"net/http"
	"errors"
	"encoding/json"
	"strings"
)

func GetSessionValue(req *http.Request, key interface{}) (interface{}, error) {
	s, err := SessionStorage.Get(req, USER_AUTH_SESSION)
	if err != nil { return nil, err }

	return s.Values[key], nil
}
func SetSessionValue(req *http.Request, resp http.ResponseWriter, key, value interface{}) error {
	//Ignore the error since sometimes the browser side coolie storage is broken
	//But we still can assign new cookies
	s, _ := SessionStorage.Get(req, USER_AUTH_SESSION)
	if s == nil { return errors.New("Session " + USER_AUTH_SESSION + " not available") }

	s.Values[key] = value
	return s.Save(req, resp)
}

func ResponseOkAsJson(resp http.ResponseWriter, value interface{}) (int, error){
	return ResponseStatusAsJson(resp, 200, value)
}
func ResponseStatusAsJson(resp http.ResponseWriter, status int, value interface{}) (int, error){
	if j_bytes, err := json.Marshal(value); err != nil {
		resp.WriteHeader(500)
		return 500, err
	}else{
		//Restore '&'
		str := string(j_bytes)
		str = strings.Replace(str, `\u0026`, "&", -1)

		resp.Header().Set("Content-Type", "application/json; charset=utf-8")
		resp.WriteHeader(status)
		_, err = resp.Write([]byte(str))
		return status, err
	}
}
