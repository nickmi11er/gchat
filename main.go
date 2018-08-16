package main

import (
	"gopkg.in/mgo.v2"

	"./chat"
	"./documents"
	"./models"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"log"
	"net/http"
)

var usersCollection *mgo.Collection

//var channelsCollection *mgo.Collection

func respondWithError(code int, message string, w http.ResponseWriter) {
	r := models.Response{Code: code, Status: "error", Params: map[string]string{"message": message}}
	w.Write(r.Json())
}

func respond(w http.ResponseWriter, code int, s string, p map[string]string) {
	r := models.Response{Code: code, Status: s, Params: p}
	w.Write(r.Json())
}

func methodNotAllowed(w http.ResponseWriter) {
	respond(w, 405, "method not allowed", map[string]string{})
}

// Auth middleware decorator
func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.FormValue("session_token")

		if token == "" {
			log.Println("User is not authorized")
			respond(w, 401, "unauthorized", map[string]string{})
			return
		}

		var user documents.UserDocument
		err := usersCollection.Find(bson.M{"session_token": token}).One(&user)
		if err != nil {
			log.Println("Invalid session token")
			respond(w, 401, "invalid session token", map[string]string{})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		methodNotAllowed(w)
	} else if r.Method == "POST" {
		r.ParseForm()
		log.Println(r.Form)
		//decoder := json.NewDecoder(r.Body)
		//decoder.Decode(map[string]string{"login":"", "password":""})

		login := r.Form.Get("username")
		password := r.Form.Get("password")
		pHash := PasswordHash(password)

		var user documents.UserDocument
		err := usersCollection.Find(bson.M{"username": login}).One(&user)
		if err != nil {
			log.Println("User not found")
			token := GenerateToken()
			userDoc := documents.UserDocument{Username: login, PasswordHash: pHash, SessionToken: token}
			usersCollection.Insert(userDoc)
			respond(w, 201, "created", userDoc.ToMap())
		} else if user.PasswordHash == pHash {
			respond(w, 200, "success", user.ToMap())
		} else {
			respondWithError(400, "Invalid login or password", w)
		}
	}
}

func chatHandler(w http.ResponseWriter, _ *http.Request) {
	t, err := template.ParseFiles("templates/chat.html")
	if err != nil {
		log.Println(err)
	}
	t.ExecuteTemplate(w, "chat", nil)
}

func pingHandler(w http.ResponseWriter, _ *http.Request) {
	respond(w, 200, "success", map[string]string{"action": "pong"})
}

func main() {
	log.SetFlags(log.Lshortfile)
	initDb()

	// websocket server
	server := chat.NewServer("/entry")
	go server.Listen()

	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/auth", authHandler)
	http.Handle("/chat", auth(http.HandlerFunc(chatHandler)))
	http.Handle("/ping", auth(http.HandlerFunc(pingHandler)))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initDb() {

	db := documents.Connection().DB
	usersCollection = db.C("users")
	//channelsCollection = db.C("channels")
}
