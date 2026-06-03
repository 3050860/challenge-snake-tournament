package models

//import (
//	"encoding/json"
//)
//
//type Error struct {
//	Err     error  `json:"-"`
//	Message string `json:"message,omitempty"`
//	Code    int    `json:"-"`
//}
//
//func (e *Error) Error() string {
//	return e.Err.Error()
//}
//
//func (e *Error) Unwrap() error { return e.Err }
//
//func (e *Error) Marshal() []byte {
//	bytes, err := json.Marshal(e)
//	if err != nil {
//		return nil
//	}
//	return bytes
//}
//
//func (e *Error) Status() int {
//	return e.Code
//}
