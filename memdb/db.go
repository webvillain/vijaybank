package memdb

import (
	"bank/db"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Db *sql.DB
}

func NewDatabase() *Database {
	db, err := sql.Open("sqlite3", "./Bank.db")
	if err != nil {
		fmt.Println("Error While Connecting To Database !")
		log.Fatal(err)
	}
	fmt.Println("Connection To Database Is Successfull.")
	// to do
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS USERS (Id INTEGER PRIMARY KEY AUTOINCREMENT , Name TEXT , Email TEXT )")
	if err != nil {
		fmt.Println("Error While Creating Table")
		log.Fatal(err)
	}
	res, err := stmt.Exec()
	if err != nil {
		fmt.Println(err)
	}
	res.RowsAffected()
	fmt.Println("Table Is Crated Successfully into Database.")

	return &Database{
		Db: db,
	}
}

var ErrNotImplemented = errors.New("not implemnted")

func (d *Database) ListAllUsers(ctx context.Context) ([]*db.User, error) {
	// users := make([]*db.User, 0, len(d.users))
	// for _, user := range d.users {
	// 	users = append(users, user)
	// }
	// return users, nil
	var users []*db.User
	rows, err := d.Db.Query("SELECT * FROM USERS")

	if err != nil {
		log.Fatal(err)
	}
	var Id int
	var Name string
	var Email string
	for rows.Next() {

		rows.Scan(&Id, &Name, &Email)
	}
	users = append(users, &db.User{Id: Id, Name: Name, Email: Email})

	// defer d.Db.Close()
	// defer rows.Close()

	return users, nil
}

func (d *Database) UserById(ctx context.Context, Id int) (*db.User, error) {
	row, err := d.Db.Query("SELECT * FROM USERS WHERE Id = ?", Id)

	if err != nil {
		log.Fatal(err)
	}
	var id int
	var Name string
	var Email string
	if row.Next() {
		row.Scan(&id, &Name, &Email)
	}

	var user *db.User = &db.User{Id: id, Name: Name, Email: Email}
	// defer row.Close()

	return user, nil

	// u, ok := d.users[id]
	// if !ok {
	// 	return nil, db.ErrNotFound
	// }
	// return u, nil
}

/*
{"Name":"x","Email":"Y"}
*/
func (d *Database) CreateUser(ctx context.Context, u *db.User) (*db.User, error) {
	// u.ID = uuid.New().String()
	// d.users[u.ID] = u
	// return u, nil
	stmt, err := d.Db.Prepare("INSERT INTO USERS (Name , Email ) VALUES (?,?)")
	if err != nil {
		log.Fatal(err)
	}

	// we are not inserting Id Here Id Will Be Automatically Created In Database Because We Already Defines That Into Our Table Schema .
	stmt.Exec(u.Name, u.Email)

	return u, nil
}

/*
{"Name":"x1","Email":"Y1"}
*/
func (d *Database) UpdateUser(ctx context.Context, u *db.User) error {
	// _, ok := d.users[u.ID]
	// if !ok {
	// 	return db.ErrNotFound
	// }
	// d.users[u.ID] = u
	// return nil

	// to do
	stmt, err := d.Db.Prepare("UPDATE USERS SET Name = ? , Email = ? WHERE Id = ?")
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec(u.Name, u.Email, u.Id)

	return nil
}

func (d *Database) DeleteUser(ctx context.Context, Id int) error {
	// _, ok := d.users[id]
	// if !ok {
	// 	return db.ErrNotFound
	// }
	// delete(d.users, id)
	// return nil
	stmt, err := d.Db.Prepare("DELETE FROM USERS WHERE Id = ?")
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec(Id)
	return nil
}
