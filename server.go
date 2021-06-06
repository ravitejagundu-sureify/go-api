package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("our-secrete-key"))

type details struct {
	Username string
	Mailid   string
	Role     string
	Dept     string
	Jdate    string
	Company  string
	Quali    string
	Phone    string
	Address  string
}

func RenderTemplate(w http.ResponseWriter, filename string, p *details) error {
	t, err := template.ParseFiles(filename)
	// fmt.Println(p.company, p.username)

	if err != nil {

		return err
	}
	err = t.Execute(w, p)
	if err != nil {

		return err
	}

	return nil
}

func RegisterHandle(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "this page handles all registration requests")

	// check for any current session
	session, _ := store.Get(r, "sessions")
	_, ok := session.Values["username"]

	if ok {
		http.Redirect(w, r, "/dashboard", http.StatusFound)
		return
	}
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v \n", err)
		return
	}
	// fmt.Fprintf(w, "Post request succesful\n")
	username := r.FormValue("username")
	mailid := r.FormValue("mailid")
	password := r.FormValue("password")

	db, err := sql.Open("mysql", "root:raviteja@sureify@(172.17.0.2:3306)/php-app?parseTime=true")
	defer db.Ping()
	defer db.Close()

	if err != nil {
		fmt.Println("Database connection failed")
		return
	}

	query := `SELECT mail_id FROM users WHERE mail_id =?`

	rows, _ := db.Query(query, mailid)

	if rows == nil {

		query := `insert into users(mail_id,username,password) values(?,?,?)`

		_, err = db.Exec(query, mailid, username, password)

		if err != nil {
			fmt.Println("Insertion failed", err)
			return
		}
		fmt.Fprint(w, "<p>User succesfuly registered \n Head back to <a href='../login.html'>login page</a></p>")
	} else {
		fmt.Fprint(w, "<p>User Already registered \n Head back to <a href='../login.html'>login page</a></p>")
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "this page handles all the login requests")
	session, _ := store.Get(r, "sessions")
	_, ok := session.Values["username"]
	// fmt.Println(ok, "values")
	if !ok {

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v \n", err)
			return
		}

		mailid := r.FormValue("mail")
		password := r.FormValue("password")

		//establish a database connection
		db, err := sql.Open("mysql", "root:raviteja@sureify@(172.17.0.2:3306)/php-app?parseTime=true")

		if err != nil {
			fmt.Fprintf(w, "Database connection problem")
			return
		}
		// fmt.Println(mailid)
		query := `select password from users where mail_id = ?`
		var pw string
		err = db.QueryRow(query, mailid).Scan(&pw)

		if err != nil {
			fmt.Fprintf(w, "User email is not Registered. Error: %v", err)
			return
		}
		if pw != password {
			fmt.Fprintf(w, "Login unsuccesfull(Wrong Password)")
			return
		}
		session.Values["username"] = mailid
		session.Options.MaxAge = 7200
		session.Save(r, w)

	}
	http.Redirect(w, r, "/dashboard", http.StatusFound)

}
func DashBoardHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "user loged in succesully i guess.")

	// check for any current session
	session, _ := store.Get(r, "sessions")
	val, ok := session.Values["username"]

	if !ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	//establish a database connection
	db, err := sql.Open("mysql", "root:raviteja@sureify@(172.17.0.2:3306)/php-app?parseTime=true")
	defer db.Ping()
	defer db.Close()

	if err != nil {
		fmt.Fprintf(w, "Database connection problem")
		return
	}
	query := `select username from users where mail_id = ?`
	mailid := val.(string)
	data := &details{}
	_ = db.QueryRow(query, mailid).Scan(&data.Username)

	query = `select role,dept,company,qualification from Details where mail_id = ?`
	_ = db.QueryRow(query, mailid).Scan(&data.Role, &data.Dept, &data.Company, &data.Quali)
	data.Mailid = mailid
	// fmt.Fprintf(w, "<p>Username : %s\nMailid : %s\nCompany : %s<br>Role : %s<br>Department :%s<br>Qualification :%s<br></p>", username, mailid, company, role, dept, quali)

	// fmt.Println(data.username, data.company)
	RenderTemplate(w, "./files/dashboard.html", data)

}
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "This handles logut requests")

	// check for any current session
	session, _ := store.Get(r, "sessions")
	_, ok := session.Values["username"]

	if ok {
		session.Options.MaxAge = -1
		session.Save(r, w)
	}
	fmt.Fprintf(w, "<h1>Logut successfull</h1><a href='/'>Return to Home page</a>")
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Print("Update reuqests are handled here.")

	// check for any current session
	session, _ := store.Get(r, "sessions")
	name, ok := session.Values["username"]

	if !ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v \n", err)
		return
	}

	data := details{Mailid: name.(string)}
	//Retrieve form values
	data.Role = r.FormValue("role")
	data.Dept = r.FormValue("dept")
	data.Jdate = r.FormValue("jdate")
	data.Company = r.FormValue("company")
	data.Quali = r.FormValue("quali")
	data.Phone = r.FormValue("phone")
	data.Address = r.FormValue("address")

	//establish a database connection
	db, err := sql.Open("mysql", "root:raviteja@sureify@(172.17.0.2:3306)/php-app?parseTime=true")
	defer db.Ping()
	defer db.Close()

	if err != nil {
		fmt.Fprintf(w, "Database connection problem")
		return
	}
	query := `SELECT mail_id FROM Details WHERE mail_id =?`

	rows, _ := db.Query(query, data.Mailid)

	if rows == nil {

		query := `INSERT INTO Details VALUES(?,?,?,?,?,?,?,?)`

		_, err = db.Exec(query, data.Mailid, data.Role, data.Dept, data.Jdate, data.Company, data.Quali, data.Phone, data.Address)

		if err != nil {
			fmt.Println(err)
		}
	} else {
		query = `UPDATE Details SET role = ?,dept =?,jdate =?,company =?,qualification=?,phone_number =?,Address =? WHERE mail_id = ?`
		_, err = db.Exec(query, data.Role, data.Dept, data.Jdate, data.Company, data.Quali, data.Phone, data.Address, data.Mailid)
		if err != nil {
			fmt.Println(err)
		}
	}
	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "sessions")
	// check for any current session
	val, ok := session.Values["username"]
	if !ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	mailid := val.(string)

	//establish a database connection
	db, err := sql.Open("mysql", "root:raviteja@sureify@(172.17.0.2:3306)/php-app?parseTime=true")
	defer db.Ping()
	defer db.Close()
	if err != nil {
		fmt.Fprintf(w, "Database connection problem")
		return
	}

	// Delete perticular rows from users and details tables
	query := `DELETE FROM users WHERE mail_id = ?`
	_, _ = db.Exec(query, mailid)
	query = `DELETE FROM Details WHERE mail_id = ?`
	_, _ = db.Exec(query, mailid)

	// Delete session for that user
	session.Options.MaxAge = -1
	session.Save(r, w)
	fmt.Fprintf(w, "<h1>User Account deleted</h1><a href='/'>Return to Home page</a>")
}

func main() {

	http.Handle("/", http.FileServer(http.Dir("./files")))
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("./files"))))

	http.HandleFunc("/register", RegisterHandle)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/dashboard", DashBoardHandler)
	http.HandleFunc("/logout", LogoutHandler)
	http.HandleFunc("/update", UpdateHandler)
	http.HandleFunc("/delete", DeleteAccountHandler)

	fmt.Println("Starting server at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
