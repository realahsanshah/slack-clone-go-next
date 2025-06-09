package auth

import (
	"fmt"
	"slack-clone-go-next/utils"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	utils.Created(c, gin.H{"message": "Register"}, "Success")
}

func Login(c *gin.Context) {
	fmt.Println("Login")
	utils.OK(c, gin.H{"message": "Login"}, "Success")
}
