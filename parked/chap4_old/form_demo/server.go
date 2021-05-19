package main

import (
	"fmt"
	"log"
	"net/http"
)

func formHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "Form: %#v\n", r.Form)
	fmt.Fprintf(w, "PostForm: %#v\n", r.PostForm)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	htmlData := `<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8" />
		<title>Your first HTML form</title>
	</head>

	<body>
		<form action="/form" method="POST">
			<ul>
				<li>
					<label for="name">Name:</label>
					<input type="text" id="name" name="user_name" />
				</li>
				<li>
					<label for="mail">E-mail:</label>
					<input type="email" id="mail" name="user_mail" />
				</li>
				<li>
					<label for="msg">Message:</label>
					<textarea id="msg" name="user_message"></textarea>
				</li>
				<li class="button">
					<button type="submit">Send your message</button>
				</li>
			</ul>
		</form>
	</body>
</html>
`
	fmt.Fprintf(w, htmlData)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/form", formHandler)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
