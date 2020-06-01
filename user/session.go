package user

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// SessionClient fires fn when a /session connection.
func SessionClient(c *gin.Context) {
	if SessionNotFound() {
	}
}

// CreateSessionMiddleware create session middleware.
func CreateSessionMiddleware() gin.HandlerFunc {
	store := cookie.NewStore([]byte("secret"))
	mid := sessions.Sessions("chatuser", store)

	return mid
}

// SessionNotFound check a session existance.
func SessionNotFound() bool {
	b := true
	return b
}

// CreateSession create session based on a access token.
func CreateSession(c *gin.Context, t string) {
	session := sessions.Default(c)

	session.Set(t, "hello!(temporary)")
	session.Save()

}

// DeleteSession delete session.
func DeleteSession() {}
