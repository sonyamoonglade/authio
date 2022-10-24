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
		Build()

	handler := new(MyHandler)
	handler.auth = auth

	authRequired := auth.HTTPGetSessionWithLabel //change for wrapped handler
	http.HandleFunc("/protected", authRequired(handler.greeting, MyLabel))
	http.HandleFunc("/register", handler.register)

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
