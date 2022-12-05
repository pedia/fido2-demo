package main

import (
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
		RPOrigin:      "https://wdemo.com:8443",             // The origin URL for WebAuthn requests
		RPIcon:        "https://wdemo.com:8443/s/logo.jpeg", // Optional icon URL for your site
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

	r.LoadHTMLGlob("templates/*.html")

	// navigator.credentials need tls
	r.RunTLS(":8443", "demo.pem", "demo.key")
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
	authSelect := protocol.AuthenticatorSelection{
		ResidentKey:        protocol.ResidentKeyRequirementPreferred,
		RequireResidentKey: protocol.ResidentKeyUnrequired(),
		UserVerification:   protocol.VerificationRequired, //
	}

	user := datastore.GetUser(name) // Find or create the new user
	options, sessionData, err := web.BeginRegistration(
		user,
		wan.WithAuthenticatorSelection(authSelect),
	)
	// handle errors if present
	// webauthn.SessionData{Challenge:"UILZsLXDM92UVGiKXrPT4-LRm0M6w4MCSEoDuy57hjg",
	// UserID:[]uint8{0x69, 0x64, 0x2d, 0x66, 0x6f, 0x6f},
	// AllowedCredentialIDs:[][]uint8(nil),
	// UserVerification:"required", Extensions:protocol.AuthenticationExtensions(nil)}
	datastore.SaveSession(sessionData)
	_ = err

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

	log.Printf("%#v\n%#v", options, sessionData)

	// options.publicKey contain our registration options
	c.HTML(http.StatusOK, "register.html", options)
}

func FinishRegistration(c *gin.Context) {
	// next: redirect url after finished
	name := c.Query("username")
	if name == "" {
		name = "foo"
	}

	// log.Printf("headers: %v", c.Request.Header)

	// TODO: read from body
	// var p []byte
	// if n, err := c.Request.Body.Read(p); err == nil {
	// 	 log.Printf("body: %d\n %s\n", n, string(p))
	// }

	user := datastore.GetUser(name) // Get the user
	sd := datastore.GetSession(name)
	if len(sd) == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// Get the session data stored from the function above
	// using gorilla/sessions it could look like this
	// sessionData := wan.SessionData{} // store.Get(r, "registration-session")
	parsedResponse, err := protocol.ParseCredentialCreationResponseBody(c.Request.Body)
	log.Printf("response: %#v", parsedResponse)
	_ = err

	credential, err := web.CreateCredential(user, sd[0], parsedResponse)
	log.Printf("credential: %s\n%#v", err, credential)

	// &webauthn.Credential{
	// 	ID:[]uint8{0xc1, 0x9c, 0xd8, 0xd1, 0x68, 0xc, 0xb6, 0x30, 0xa0, 0x3a, 0xa1, 0x7c, 0x3c, 0x6c, 0x59, 0xad},
	// 	PublicKey:[]uint8{0xa5, 0x1, 0x2, 0x3, 0x26, 0x20, 0x1, 0x21, 0x58, 0x20, 0x44, 0xef, 0xc2, 0x64, 0x33, 0xb2, 0x57, 0x31, 0x95, 0xbd, 0xaf, 0xd0, 0x5a, 0x32, 0x0, 0x8f, 0x0, 0x52, 0x7, 0x5a, 0xe1, 0xcc, 0xc7, 0xa3, 0x19, 0x4f, 0xf, 0xab, 0xc6, 0x7c, 0xb4, 0x2e, 0x22, 0x58, 0x20, 0x86, 0x55, 0x60, 0x34, 0xd1, 0x67, 0x22, 0x9a, 0x25, 0xdd, 0x24, 0x93, 0x23, 0x61, 0x4, 0x2c, 0x6b, 0xad, 0x28, 0xae, 0x88, 0x75, 0x3, 0xc3, 0xea, 0xff, 0xa3, 0x65, 0x71, 0x47, 0x74, 0xd4},
	// 	AttestationType:"none",
	// 	Authenticator:webauthn.Authenticator{
	// 			AAGUID:[]uint8{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	// 			SignCount:0x0,
	// 			CloneWarning:false
	// 		}
	// 	}
	datastore.SaveCredential(user, credential)

	// Handle validation or input errors
	// If creation was successful, store the credential object
	c.JSON(http.StatusOK, "Registration Success") // Handle next steps
}

func BeginLogin(c *gin.Context) {
	name := c.Query("username")
	if name == "" {
		name = "foo"
	}

	user := datastore.GetUser(name) // Find the user
	options, sessionData, err := web.BeginLogin(user)
	// handle errors if present
	_ = sessionData
	_ = err
	// store the sessionData values
	// c.JSON(http.StatusOK, options) // return the options generated
	// options.publicKey contain our registration options
	c.HTML(http.StatusOK, "login.html", options)
}

func FinishLogin(c *gin.Context) {
	// next:

	user := datastore.GetUser("") // Get the user
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
