package keycloak

import (
	"context"
	"fmt"
	"slices"

	"github.com/Nerzal/gocloak/v13"
)

type KeyCloak struct {
	client *gocloak.GoCloak;
	ctx context.Context;
	token *gocloak.JWT;
}

type KeyCloakIface interface {
	LoginAdmin(username, password, realm string) (error);
	LogoutAllSessions(realm, userID string) error;
	CreateUser(realm string, user gocloak.User) (string, error);
	AddRealmRoleToUser(realm, userID string, roles []string) error;
	SetPassword(userID, realm, password string, temporary bool) error;
	DeleteUser(realm, userID string) error;
	GetUserByName(realm, username string) (*gocloak.User, error);
}

// Verify struct implements interface
var _ KeyCloakIface = &KeyCloak{}

func NewClient(ctx context.Context, basePath string, options ...func(*KeyCloak)) *KeyCloak {
	c := KeyCloak{
		client: gocloak.NewClient(basePath),
		ctx   : ctx,
	}
	return &c
}

func (k *KeyCloak) LoginAdmin(username, password, realm string) (error) {
	token, err := k.client.LoginAdmin(k.ctx, username, password, realm)

	if err != nil {
		return err
	}

	k.token = token
	return nil
}

func (k *KeyCloak) LogoutAllSessions(realm, userID string) error {
	err := k.client.LogoutAllSessions(
		k.ctx,
		k.token.AccessToken,
		realm,
		userID)

	if err == nil {
		k.token = nil
	}

	return err
}

func (k *KeyCloak) CreateUser(realm string, user gocloak.User) (token string, err error) {
	token, err = k.client.CreateUser(k.ctx, k.token.AccessToken, realm, user)
	return
}

func (k *KeyCloak) AddRealmRoleToUser(realm, userID string, roles []string) error {
	keycloakRoles, err := k.client.GetRealmRoles(
		k.ctx,
		k.token.AccessToken,
		realm,
		gocloak.GetRoleParams{})

	if err != nil {
		return err
	}

	var matchingRoles []gocloak.Role

	for _, role := range keycloakRoles {
		if slices.Contains(roles, *role.Name) {
			matchingRoles = append(matchingRoles, *role)
		}
	}

	if len(matchingRoles) != len(roles) {
		return fmt.Errorf("not all of the roles are present")
	}

	err = k.client.AddRealmRoleToUser(
		k.ctx,
		k.token.AccessToken,
		realm,
		userID,
		matchingRoles)

	return err
}

func (k *KeyCloak) SetPassword(userID, realm, password string, temporary bool) error {
	err := k.client.SetPassword(
		k.ctx,
		k.token.AccessToken,
		userID,
		realm,
		password,
		temporary)

	return err
}

func (k *KeyCloak) DeleteUser(realm, userID string) error {
	err := k.client.DeleteUser(k.ctx, k.token.AccessToken, realm, userID)
	return err
}

func (k *KeyCloak) GetUserByName(realm, username string) (*gocloak.User, error) {
	users, err := k.client.GetUsers(k.ctx, k.token.AccessToken, realm,
		gocloak.GetUsersParams{
			Username: &username,
		},
	)

	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	if len(users) > 1 {
		return nil, fmt.Errorf("multiple users found")
	}

	return users[0], err
}