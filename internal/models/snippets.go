package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID 		int
	Title 	string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error){

	q := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY));`
	result, err := m.DB.Exec(q, title, content, expires)
	if err != nil{
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil{
		return 0, err
	}
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (Snippet, error){
	//Initialize snippet
	var snip Snippet
	q := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`
	if err := m.DB.QueryRow(q, id).Scan(&snip.ID, &snip.Title, &snip.Content, &snip.Created, &snip.Expires); err != nil{
		if errors.Is(err, sql.ErrNoRows){
			return snip, ErrNoRecord
		}else{
			return snip, err
		}
	}
	return snip, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error){
	//Initialize snippet
	var s []Snippet
	q := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`
	rows, err := m.DB.Query(q)
	if err != nil{
		return s, err
	}
	defer rows.Close()
	for rows.Next() {
        var snip Snippet
        if err := rows.Scan(&snip.ID, &snip.Title, &snip.Content, &snip.Created, &snip.Expires); err != nil{
            return nil, err
        }
        s = append(s, snip)
    }
	if err = rows.Err(); err != nil {
        return nil, err
    }
	return s, nil
}