package dbs

import (
	"database/sql"
	"fmt"
	_ "strconv"

	_ "github.com/go-sql-driver/mysql"
)

type (
	User struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
		Address  string `json:"address"`
	}

	Content struct {
		ContentPath string `json:"content"`
		ContentHash string `json:"content_hash"`
		Address     string `json:"address"`
		TokenId     string `json:"token_id"`
	}
)

var DBConn *sql.DB

func init() {
	DBConn = InitDB("root:system@tcp(127.0.0.1:3306)/standingbook?charset=utf8", "mysql")
}

func InitDB(connstr, Driver string) *sql.DB {
	db, error := sql.Open(Driver, connstr)
	if error != nil {
		fmt.Println("Error", error.Error())
	}
	return db
}

func (u User) Add() error {
	_, error := DBConn.Exec("insert into user(email,username,password,address) values(?,?,?,?)",
		u.Email, u.Username, u.Password, u.Address)
	if error != nil {
		fmt.Println("Failed to insert into t_user", error)
		return error
	}
	return nil
}

func (u *User) Query() (bool, error) {
	rows, err := DBConn.Query("select email, address from user where username=? and password=?", u.Username, u.Password)
	if err != nil {
		fmt.Println("Failed to find corresponding user", err)
		return false, err
	}
	if rows.Next() {
		err = rows.Scan(&u.Email, &u.Address)
		if err != nil {
			fmt.Println("Failed to bind properties to user", err)
			return false, err
		}
		return true, nil
	}
	return false, err
}

func (content *Content) AddContent() error {
	_, err := DBConn.Exec("insert into content(content,content_hash,address,token_id) values(?,?,?,?)",
		content.ContentPath, content.ContentHash, content.Address, content.TokenId)
	if err != nil {
		fmt.Println("failed to insert content", err)
		return err
	}
	return nil
}

func QueryContent(address string) ([]Content, error) {
	s := []Content{}
	rows, err := DBConn.Query("select content, content_hash, token_id from content where address=?", address)
	if err != nil {
		fmt.Println("failed to Query content via address", err)
		return s, err
	}

	var c Content
	for rows.Next() {
		err = rows.Scan(&c.ContentPath, &c.ContentHash, &c.TokenId)
		if err != nil {
			fmt.Println("failed to scan content properties", err)
			return s, err
		}
		s = append(s, c)
	}
	return s, nil

}
