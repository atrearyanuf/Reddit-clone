package main

import (
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

// Define Messages
type RegisterUser struct {
	Username string
	Password string
}
type UserRegistered struct {
	User *User
}

type CreateSubreddit struct {
	Name string
}
type SubredditCreated struct {
	Subreddit *Subreddit
}

type JoinSubreddit struct {
	UserName      string
	SubredditName string
}
type LeaveSubreddit struct {
	Username      string
	SubredditName string
}
type GetFeed struct {
	SubredditName string
}
type CreatePost struct {
	Username      string
	SubredditName string
	Content       string
}
type PostCreated struct {
	Post *Post
}

type UpvotePost struct {
	Username string
	PostID   string
}
type DownvotePost struct {
	Username string
	PostID   string
}

type AddComment struct {
	User    *User
	Content string
	PostID  string
	Parent  *Comment
}
type CommentAdded struct {
	Comment *Comment
}

type SendDirectMessage struct {
	Sender   *User
	Receiver *User
	Content  string
}
type GetDirectMessages struct {
	UserID string
}

type User struct {
	ID         string
	Username   string
	Karma      int
	Subreddits map[string]*Subreddit
}

type Subreddit struct {
	ID      string
	Name    string
	Users   map[string]bool
	Posts   map[string]*Post
	Created time.Time
}

type Post struct {
	ID            string
	Username      string
	SubredditName string
	Content       string
	Upvotes       int
	Downvotes     int
	Comments      map[string]*Comment
}

type Comment struct {
	ID        string
	User      *User
	Content   string
	PostID    string
	Parent    *Comment
	Upvotes   int
	Downvotes int
}

type UpvoteComment struct {
	UserID    string // ID of the user performing the upvote
	CommentID string // ID of the comment being upvoted
}

type DownvoteComment struct {
	UserID    string // ID of the user performing the downvote
	CommentID string // ID of the comment being downvoted
}

type DirectMessage struct {
	ID       string
	Sender   *User
	Receiver *User
	Content  string
	Read     bool
}

type DirectMessages struct {
	Messages []DirectMessage
}

type DirectMessageActor struct {
	directMessages []DirectMessage
}

// Actor Implementation
type UserActor struct {
	users map[string]*User
}

// Define Actor for Posts
type PostActor struct {
	posts map[string]*Post
}

type SubredditActor struct {
	subreddits map[string]*Subreddit
}

// Define Actor for Comments
type CommentActor struct {
	comments map[string]*Comment
}
type GetUser struct {
	Username string
}

type GetUserFeed struct {
	Username string
}

// RegisterUser
func (state *UserActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *RegisterUser:
		user := &User{
			ID:         fmt.Sprintf("user-%d", len(state.users)+1),
			Username:   msg.Username,
			Karma:      0,
			Subreddits: make(map[string]*Subreddit),
		}
		state.users[user.Username] = user
		fmt.Println("Created User : " + user.Username)
		ctx.Respond(user)
	case *GetUser:
		user := state.users[msg.Username]
		ctx.Respond(user)
	case *GetUserFeed:
		var subredditfeed []*Subreddit
		for _, value := range state.users[msg.Username].Subreddits {
			subredditfeed = append(subredditfeed, value)
		}
		ctx.Respond(subredditfeed)
	}

}

