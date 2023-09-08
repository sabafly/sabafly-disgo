package main

import (
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/disgoorg/json"
	"github.com/disgoorg/log"

	"github.com/sabafly/sabafly-disgo"
	"github.com/sabafly/sabafly-disgo/bot"
	"github.com/sabafly/sabafly-disgo/discord"
	"github.com/sabafly/sabafly-disgo/oauth2"
)

var (
	letters      = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	token        = os.Getenv("disgo_token")
	clientSecret = os.Getenv("disgo_client_secret")
	baseURL      = os.Getenv("disgo_base_url")
	client       bot.Client
	oAuth2Client oauth2.Client
)

func main() {
	log.SetLevel(log.LevelDebug)
	log.Info("starting example...")
	log.Infof("disgo %s", disgo.Version)

	var err error
	client, err = disgo.New(token)
	if err != nil {
		log.Panic(err)
	}

	_, _ = client.Rest().UpdateApplicationRoleConnectionMetadata(client.ApplicationID(), []discord.ApplicationRoleConnectionMetadata{
		{
			Type:        discord.ApplicationRoleConnectionMetadataTypeIntegerGreaterThanOrEqual,
			Key:         "cookies_eaten",
			Name:        "Cookies Eaten",
			Description: "How many cookies have you eaten?",
		},
	})

	oAuth2Client = oauth2.New(client.ApplicationID(), clientSecret)

	mux := http.NewServeMux()
	mux.HandleFunc("/verify", handleVerify)
	mux.HandleFunc("/callback", handleCallback)
	_ = http.ListenAndServe(":6969", mux)
}

func handleVerify(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, oAuth2Client.GenerateAuthorizationURL(baseURL+"/callback", discord.PermissionsNone, 0, false, discord.OAuth2ScopeIdentify, discord.OAuth2ScopeRoleConnectionsWrite), http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	var (
		query = r.URL.Query()
		code  = query.Get("code")
		state = query.Get("state")
	)
	if code != "" && state != "" {
		session, _, err := oAuth2Client.StartSession(code, state)
		if err != nil {
			writeError(w, "error while starting session", err)
			return
		}

		user, err := oAuth2Client.GetUser(session)
		if err != nil {
			writeError(w, "error while getting user", err)
			return
		}

		_, err = oAuth2Client.UpdateApplicationRoleConnection(session, client.ApplicationID(), discord.ApplicationRoleConnectionUpdate{
			PlatformName:     json.Ptr("Cookie Monster " + user.Username),
			PlatformUsername: json.Ptr("Cookie Monster " + user.Tag()),
			Metadata: &map[string]string{
				"cookies_eaten": strconv.Itoa(rand.Intn(100)),
			},
		})
		if err != nil {
			writeError(w, "error while updating role connection", err)
			return
		}

		metadata, err := oAuth2Client.GetApplicationRoleConnection(session, client.ApplicationID())
		if err != nil {
			writeError(w, "error while getting role connection", err)
			return
		}

		data, _ := json.MarshalIndent(metadata, "", "\t")
		_, _ = w.Write([]byte("updated role connection:\n" + string(data)))

	}
}

func writeError(w http.ResponseWriter, text string, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(text + ": " + err.Error()))
}
