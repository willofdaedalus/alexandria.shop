package main

import (
	"database/sql"
	"fmt"
	"log"
)

func main() {
    db, err  := sql.Open("sqlite3", "./users.db")
    if err != nil {
        log.Fatalf("open err: %s", err.Error())
    }

    err = createUsersTable(db)
    if err != nil {
        log.Fatalf("table err: %s", err.Error())
    }

    myUser := user{username: "manny", password: "mypassword"}
    err = signupUser(db, myUser)
    if err != nil {
        log.Fatalf("signup err: %s", err.Error())
    }
    fmt.Print("created new user\n")

    err = loginUser(db, myUser)
    if err != nil {
        log.Fatalf("signup login: %s", err.Error())
    }
    fmt.Print("user logged in successfully")
}
