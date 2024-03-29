package types

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

type UpdateUserParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (p UpdateUserParams) ToBSON() bson.M {
	m := bson.M{}
	if len(p.Username) > 0 {
		m["username"] = p.Username
	}
	if len(p.Password) > 0 {
		m["password"] = p.Password
	}
	return m
}

type CreateUserParams struct {
	Username     string               `json:"Username"`
	Password     string               `json:"password"`
	Description  string               `json:"description"`
	Email        string               `json:"email"`
	OwnedServers []primitive.ObjectID `json:"ownedServers"`
	/*
		TODO: need to add (status(enum), []friends, []servers, []blocked, []activities, nitro boolean,[]badges, OwnedServers []Server.
	*/
}

func (params CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.Username) < minUsernameLen {
		errors["username"] = fmt.Sprintf("username length should be at least %d characters", minUsernameLen)
	}
	if len(params.Description) > descriptionLen {
		errors["description"] = fmt.Sprintf("description length should be at least %d characters", descriptionLen)
	}
	if len(params.Password) < minPasswordLen {
		errors["password"] = fmt.Sprintf("password length should be at least %d characters", minPasswordLen)
	}
	if !isEmailValid(params.Email) {
		errors["email"] = fmt.Sprintf("email is invalid")
	}
	return errors
}

func isEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}
func IsValidPassword(encpw, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encpw), []byte(pw)) == nil
}

type User struct {
	ID                primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Username          string               `bson:"username" json:"username"`
	EncryptedPassword string               `bson:"EncryptedPassword" json:"-"`
	Description       string               `bson:"description" json:"description"`
	Email             string               `bson:"email" json:"email"`
	OwnedServers      []primitive.ObjectID `bson:"ownedServers" json:"ownedServers"` // identification of owned servers
}

func NewUser(params CreateUserParams) (*User, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), passwordEncryptionLevel)
	if err != nil {
		return nil, err
	}
	return &User{
		Username:          params.Username,
		Description:       params.Description,
		Email:             params.Email,
		EncryptedPassword: string(encryptedPassword),
		OwnedServers:      params.OwnedServers,
	}, nil
}
