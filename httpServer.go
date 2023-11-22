package main

import (
    "net/http"
	"html/template"
	"database/sql"
	"log"
	_ "modernc.org/sqlite"
)

type LoginInfo struct{
	Email string
	Password string
}

func main() {
	//connect to database
	db, err := sql.Open("sqlite", "test.db")
	if err != nil{
		log.Println("Error connecting to database!", err)
	}
	defer db.Close()
	//create table if it does not exist
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS user(email VARCHAR PRIMARY KEY, password VARCHAR)")
	if err != nil{
		log.Println("Error creating table!", err)
	}
	stmt.Exec()
	//insert value
	result, err := db.Exec("INSERT OR IGNORE INTO user (email, password) VALUES (?, ?)", "abc@gmail.com","123")
	if err != nil{
		log.Println("Error inserting value!", err)
	}
	_ = result
	//select/query data in database
	rows, err := db.Query("SELECT email, password FROM user")
	if err != nil{
		log.Println("Error selecting query!", err)
	}
	defer rows.Close()

	//create login template, handle form request
	tmpl := template.Must(template.ParseFiles("login.html"))
    http.HandleFunc("/login", func (w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
			tmpl.Execute(w,nil)
			return
		}
	//get value from user input
		logInfo := LoginInfo{
			Email: r.FormValue("email"),
			Password: r.FormValue("password"),
		}

	//compare user input with database
		var dbEmail, dbPassword string
		tmplResult := struct{ Success, Error_em, Error_pw bool }{false, false, false}

		err := db.QueryRow("SELECT email, password FROM user WHERE email=?", logInfo.Email).Scan(&dbEmail, &dbPassword)
		if err != nil {
			tmplResult.Error_em = true
		} else if logInfo.Password == dbPassword {
			tmplResult.Success = true
		} else {
			tmplResult.Error_pw = true
		}

		tmpl.Execute(w, tmplResult)

    })
	//create server
    fs := http.FileServer(http.Dir("static/"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    http.ListenAndServe(":5500", nil)
}

