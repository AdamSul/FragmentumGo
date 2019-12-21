// model.go
package main

import (
	"database/sql"
	"fmt"
)

type fragment struct {
	ID      int            `json: "id"`
	content sql.NullString `json: "content"`
	pre     sql.NullString `json: "pre"`
	post    sql.NullString `json: "post"`
	parent  sql.NullInt64  `json: "parent"`
}

func getFragment(db *sql.DB, fragmentID int) (fragment, error) {
	statement := fmt.Sprintf("SELECT id, pre, content, post, parent FROM fragment where id=%d and context=1", fragmentID)
	var u fragment
	err := db.QueryRow(statement).Scan(&u.ID, &u.pre, &u.content, &u.post, &u.parent)

	if err != nil {
		if err == sql.ErrNoRows {
			return u, err // there were no rows, but otherwise no error occurred
		} else {
			return u, fmt.Errorf("getFragment 26: %v", err)
		}
	}

	return u, nil
}

func getSubfragments(db *sql.DB, parentID int) ([]int, error) {
	statement := fmt.Sprintf("SELECT id FROM fragment where parent =%d and context=1 order by sequenceaschild", parentID)
	row, err := db.Query(statement)
	var u fragment

	if err != nil {
		return nil, fmt.Errorf("getSubfragments 39: %v", err)
	}

	var s []int
	//debug: fmt.Println(statement)
	defer row.Close()

	for row.Next() {
		if err := row.Scan(&u.ID); err != nil {
			return s, fmt.Errorf("getSubfragments 48: %v", err)
		}
		s = append(s, u.ID)
		//debug: fmt.Println(s)
		//debug: fmt.Println(u)
	}
	err = row.Err()
	if err != nil {
		return s, fmt.Errorf("getSubfragments 56: %v", err)
	}
	return s, nil
}
