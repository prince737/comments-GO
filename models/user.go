package models

import (
	"errors"
	"fmt"

	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
)

var (
	//ErrUserNotFound user not registered
	ErrUserNotFound = errors.New("User not found")
	//ErrInvalidLogin incorrect password or username
	ErrInvalidLogin = errors.New("Invalid login")
	//ErrUsernameTaken user already exists
	ErrUsernameTaken = errors.New("Username taken")
)

//User identifies a user
type User struct {
	id int64
}

//NewUser user constructor
func newUser(username string, hash []byte) (*User, error) {
	exists, err := client.HExists("user:by-username", username).Result()
	fmt.Println("here")
	if exists {
		fmt.Println("here")
		return nil, ErrUsernameTaken
	}
	id, err := client.Incr("user:next-id").Result()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	key := fmt.Sprintf("user:%d", id)
	pipe := client.Pipeline()
	pipe.HSet(key, "id", id)
	pipe.HSet(key, "username", username)
	pipe.HSet(key, "hash", hash)
	pipe.HSet("user:by-username", username, id)
	_, err = pipe.Exec()
	if err != nil {
		return nil, err
	}
	return &User{id}, nil
}

//GetID from pipeline using key
func (user *User) GetID() (int64, error) {
	return user.id, nil
}

//GetUsername from pipeline using key
func (user *User) GetUsername() (string, error) {
	key := fmt.Sprintf("user:%d", user.id)
	return client.HGet(key, "username").Result()
}

//GetHash from pipeline using key
func (user *User) GetHash() ([]byte, error) {
	key := fmt.Sprintf("user:%d", user.id)
	return client.HGet(key, "hash").Bytes()
}

//Authenticate authenticates a user
func (user *User) authenticate(password string) error {
	hash, err := user.GetHash()
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return ErrInvalidLogin
	}
	return err
}

//GetUserByID gets a user from userid
func GetUserByID(id int64) (*User, error) {
	return &User{id}, nil
}

//GetUserByUsername gets a user by username
func GetUserByUsername(username string) (*User, error) {
	id, err := client.HGet("user:by-username", username).Int64()
	if err == redis.Nil {
		return nil, ErrUserNotFound
	} else if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return GetUserByID(id)

}

//AuthenticateUser autheoticates a user
func AuthenticateUser(username, password string) (*User, error) {
	user, err := GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	return user, user.authenticate(password)
}

//RegisterUser registers a user
func RegisterUser(username, password string) error {
	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return err
	}
	_, err = newUser(username, hash)
	return err
}
