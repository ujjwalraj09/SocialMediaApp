package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string `json:"name"`
	Email        string `json:"email" gorm:"unique"`
	Password     string `json:"password"`
	Posts        []Post `gorm:"foreignKey:UserID"` // One-to-many relationship with Post model
	ProfileImage string `json:"profile_image"`
}

type Post struct {
	gorm.Model
	UserID      string      `json:"user_id" gorm:"foreignKey:UserID"`
	PostID      string      `gorm:"unique"`
	PostContent string      `json:"post_content"`
	LikeCount   uint        `json:"like_count"`
	PostImages  []PostImage `gorm:"foreignKey:PostID"`
	Comments    []Comment   `gorm:"foreignKey:PostID"`
	PostTags    []Tags      `gorm:"many2many:post_tags;foreignKey:PostID"` // Many-to-many relationship with Tags
	Likes       []Like      `gorm:"foreignKey:PostID"`
}

type PostImage struct {
	gorm.Model
	PostID  string `json:"post_id"`
	ImageID string `json:"image_id"`
	Image   []byte `json:"image_data"`
}

type Comment struct {
	gorm.Model

	PostID      string `json:"post_id" gorm:"foreignKey:PostID"`
	CommentID   string `json:"comment_id"`
	CommentText string `json:"comment_text"`
}

type Like struct {
	gorm.Model
	UserID string `json:"user_id" gorm:"foreignKey:UserID"` //ID of the user who is liking the post
	PostID string `json:"post_id" gorm:"foreignKey:PostID"` //PostID to be Liked
}

type Tags struct {
	gorm.Model
	Tag    string `json:"tag"`
	PostID string
}

type BookMark struct {
	gorm.Model
	UserID    string `gorm:"primaryKey"`
	UserPosts []Post `gorm:"foreignKey:UserID"`
}

type Follow struct {
	gorm.Model
	FollowerID string `json:"follower_id"`
	FollowedID string `json:"followed_id"`
}
