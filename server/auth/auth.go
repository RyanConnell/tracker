package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"tracker/server/page"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

const (
	clientID     = "933303928886-nm68cv5c3rucdjjk8tvcrntlsk87u3u3.apps.googleusercontent.com"
	clientSecret = "XYx79SpBg9TlnTdpjSFL_zT-"
)

var conf = &oauth2.Config{
	ClientID:     clientID,
	ClientSecret: clientSecret,
	RedirectURL:  "http://localhost:8080/auth",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
	},
	Endpoint: google.Endpoint,
}

var sessionKey = randomString()
var store = sessions.NewCookieStore([]byte(sessionKey))

func RegisterHandlers() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/login", loginRequest)
	rtr.HandleFunc("/logout", logoutRequest)
	rtr.HandleFunc("/auth", authRequest)
	http.Handle("/", rtr)
}

func GetSession(r *http.Request, name string) (*sessions.Session, error) {
	return store.Get(r, name)
}

// State must be a randomly generated hash string.
func getLoginURL(state string) string {
	return conf.AuthCodeURL(state)
}

func randomString() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func loginRequest(w http.ResponseWriter, r *http.Request) {
	state := randomString()
	session, err := store.Get(r, "tracker")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	session.Values["state"] = state
	session.Save(r, w)

	http.Redirect(w, r, getLoginURL(state), http.StatusSeeOther)
}

func logoutRequest(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "tracker")
	if err != nil {
		serveError(err, w)
		return
	}

	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		serveError(err, w)
		return
	}

	http.Redirect(w, r, "http://localhost:8080/show", http.StatusSeeOther)
}

func authRequest(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "tracker")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	retrievedState := session.Values["state"]
	if retrievedState != r.URL.Query().Get("state") {
		serveError(fmt.Errorf("Retrieved State != Returned State"), w)
		return
	}

	token, err := conf.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))
	if err != nil {
		serveError(fmt.Errorf("Error exchanging token: %v\n", err), w)
		return
	}

	client := conf.Client(oauth2.NoContext, token)
	info, err := gatherUserInfo(client)
	if err != nil {
		serveError(err, w)
		return
	}

	user, err := LoadUser(info.Email)
	if err != nil {
		serveError(err, w)
		return
	}

	if user == nil {
		user, err = CreateUser(info)
		if err != nil {
			serveError(err, w)
			return
		}
	}

	fmt.Printf("User: %v\n", user)
	session.Options.Path = "/"
	session.Options.MaxAge = 86400 * 7

	// TODO: Change this to UUID instead of email?
	session.Values["user-id"] = user.Email
	session.Save(r, w)

	http.Redirect(w, r, "http://localhost:8080/show", http.StatusSeeOther)
}

func gatherUserInfo(client *http.Client) (*GoogleUserInfo, error) {
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve userinfo: %v", err)
	}
	defer resp.Body.Close()

	decode := json.NewDecoder(resp.Body)
	var jsonRep GoogleUserInfo
	decode.Decode(&jsonRep)
	return &jsonRep, nil
}

func serveError(err error, w http.ResponseWriter) {
	p := page.Page{[]byte(fmt.Sprintf("Error occured: %v", err.Error()))}
	p.ServePage(w)
}
