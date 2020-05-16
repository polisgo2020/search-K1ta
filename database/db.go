package database

import (
	"database/sql"
	"fmt"
)

type DB struct {
	*sql.DB
}

type TransactionErr struct {
	ExecErr     string
	RollbackErr string
}

const (
	addTitle   = "insert into titles (title) values ($1) on conflict (title) do update SET title = $1 returning id"
	addWord    = "insert into words (word) values ($1) on conflict (word) do update SET word = $1 returning id"
	getIndices = "select title_id from word_title where word_id = (select id from words where word = $1)"
	getTitle   = "select title from titles where id = $1"
	dropAll    = "drop table if exists word_title; drop table if exists words; drop table if exists  titles"
)

func Connect(host string, port string, user string, password string, dbName string) (*DB, error) {
	pgConString := fmt.Sprintf("port=%s host=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		port, host, user, password, dbName)
	db, err := sql.Open("postgres", pgConString)
	if err != nil {
		return &DB{}, err
	}
	err = db.Ping()
	return &DB{db}, err
}

func (db *DB) Init() error {
	_, err := db.Exec(`create table if not exists words
(
	id serial not null
		constraint words_pk
			primary key,
	word text
);

alter table words owner to postgres;

create unique index if not exists words_word_uindex
	on words (word);

create table if not exists titles
(
	id serial not null
		constraint titles_pk
			primary key,
	title text
);

alter table titles owner to postgres;

create unique index if not exists titles_title_uindex
	on titles (title);

create table if not exists word_title
(
	word_id integer not null
		constraint word_title_words_id_fk
			references words,
	title_id integer not null
		constraint word_title_titles_id_fk
			references titles,
	constraint word_title_pk
		primary key (word_id, title_id)
);

alter table word_title owner to postgres;
`)
	return err
}

func (db *DB) DropAll() (err error) {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("cannot begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("%s; cannot rollback: %w", err, rollbackErr)
			}
			err = fmt.Errorf("error on transaction: %w", err)
		}
	}()
	_, err = tx.Exec(dropAll)
	if err != nil {
		return
	}
	return tx.Commit()
}

func (db *DB) AddTitle(title string) (int64, error) {
	lastInsertedId := int64(-1)
	err := db.QueryRow(addTitle, title).Scan(&lastInsertedId)
	return lastInsertedId, err
}

func (db *DB) AddWord(word string) (int64, error) {
	lastInsertedId := int64(-1)
	err := db.QueryRow(addWord, word).Scan(&lastInsertedId)
	return lastInsertedId, err
}

func (db *DB) AddWordsIndices(wordId int64, indices []int64) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("cannot begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("%s; cannot rollback: %w", err, rollbackErr)
			}
			err = fmt.Errorf("error on transaction: %w", err)
		}
	}()
	for _, titleId := range indices {
		if _, err = tx.Exec("insert into word_title (word_id, title_id) values ($1, $2)", wordId, titleId); err != nil {
			return fmt.Errorf("cannot insert: %w", err)
		}
	}
	err = tx.Commit()
	return
}

func (db *DB) GetWordIndiced(word string) ([]int64, error) {
	rows, err := db.Query(getIndices, word)
	if err != nil {
		return nil, fmt.Errorf("error on get indices: %w", err)
	}
	res := make([]int64, 0)
	for rows.Next() {
		var index int64
		err = rows.Scan(&index)
		if err != nil {
			return nil, fmt.Errorf("error on scan: %w", err)
		}
		res = append(res, index)
	}
	return res, nil
}

func (db *DB) GetTitleById(id int64) (string, error) {
	var title string
	err := db.QueryRow(getTitle, id).Scan(&title)
	return title, err
}
