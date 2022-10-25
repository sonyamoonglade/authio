package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sonyamoonglade/authio"
	"github.com/sonyamoonglade/authio/cookies"
	"github.com/sonyamoonglade/authio/gcmcrypt"
	"github.com/sonyamoonglade/authio/session"
	"github.com/sonyamoonglade/authio/store"
)

const MyLabel = "your-great-label"

func main() {

	setting := &cookies.Setting{
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

	auth := authio.NewAuthBuilder().
		UseLogger(nil).
		AddCookieSetting(setting).
		UseStore(store.NewInMemoryStore(&store.Config{
			EntryTTL:         time.Hour * 1,
			OverflowStrategy: store.LRU,
			ParseFunc:        store.ToInt64,
		}, &store.InMemoryConfig{
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

	// this factory will make middlewares based on *MyLabel* setting
	authRequired := auth.AuthRequired(MyLabel)

	//someOtherSettingAuthRequired := auth.AuthRequired(MyVerySecureLabel)
	//...

	// this factory will redirect users if they try to reach auth-unprotected endpoit while being authed.
	redirectAuthed := auth.RedirectAuthed(MyLabel)

	http.HandleFunc("/home", handler.home)
	http.HandleFunc("/register", redirectAuthed(handler.register))
	http.HandleFunc("/protected", authRequired(handler.greeting))
	http.HandleFunc("/logout", authRequired(handler.logout))

	fmt.Println("running at :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

type User struct {
	ID   int64
	Name string
}

type MyHandler struct {
	auth *authio.Auth
}

func (h *MyHandler) home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello darling!")
}

func (h *MyHandler) greeting(w http.ResponseWriter, r *http.Request) {

	userID, ok := authio.ValueFromContext[int64](r.Context())
	if !ok {
		///think of
	}

	fmt.Fprintf(w, "hello %d!\n", userID)
}

func (h *MyHandler) register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}

	name := "Bob" // e.g. passed with body
	var ID int64 = 5432
	u := &User{
		ID:   ID,
		Name: name,
	}

	err := h.auth.SaveSession(w, MyLabel, session.FromInt64(ID))
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

	authSession := authio.SessionFromContext(r.Context())

	//if you need to get a session value that's asosiated with authSession
	//method 1:
	userID, ok := authio.ValueFromContext[int64](r.Context())
	if !ok {
		panic("internal error!!")
	}

	//method 2:
	userID = authSession.Raw().(int64) //same as authSession.Value.Raw().(int64)

	fmt.Printf("user %d is logging out...\n", userID)

	//This will remove a cookie and value inside an auth.store
	h.auth.InvalidateSession(w, MyLabel, authSession.ID)

	fmt.Fprintf(w, "goodbye, user %d!\n", userID)
}
