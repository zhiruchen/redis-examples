package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zhiruchen/redis-examples/db"
)

func TestCreateUser(t *testing.T) {
	setUpDB()
	defer tearDown()

	type bodyType struct {
		Code int    `json:"code"`
		Msg  string `json:"message"`
	}

	cases := []struct {
		body string
		code int
	}{
		{
			body: `{"name": "test_user"}`,
			code: http.StatusOK,
		},
		{
			body: `{"name": ""}`,
			code: invalidParamCode,
		},
	}

	for _, cc := range cases {
		reader := bytes.NewReader([]byte(cc.body))
		req, _ := http.NewRequest("POST", "/users", reader)

		w := httptest.NewRecorder()
		r := getRouter()
		r.ServeHTTP(w, req)

		bodyS := w.Body.String()
		t.Logf("body: %s\n", bodyS)

		resp := &bodyType{}
		err := json.Unmarshal([]byte(bodyS), &resp)
		if err != nil {
			t.Errorf("parse body error: %v\n", err)
		}

		if resp.Code != cc.code {
			t.Errorf("expect: %d, get: %d\n", cc.code, w.Code)
		}
	}
}

func TestGetUser(t *testing.T) {

}

func TestUpdateUser(t *testing.T) {

}

func setUpDB() {
	if err := db.InitMysql(); err != nil {
		panic(err)
	}

	if err := db.NewRedisClient(); err != nil {
		panic(err)
	}
}

func tearDown() {
	db.ORM.Exec("TRUNCATE TABLE user")
}
