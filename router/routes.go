package router

import (
	"github.com/gin-gonic/gin"
	"github.com/matac42/etChat/chatapp"
	"github.com/matac42/etChat/oauth"
	"github.com/matac42/etChat/user"
)

// Router creates a gin router.
func Router() {
	r := gin.Default()
	r.LoadHTMLGlob("html/*.html")

	cmelody := chatapp.CreateMelodyHandler()

	smiddle := user.CreateSessionMiddleware()
	r.Use(smiddle)

	v1 := r.Group("/")
	{
		v1.GET("chat", chatapp.ChatClient)
		v1.GET("ws", cmelody.MelodyClient)
		v1.GET("login", oauth.LogInClient)
		v1.GET("oauth", oauth.RedirectAuthenticateClient)
		v1.GET("callback", oauth.GetAccessTokenClient)
		v1.GET("session", user.SessionClient)
	}
	r.Run(":8080")
}
