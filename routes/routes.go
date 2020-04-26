package routes

import (
	"net/http"

	"../middleware"
	"../models"
	"../sessions"
	"../utils"
	"github.com/gorilla/mux"
)

func indexhandler(w http.ResponseWriter, r *http.Request) {
	comments, err := models.GetAllComments()
	if len(comments) == 0 {
		comments = nil
	}
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	utils.ExecuteTemplate(w, "index.html", struct {
		Title       string
		Comments    []*models.Comment
		DisplayForm bool
	}{"Comments", comments, true})
}

func indexPostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "session")
	untypedUserID := session.Values["user_id"]
	userID, ok := untypedUserID.(int64)
	if !ok {
		utils.InternalServerError(w)
		return
	}

	r.ParseForm()
	comment := r.PostForm.Get("comment")
	err := models.PostComment(userID, comment)
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	http.Redirect(w, r, "/index", 302)
}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteTemplate(w, "login.html", nil)
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	user, err := models.AuthenticateUser(username, password)
	if err != nil {
		switch err {
		case models.ErrUserNotFound:
			utils.ExecuteTemplate(w, "login.html", "unknown user")
		case models.ErrInvalidLogin:
			utils.ExecuteTemplate(w, "login.html", "Incorrect password")
		default:
			utils.InternalServerError(w)
		}
		return
	}
	userID, err := user.GetID()
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	session, _ := sessions.Store.Get(r, "session")
	session.Values["user_id"] = userID
	session.Save(r, w)
	http.Redirect(w, r, "/index", 302)
}

func registerGetHandler(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteTemplate(w, "register.html", nil)
}

func registerPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	err := models.RegisterUser(username, password)
	if err == models.ErrUsernameTaken {
		utils.ExecuteTemplate(w, "register.html", "Username taken")
		return
	} else if err != nil {
		utils.InternalServerError(w)
		return
	}
	http.Redirect(w, r, "/login", 302)
}

func userGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "session")
	untypedUserID := session.Values["user_id"]
	sessUserID, ok := untypedUserID.(int64)
	if !ok {
		utils.InternalServerError(w)
		return
	}
	vars := mux.Vars(r)
	username := vars["username"]
	user, err := models.GetUserByUsername(username)
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	userID, err := user.GetID()
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	comments, err := models.GetComments(userID)
	if len(comments) == 0 {
		comments = nil
	}
	if err != nil {
		utils.InternalServerError(w)
		return
	}
	utils.ExecuteTemplate(w, "index.html", struct {
		Title       string
		Comments    []*models.Comment
		DisplayForm bool
	}{username, comments, sessUserID == userID})
}

func logoutGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Store.Get(r, "session")
	delete(session.Values, "user_id")
	session.Save(r, w)
	http.Redirect(w, r, "/login", 302)
}

//NewRouter This is a new router
func NewRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/index", middleware.AuthRequired(indexhandler)).Methods("GET")
	router.HandleFunc("/index", indexPostHandler).Methods("POST")
	router.HandleFunc("/login", loginGetHandler).Methods("GET")
	router.HandleFunc("/login", loginPostHandler).Methods("POST")
	router.HandleFunc("/logout", logoutGetHandler).Methods("GET")
	router.HandleFunc("/register", registerGetHandler).Methods("GET")
	router.HandleFunc("/register", registerPostHandler).Methods("POST")
	router.HandleFunc("/{username}", middleware.AuthRequired(userGetHandler)).Methods("GET")

	fs := http.FileServer(http.Dir("./static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	return router
}
