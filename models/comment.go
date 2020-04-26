package models

import (
	"fmt"
	"strconv"
)

//Comment identifies a comment
type Comment struct {
	id int64
}

//newComment comment constructor
func newComment(userID int64, comment string) (*Comment, error) {
	id, err := client.Incr("comment:next-id").Result()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	key := fmt.Sprintf("comment:%d", id)
	pipe := client.Pipeline()
	pipe.HSet(key, "id", id)
	pipe.HSet(key, "body", comment)
	pipe.HSet(key, "user_id", userID)
	pipe.LPush("comments", id)
	pipe.LPush(fmt.Sprintf("user:%d:comments", userID), id)
	_, err = pipe.Exec()
	if err != nil {
		return nil, err
	}
	return &Comment{id}, nil
}

//GetBody returns the body of a commment
func (comment *Comment) GetBody() (string, error) {
	key := fmt.Sprintf("comment:%d", comment.id)
	return client.HGet(key, "body").Result()
}

//GetUser returns the user object that posted the comment
func (comment *Comment) GetUser() (*User, error) {
	key := fmt.Sprintf("comment:%d", comment.id)
	userID, err := client.HGet(key, "user_id").Int64()
	if err != nil {
		return nil, err
	}
	return GetUserByID(userID)
}

func queryComments(key string) ([]*Comment, error) {
	commentIDs, err := client.LRange(key, 0, 10).Result()
	if err != nil {
		return nil, err
	}
	comments := make([]*Comment, len(commentIDs))
	for i, strID := range commentIDs {
		id, err := strconv.Atoi(strID)
		if err != nil {
			return nil, err
		}
		comments[i] = &Comment{int64(id)}
	}
	return comments, nil
}

//GetAllComments Gets all comments from redis db
func GetAllComments() ([]*Comment, error) {
	return queryComments("comments")
}

//GetComments Gets all comments from redis db
func GetComments(userID int64) ([]*Comment, error) {
	key := fmt.Sprintf("user:%d:comments", userID)
	return queryComments(key)
}

//PostComment posts a comment to the redis db
func PostComment(userID int64, comment string) error {
	_, err := newComment(userID, comment)
	return err
}
