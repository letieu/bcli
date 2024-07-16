package view

import (
	"bcli/api"
	"fmt"
)

func PrintUser(userInfo *api.UserInfo) {
    fmt.Println("Current user information:")
    fmt.Printf(" - User ID: %s\n", userInfo.UsrID)
    fmt.Printf(" - User name: %s\n", userInfo.UsrNm)
    fmt.Printf(" - User email: %s\n", userInfo.UsrEml)
}
