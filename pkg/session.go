package pkg

import (
"github.com/gin-gonic/contrib/sessions"
"github.com/gin-gonic/gin"
"log"
"net/http"
)

const (
	userkey = "user"
)

// gin session key
// const userkey = "user"

func EnableCookieSession() gin.HandlerFunc {
	// store := sessions.NewCookieStore([]byte("secret"))
	// store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	store, err := sessions.NewRedisStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	if err != nil {
		log.Fatalf("sessions.NewRedisStore: %v", err)
	}
	return sessions.Sessions("mysession", store)
}

func AuthSessionMiddle(c *gin.Context) {
	session := sessions.Default(c)
	sessionValue := session.Get("user")
	if sessionValue == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		c.Abort()
		return
	}
	c.Next()
	return
}

// func AuthRequired(c *gin.Context) {
//     session := sessions.Default(c)
//     user := session.Get("user")
//     if user == nil {
//         // Abort the request with the appropriate error code
//         c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
//         return
//     }
//     // Continue down the chain to handler etc
//     c.Next()
// }

func SaveAuthSession(c *gin.Context, username string) {
	session := sessions.Default(c)
	session.Set("user", username)
	err := session.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save session"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "save session ok"})
}

func ClearAuthSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
}

func HasSession(c *gin.Context) bool {
	session := sessions.Default(c)
	if sessionValue := session.Get("userId"); sessionValue == nil {
		return false
	}
	return true
}

func GetSessionUserId(c *gin.Context) uint {
	session := sessions.Default(c)
	sessionValue := session.Get("userId")
	if sessionValue == nil {
		return 0
	}
	return sessionValue.(uint)
}

// func GetUserSession(c *gin.Context) map[string]interface{} {
//
//     hasSession := HasSession(c)
//     userName := ""
//     if hasSession {
//         userId := GetSessionUserId(c)
//         userName = models.UserDetail(userId).Name
//     }
//     data := make(map[string]interface{})
//     data["hasSession"] = hasSession
//     data["userName"] = userName
//     return data
// }


// AuthRequired is a simple middleware to check the session
func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		// Abort the request with the appropriate error code
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// Continue down the chain to handler etc
	c.Next()
}



func Logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}
	session.Delete(userkey)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func Me(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "You are logged in"})
}
