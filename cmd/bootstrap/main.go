package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/egandro/keycloak-demo/pkg/keycloak"
	"github.com/egandro/keycloak-demo/pkg/tools"
)


func main() {
	ctx := context.Background()

	url:= "http://localhost:8080"

	realm := "master"
	bootstrapAdmin := "admin"
	bootstrapAdminPassword := "admin"

	admin := "CoolGuy"
	adminPassword := "Secret!123"
	adminGroup := "admin"

	err := tools.WaitForURL(url, 10, 10*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}

	k := keycloak.NewClient(ctx, url)

	err = k.LoginAdmin(bootstrapAdmin,bootstrapAdminPassword, realm)
	if err != nil {
		fmt.Println(err)
		return
	}

	user := gocloak.User{
		FirstName: gocloak.StringP("Bob"),
		LastName:  gocloak.StringP("Uncle"),
		Email:     gocloak.StringP("something@really.wrong"),
		Enabled:   gocloak.BoolP(true),
		Username:  gocloak.StringP(admin),
		EmailVerified: gocloak.BoolP(true),
	}

	userID, err := k.CreateUser(realm, user)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Created admin user:", userID)

	err = k.AddRealmRoleToUser(realm, userID, []string{adminGroup})
	if err != nil {
		fmt.Println(err)
		return
	}

	err = k.SetPassword(userID, realm, adminPassword, false)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = k.LogoutAllSessions(realm, userID)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("logged out all sessions for bootstrap admin:", bootstrapAdmin)

	err = k.LoginAdmin(admin, adminPassword, realm)
	if err != nil {
		fmt.Println(err)
		return
	}

	bootstrapAdminUser, err := k.GetUserByName(realm, bootstrapAdmin)
	if err != nil {
		fmt.Println(err)
		return
	}

	if bootstrapAdminUser == nil {
		fmt.Println("Bootstrap Admin User not found")
		return
	}

	err = k.DeleteUser(realm, *bootstrapAdminUser.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Bootstrap Admin User removed successfully - new admin is: ", admin)

	err = k.LogoutAllSessions(realm, userID)
	if err != nil {
		fmt.Println(err)
		return
	}
}
