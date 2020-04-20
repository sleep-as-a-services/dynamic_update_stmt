package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	_ "github.com/lib/pq" // here
)

func main() {
	db, err := sqlx.Open("postgres", "postgres:///sneak?sslmode=disable")
	if err != nil {
		panic(err)
	}

	req := &Request{
		ID:       "1",
		Username: "userPassword",
		Password: "simplePassword",
		Flag:     new(bool),
	}

	m, err := toMap(req)
	if err != nil {
		panic(err)
	}

	userRepo := UserRepo{db}
	userRepo.UpdateUserItems(context.Background(), m)
}

type UserRepo struct {
	db *sqlx.DB
}

func (u *UserRepo) UpdateUserItems(ctx context.Context, req map[string]interface{}) error {
	var (
		supportColumns = []string{"username", "password", "flag"}
		updateStmt     string
	)

	for _, column := range supportColumns {
		for key, value := range req {
			if column == key {
				updateStmt += fmt.Sprintf("%v=%#v,", column, value)
			}
		}
	}

	rawStmt := strings.Replace(fmt.Sprintf("update test_table set %s where id='%s';", updateStmt[:len(updateStmt)-1], req["id"]), "\"", "'", -1)

	stmt, err := u.db.Prepare(rawStmt)
	if err != nil {
		return errors.Wrap(err, "prepare stmt")
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return errors.Wrap(err, "update")
	}

	return nil
}

/* important! please aware. If update with zero values use pointer instead. */
type Request struct {
	ID       string `json:"id"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Flag     *bool  `json:"flag,omitempty"`
}

func toMap(obj *Request) (map[string]interface{}, error) {
	var req = make(map[string]interface{})

	if (Request{}) == *obj {
		//TODO: handle if empty
		panic("empty")
	}

	raw, err := json.Marshal(obj)
	if err != nil {
		return nil, errors.Wrap(err, "marshal request")
	}

	err = json.Unmarshal(raw, &req)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal to map")
	}

	return req, nil
}
