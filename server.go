package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	http.Handle("/", http.FileServer(http.Dir("./files")))
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("./files"))))

	http.HandleFunc("/register", RegisterHandle)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/dashboard", DashBoardHandler)
	http.HandleFunc("/logout", LogoutHandler)
	http.HandleFunc("/update", UpdateHandler)
	http.HandleFunc("/delete", DeleteAccountHandler)
	// http.HandleFunc("/login.html", LoginHtmlhandler)

	fmt.Println("Starting server at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
