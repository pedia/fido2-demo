package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/protocol/webauthncose"
	wan "github.com/duo-labs/webauthn/webauthn"
	"github.com/gin-gonic/gin"
)

var (
	web       *wan.WebAuthn
	err       error
	datastore Store  = NewStore()
	site      string = "https://9aedu.net"
)

func session_from(c *gin.Context) (*User, *wan.SessionData) {
	if is, err := c.Cookie("sid"); err == nil {
		if sid, err := strconv.Atoi(is); err == nil {
			return datastore.GetSession(sid)
		}
	}
	return nil, nil
}

func session_to(c *gin.Context, sid int) {
	c.SetCookie("sid", fmt.Sprintf("%d", sid), 0, "", "", true, true)
}

// Your initialization function
func main() {
	web, err = wan.New(&wan.Config{
		RPDisplayName: "Demo(site)",                    // Display Name for your site
		RPID:          "9aedu.net",                     // Generally the FQDN for your site
		RPOrigin:      "https://9aedu.net",             // The origin URL for WebAuthn requests
		RPIcon:        "https://9aedu.net/s/logo.jpeg", // Optional icon URL for your site
		Debug:         true,
		Timeout:       360000, // 6 minutes
	})
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	// TODO: not work
	// r.StaticFileFS("/s", "logo.jpeg", gin.Dir("s", false))

	r.Static("/s", "s")

	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/register", BeginRegistration)
	r.POST("/register", FinishRegistration)

	r.GET("/login", BeginLogin)
	r.POST("/login", FinishLogin)
	// log.Println("listen on http://0.0.0.0:8443")

	r.LoadHTMLGlob("templates/*")

	// navigator.credentials need tls
	r.RunTLS(":443", "demo.pem", "demo.key")
}

func BeginRegistration(c *gin.Context) {
	name := c.Query("username")
	if name == "" {
		name = "foo"
	}

	// chrome need this
	// "authenticatorSelection": {
	//    "residentKey": "preferred",
	//    "requireResidentKey": false,
	//    "userVerification": "required"
	// }

	// https://demo.yubico.com/webauthn-technical
	// "authenticatorSelection": {
	//	  "requireResidentKey": false,
	//	  "userVerification": "discouraged"
	//  },
	authenticatorSelection := protocol.AuthenticatorSelection{
		// ResidentKey:        protocol.ResidentKeyRequirementPreferred,
		RequireResidentKey: protocol.ResidentKeyUnrequired(),
		UserVerification:   protocol.VerificationDiscouraged, //
	}

	user := datastore.GetUser(name) // Find or create the new user
	options, sessionData, err := web.BeginRegistration(
		user,
		wan.WithAuthenticatorSelection(authenticatorSelection),
	)
	if err != nil {
		c.HTML(http.StatusOK, "index.html", err)
		return
	}

	// webauthn.SessionData{Challenge:"UILZsLXDM92UVGiKXrPT4-LRm0M6w4MCSEoDuy57hjg",
	// UserID:[]uint8{0x69, 0x64, 0x2d, 0x66, 0x6f, 0x6f},
	// AllowedCredentialIDs:[][]uint8(nil),
	// UserVerification:"required", Extensions:protocol.AuthenticationExtensions(nil)}
	sid, err := datastore.SaveSession(sessionData)
	if err != nil {
		c.HTML(http.StatusOK, "index.html", err)
		return
	}

	session_to(c, sid)

	// store the sessionData values
	options.Response.Parameters = []protocol.CredentialParameter{
		{
			Type:      protocol.PublicKeyCredentialType,
			Algorithm: webauthncose.AlgES256, // -7
		},
		{
			Type:      protocol.PublicKeyCredentialType,
			Algorithm: webauthncose.AlgRS256, // -257
		},
	}

	options.Response.Attestation = protocol.PreferDirectAttestation

	// log.Printf("%#v\n%#v", options, sessionData)

	// options.publicKey contain our registration options
	c.HTML(http.StatusOK, "register.html", map[string]any{"Opts": options, "Username": name})
}

func to_json(o any) string {
	if bs, erj := json.Marshal(o); erj == nil {
		return string(bs)
	}
	return ""
}

func FinishRegistration(c *gin.Context) {
	parsedResponse, err := protocol.ParseCredentialCreationResponseBody(c.Request.Body)
	log.Printf("response: %s", to_json(parsedResponse))
	if err != nil {
		c.HTML(http.StatusOK, "index.html", err)
		return
	}

	// Get the session data stored from the function above
	// using gorilla/sessions it could look like this
	u, sd := session_from(c)
	if sd == nil {
		c.HTML(http.StatusOK, "index.html", errors.New("no session"))
		return
	}

	credential, err := web.CreateCredential(u, *sd, parsedResponse)
	// log.Printf("credential: %s\n%#v", err, credential)

	if err == nil {
		// Handle validation or input errors
		// If creation was successful, store the credential object
		datastore.SaveCredential(u, credential)

		c.JSON(http.StatusOK, "Registration Success") // Handle next steps
		return
	}

	c.JSON(http.StatusOK, "Registration Failed") // Handle next steps
}

func BeginLogin(c *gin.Context) {
	name := c.Query("username")
	if name == "" {
		name = "foo"
	}

	user := datastore.GetUser(name) // Find the user
	options, sessionData, err := web.BeginLogin(
		user,
		func(pkcro *protocol.PublicKeyCredentialRequestOptions) {
			pkcro.UserVerification = protocol.VerificationRequired
		},
	)
	// handle errors if present
	if err != nil {
		c.HTML(http.StatusOK, "index.html", err)
		return
	}

	sid, err := datastore.SaveSession(sessionData)
	if err != nil {
		c.HTML(http.StatusOK, "index.html", err)
		return
	}

	session_to(c, sid)

	// store the sessionData values
	// c.JSON(http.StatusOK, options) // return the options generated
	// options.publicKey contain our registration options
	c.HTML(http.StatusOK, "login.html", map[string]any{"Opts": options, "Username": name})
}

func FinishLogin(c *gin.Context) {
	u, sd := session_from(c)
	if sd == nil {
		c.HTML(http.StatusOK, "index.html", errors.New("no session"))
		return
	}

	if parsedResponse, err := protocol.ParseCredentialRequestResponseBody(
		c.Request.Body,
	); err == nil {
		credential, err := web.ValidateLogin(
			u, *sd, parsedResponse,
		)
		log.Printf("login %s\n%s", err, to_json(credential))
		if err == nil {
			c.JSON(http.StatusOK, "Login Success")
			return
		}
	}

	c.JSON(http.StatusOK, "Login Failed")
}
