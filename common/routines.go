package common

import (
	"net/http"
	"errors"
	"encoding/json"
	"strings"
	"github.com/dchest/uniuri"
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

	//Just return status code
	if(value == nil){
		resp.WriteHeader(status)
		return status, nil
	}

	if j_bytes, err := json.Marshal(value); err != nil {
		resp.WriteHeader(status)
		return status, err
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

func StringJoin(sep string, elements ...string) string{ return strings.Join(elements, sep) }
func PathJoin(segs ...string) string {return StringJoin("/", segs...)}

func GetDefaultSecureHash() string { return uniuri.New() }
func GetSecureHash(length int) string { return uniuri.NewLen(length) }