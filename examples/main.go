package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/sonyamoonglade/authio"
	"github.com/sonyamoonglade/authio/internal/gcmcrypt"
)

type DB struct {
	data map[int64]*User
}

func (db *DB) Add(u *User) int64 {
	id := int64(len(db.data) + 1)
	u.ID = id
	db.data[id] = u
	return id
}

func (db *DB) Delete(userID int64) {
	delete(db.data, userID)
}

func (db *DB) Get(userID int64) *User {
	return db.data[userID]
}

type User struct {
	ID      int64
	Name    string
	IsAdmin bool
}

func NewUser(name string, isAdmin bool) *User {
	return &User{Name: name, IsAdmin: isAdmin}
}

// Dont use global vars in your code.
var db *DB = &DB{}

// For the sake of simplicity
func init() {
	db.data = make(map[int64]*User)
}

const MyLabel = "your-great-label"
const MyAdminLabel = "some-admin-label"

func main() {

	setting := &authio.Setting{
		Label:    MyLabel,
		Name:     "SESSION",
		Path:     "/",
		Secret:   gcmcrypt.KeyFromString("asdfjkdasfasfjasfhsfjh"),
		Signed:   true,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteDefaultMode,
		Expires:  time.Hour * 1,
	}

	adminSetting := &authio.Setting{
		Label:    MyAdminLabel,
		Name:     "SESSION_ID",
		Path:     "/",
		Secret:   gcmcrypt.KeyFromString("!@0(ZXX123XZ!@#!@*(#))"),
		Signed:   true,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteDefaultMode,
		Expires:  time.Hour * 24,
	}

	auth := authio.NewAuthBuilder().
		UseLogger(nil).
		AddCookieSetting(setting).
		AddCookieSetting(adminSetting).
		UseStore(authio.NewInMemoryStore(&authio.Config{
			EntryTTL:         time.Hour * 1,
			OverflowStrategy: authio.LRU,
			ParseFunc:        authio.ToInt64,
		}, &authio.InMemoryConfig{
			MaxItems: 100,
		})).
		UseAuthioConfig(&authio.AuthioConfig{
			Paths: struct{ OnAuthNotRequired string }{
				OnAuthNotRequired: "/home", //your homepage
			},
		}).
		Build()

	handler := new(MyHandler)
	handler.auth = auth

	// This factory will make middlewares based on *MyLabel* setting
	// TODO: combine labels if needed
	authRequired := auth.AuthRequired(MyLabel)
	authRequiredAdmin := auth.AuthRequired(MyAdminLabel)

	//someOtherSettingAuthRequired := auth.AuthRequired(MyVerySecureLabel)
	//...

	// This factory will redirect users if they try to reach auth-unprotected endpoit while being authed.
	redirectAuthed := auth.RedirectAuthed(MyLabel)

	http.HandleFunc("/home", handler.home)
	http.HandleFunc("/register", redirectAuthed(handler.register))
	http.HandleFunc("/protected", authRequired(handler.greeting))
	http.HandleFunc("/logout", authRequired(handler.logout))

	// Add your own RBAC
	http.HandleFunc("/admin/ban", authRequiredAdmin(handler.banUser))

	fmt.Println("running at :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

type MyHandler struct {
	auth *authio.Auth
}

func (h *MyHandler) home(w http.ResponseWriter, r *http.Request) {
	for _, v := range db.data {
		fmt.Println(v.ID, v.Name, v.IsAdmin)
	}
	fmt.Fprintln(w, "hello darling!")
}

func (h *MyHandler) greeting(w http.ResponseWriter, r *http.Request) {

	userID, ok := authio.ValueFromContext[int64](r.Context())
	if !ok {
		///think of
	}

	u := db.Get(userID)

	fmt.Fprintf(w, "hello %s!\n", u.Name)
}

func (h *MyHandler) register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}

	isAdmin, err := strconv.ParseBool(r.URL.Query().Get("admin"))
	if err != nil {
		panic(err)
	}

	u := NewUser("Bobbby", isAdmin)
	userID := db.Add(u)

	var label string = MyLabel
	if isAdmin {
		label = MyAdminLabel
	}

	err = h.auth.SaveSession(w, label, authio.NewValueFromInt64(userID))
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		panic(err)
	}

	return
}

func (h *MyHandler) logout(w http.ResponseWriter, r *http.Request) {

	//if you need to get a session value that's asosiated with authSession
	//method 1:
	userID, ok := authio.ValueFromContext[int64](r.Context())
	if !ok {
		panic("internal error!!")
	}

	// method 2:

	// Same as authSession.Value.Raw().(int64)
	// userID = authSession.Raw().(int64)

	u := db.Get(userID)
	fmt.Printf("user %d;%s is logging out...\n", userID, u.Name)

	//This will remove a cookie and value inside an auth.store
	authSession := authio.SessionFromContext(r.Context())
	h.auth.InvalidateSessionByID(w, MyLabel, authSession.ID)

	fmt.Fprintf(w, "goodbye, %s!\n", u.Name)

	return
}

func (h *MyHandler) banUser(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()

	userIDtoBan := q.Get("id")

	// 1. Get an ID
	userIDInt, err := strconv.ParseInt(userIDtoBan, 10, 64)
	if err != nil {
		panic(err)
	}

	bannedUser := db.Get(userIDInt)

	// User that sent a request to this endpoint
	// is not nessesarily an admin, although he
	// passed adminAuthRequired middleware.
	// This is malicious, because user could've
	// just guess settings for 'admin' cookie
	// and write it on his own.
	//
	// Currently, implementing RBAC is your own deal.
	// authio can help preventing random sessions sent with a cookie
	// and setting userID(id who made a request) in your system to ctx
	// so you can access it and implement RBAC.
	// Just check if user that sent a request
	// is actually an admin by his sessionValue

	requesterID, ok := authio.ValueFromContext[int64](r.Context())
	_ = ok

	requester := db.Get(requesterID)

	if !requester.IsAdmin {
		w.Write([]byte("you are not an admin!\n"))
		return
	}

	// 2. Get session asosiated with this value(id)
	err = h.auth.InvalidateSessionByValue(authio.NewValueFromInt64(userIDInt))
	if err != nil {
		panic(err)
	}

	// 3. User is banned...
	fmt.Println("yahoo! user is banned")

	// Up to you to delete a user from database or not
	// db.Delete(userIDInt)

	fmt.Fprintf(w, "user: %s with id: %d has been banned\n", bannedUser.Name, bannedUser.ID)
}
