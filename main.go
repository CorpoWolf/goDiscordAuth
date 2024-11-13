package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

var DiscordAuth = oauth2.Endpoint{
	AuthURL:   os.Getenv("DISCORD_CLIENT_AUTH_URL"),
	TokenURL:  "https://discord.com/api/oauth2/token",
	AuthStyle: oauth2.AuthStyleInParams,
}

var discordOAuthConfig = &oauth2.Config{
	ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
	ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
	RedirectURL:  "http://localhost:3000/callback",
	Scopes:       []string{"identify", "email"},
	Endpoint:     DiscordAuth, // Using the Discord endpoint defined above
}

type DiscordUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

func main() {
	fmt.Println("Hello Discord auth test")

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/callback", callbackHandler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	html := `<html><body><a href="/login">Login with Discord</a></body></html>`
	fmt.Fprint(w, html)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	url := discordOAuthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	if errMsg := r.URL.Query().Get("error"); errMsg != "" {
		http.Error(w, "OAuth error: "+errMsg, http.StatusInternalServerError)
		return
	}
	code := r.URL.Query().Get("code")
	token, err := discordOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Token exchange failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	client := discordOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://discord.com/api/users/@me")
	if err != nil {
		http.Error(w, "Failed to fetch user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var user DiscordUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		http.Error(w, "Failed to decode user response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("User %s has logged in via Discord with data: %+v\n", user.Username, user)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<html><body>
		<h1>Welcome, %s!</h1>
		<img src="%s" alt="Avatar of %s" />
		</body></html>`, user.Username, user.AvatarURL(), user.Username)
}

func (u *DiscordUser) AvatarURL() string {
	if u.Avatar == "" {
		return "https://cdn.discordapp.com/embed/avatars/0.png" // Default avatar URL if user has no custom avatar
	}
	return fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", u.ID, u.Avatar)
}
