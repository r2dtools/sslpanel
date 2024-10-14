package middlware

import (
	"backend/config"
	"backend/models"
	"errors"
	"fmt"
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// ErrorUnauthorized indicates that user is not authorized
var ErrorUnauthorized = errors.New("user is not authorized")

// ErrorInvalidUserData indicates that data of currently authorized user is invalid
var ErrorInvalidUserData = errors.New("user data is invalid")

type loginData struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// IdentityKey is a key for jwt identity
const IdentityKey = "id"

// User structure for payload
type User struct {
	Email string
}

// AuthMiddleware creates authMiddleware
func AuthMiddleware() *jwt.GinJWTMiddleware {
	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte(config.GetConfig().SecretKey),
		Timeout:     2 * time.Hour,
		MaxRefresh:  2 * time.Hour,
		IdentityKey: IdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					IdentityKey: v.Email,
				}
			}

			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				Email: claims[IdentityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var data loginData

			if err := c.ShouldBind(&data); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}

			email := data.Email
			password := data.Password
			user, err := models.GetUserByEmail(email)

			if err != nil {
				return nil, fmt.Errorf("could not get user: %w", err)
			}

			if user == nil {
				return nil, fmt.Errorf("user with email '%s' doest not exist", email)
			}

			if !user.IsActive() {
				return nil, fmt.Errorf("user with email '%s' is inactive", email)
			}

			if isPasswordValid := user.CheckPassword(password); !isPasswordValid {
				return nil, errors.New("invalid password")
			}

			pUser := &User{
				Email: user.Email,
			}

			return pUser, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*User); !ok {
				return false
			}

			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		log.Fatal(err)
	}

	return authMiddleware
}
