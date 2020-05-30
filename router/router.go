package router

import (
	"github.com/gin-gonic/gin"
	chatapp "github.com/matac42/etChat/chat"
	"github.com/matac42/etChat/oauth"
)

// Router creates a gin router.
func Router() {
	r := gin.Default()
	r.LoadHTMLGlob("html/*.html")

	cmelody := chatapp.CreateMelodyHandler()

	v1 := r.Group("/")
	{
		v1.GET("chat", chatapp.ChatClient)
		v1.GET("ws", cmelody.MelodyClient)
		v1.GET("login", oauth.LogInClient)
		v1.GET("oauth", oauth.RedirectAuthenticateClient)
		v1.GET("callback", oauth.GetAccessTokenClient)
	}
	r.Run(":8080")
}