// CreateSubreddit
func (state *SubredditActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *CreateSubreddit:
		subreddit := &Subreddit{
			ID:      fmt.Sprintf("subreddit-%d", len(state.subreddits)+1),
			Name:    msg.Name,
			Users:   make(map[string]bool),
			Posts:   make(map[string]*Post),
			Created: time.Now(),
		}
		state.subreddits[subreddit.Name] = subreddit
		ctx.Respond(&SubredditCreated{Subreddit: subreddit})
		fmt.Println("Created subreddit :" + subreddit.Name)
	case *JoinSubreddit:
		state.subreddits[msg.SubredditName].Users[msg.UserName] = true
		user := users[msg.UserName]
		user.Subreddits[msg.SubredditName] = state.subreddits[msg.SubredditName]
		// user.Subreddits[state.subreddits[msg.SubredditName].Name] = state.subreddits[msg.SubredditName]
		fmt.Println("User " + msg.UserName + " joined subreddit " + state.subreddits[msg.SubredditName].Name)
		ctx.Respond(msg)
	case *LeaveSubreddit:
		delete(state.subreddits[msg.SubredditName].Users, msg.Username)
		user := users[msg.Username]
		delete(user.Subreddits, msg.SubredditName)
		fmt.Println("User " + msg.Username + " left subreddit " + state.subreddits[msg.SubredditName].Name)
		ctx.Respond(msg)
	case *GetFeed:
		var posts []*Post
		for _, value := range state.subreddits[msg.SubredditName].Posts {
			posts = append(posts, value)
			fmt.Println(value)
		}
		ctx.Respond(posts)
	}
}

type GetPost struct {
	PostID string
}

// CreatePost
func (state *PostActor) Receive(ctx actor.Context) {

	switch msg := ctx.Message().(type) {
	case *CreatePost:
		// Create a post
		post := &Post{
			ID:            fmt.Sprintf("post-%d", len(state.posts)+1),
			Content:       msg.Content,
			Username:      msg.Username,
			Upvotes:       0,
			Downvotes:     0,
			SubredditName: msg.SubredditName,
			Comments:      make(map[string]*Comment),
		}
		state.posts[post.ID] = post
		subreddits[msg.SubredditName].Subreddit.Posts[post.ID] = post
		fmt.Println("User " + msg.Username + " posted " + post.ID)
		ctx.Respond(&PostCreated{Post: post})

	case *UpvotePost:
		// Upvote a post
		if post, ok := state.posts[msg.PostID]; ok {
			post.Upvotes++
			user := users[post.Username]
			user.Karma++
			ctx.Respond(post)
		} else {
			ctx.Respond(&Post{})
		}
	case *DownvotePost:
		// Downvote a post
		if post, ok := state.posts[msg.PostID]; ok {
			post.Downvotes++
			user := users[post.Username]
			user.Karma--
			ctx.Respond(post)
		} else {
			ctx.Respond(&Post{})
		}
	case *GetPost:
		ctx.Respond(state.posts[msg.PostID])
	}
}

// CreateComment
func (state *CommentActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *AddComment:
		// Create a comment
		comment := &Comment{
			ID:        fmt.Sprintf("comment-%d", len(state.comments)+1),
			User:      msg.User,
			Content:   msg.Content,
			PostID:    msg.PostID,
			Parent:    msg.Parent,
			Upvotes:   0,
			Downvotes: 0,
		}
		state.comments[comment.ID] = comment
		ctx.Respond(&CommentAdded{Comment: comment})
	case *UpvoteComment:
		if comment, ok := state.comments[msg.CommentID]; ok {
			comment.Upvotes++
			comment.User.Karma++ // Increment Karma of comment owner
		}

	case *DownvoteComment:
		if comment, ok := state.comments[msg.CommentID]; ok {
			comment.Downvotes++
			comment.User.Karma-- // Decrement Karma of comment owner
		}
	}
}

// SendMessage
func (state *DirectMessageActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *SendDirectMessage:
		// Send a direct message
		message := DirectMessage{
			ID:       fmt.Sprintf("dm-%d", len(state.directMessages)+1),
			Sender:   msg.Sender,
			Receiver: msg.Receiver,
			Content:  msg.Content,
			Read:     false,
		}
		state.directMessages = append(state.directMessages, message)

	case *GetDirectMessages:
		// Get all direct messages for a user
		var messages []DirectMessage
		for _, dm := range state.directMessages {
			if dm.Receiver.ID == msg.UserID {
				messages = append(messages, dm)
			}
		}
		ctx.Respond(&DirectMessages{Messages: messages})
	}
}
