package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/duo-labs/webauthn/protocol"
	wan "github.com/duo-labs/webauthn/webauthn"
	"github.com/gin-gonic/gin"
)

var (
	web       *wan.WebAuthn
	err       error
	datastore Store  = NewStore()
	site      string = "https://wdemo.com"
)

// Your initialization function
func main() {
	web, err = wan.New(&wan.Config{
		RPDisplayName: "Demo(site)",                         // Display Name for your site
		RPID:          "wdemo.com",                          // Generally the FQDN for your site
		RPOrigin:      "https://wdemo.com",                  // The origin URL for WebAuthn requests
		RPIcon:        "https://wdemo.com:8443/s/logo.jpeg", // Optional icon URL for your site
		Debug:         true,
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
	// log.Println("listen on http://0.0.0.0:8443")

	r.LoadHTMLGlob("templates/*")

	// navigator.credentials need tls
	r.RunTLS(":8443", "demo.pem", "demo.key")
}

func BeginRegistration(c *gin.Context) {
	// chrome need this
	// "authenticatorSelection": {
	//    "residentKey": "preferred",
	//    "requireResidentKey": false,
	//    "userVerification": "required"
	// }
	authSelect := protocol.AuthenticatorSelection{
		ResidentKey:        protocol.ResidentKeyRequirementPreferred,
		RequireResidentKey: protocol.ResidentKeyUnrequired(),
		UserVerification:   protocol.VerificationRequired,
	}

	user := datastore.GetUser() // Find or create the new user
	options, sessionData, err := web.BeginRegistration(
		user,
		wan.WithAuthenticatorSelection(authSelect),
	)
	// handle errors if present
	// webauthn.SessionData{Challenge:"UILZsLXDM92UVGiKXrPT4-LRm0M6w4MCSEoDuy57hjg",
	// UserID:[]uint8{0x69, 0x64, 0x2d, 0x66, 0x6f, 0x6f},
	// AllowedCredentialIDs:[][]uint8(nil),
	// UserVerification:"required", Extensions:protocol.AuthenticationExtensions(nil)}
	_ = sessionData
	_ = err
	_ = options
	// store the sessionData values
	// options.Response.Parameters = []protocol.CredentialParameter{
	// 	{
	// 		Type:      protocol.PublicKeyCredentialType,
	// 		Algorithm: webauthncose.AlgES256,
	// 	},
	// 	{
	// 		Type:      protocol.PublicKeyCredentialType,
	// 		Algorithm: webauthncose.AlgRS384,
	// 	},
	// }

	// options.Response.AuthenticatorSelection

	fmt.Printf("%#v\n%#v", options, sessionData)

	// options.publicKey contain our registration options
	c.HTML(http.StatusOK, "register.html", options)
}

func FinishRegistration(c *gin.Context) {
	user := datastore.GetUser() // Get the user
	// Get the session data stored from the function above
	// using gorilla/sessions it could look like this
	sessionData := wan.SessionData{} // store.Get(r, "registration-session")
	parsedResponse, err := protocol.ParseCredentialCreationResponseBody(c.Request.Body)
	_ = err
	credential, err := web.CreateCredential(user, sessionData, parsedResponse)
	_ = credential
	// Handle validation or input errors
	// If creation was successful, store the credential object
	c.JSON(http.StatusOK, "Registration Success") // Handle next steps
}

func BeginLogin(c *gin.Context) {
	user := datastore.GetUser() // Find the user
	options, sessionData, err := web.BeginLogin(user)
	// handle errors if present
	_ = sessionData
	_ = err
	// store the sessionData values
	c.JSON(http.StatusOK, options) // return the options generated
	// options.publicKey contain our registration options
}

func FinishLogin(c *gin.Context) {
	user := datastore.GetUser() // Get the user
	// Get the session data stored from the function above
	// using gorilla/sessions it could look like this
	sessionData := wan.SessionData{} // store.Get(r, "login-session")
	parsedResponse, err := protocol.ParseCredentialRequestResponseBody(c.Request.Body)
	_ = err
	credential, err := web.ValidateLogin(user, sessionData, parsedResponse)
	_ = credential
	// Handle validation or input errors
	// If login was successful, handle next steps
	c.JSON(http.StatusOK, "Login Success")
}
