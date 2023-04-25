package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

type CtxKey string

const (
	DB     CtxKey = "db"
	DBName CtxKey = "dbname"
	DBCli  CtxKey = "dbcli"
	User   CtxKey = "user"
)

const HashCost = 14

func JSONResponse(w http.ResponseWriter, m interface{}) {
	j, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// this unmarshals json
func UnmarshalObject(body io.Reader, obj interface{}) error {
	return json.NewDecoder(body).Decode(obj)
}
