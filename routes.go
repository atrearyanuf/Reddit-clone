package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func setupRoutes(r *gin.Engine) {

	r.POST("/user", registerNewUser)
	r.GET("/user/:username", getUser)
	r.POST("/subreddit", createSubreddit)
	r.POST("/subreddit/join", joinSubreddit)
	r.POST("/subreddit/leave", leaveSubreddit)
	r.POST("/subreddit/post", createPost)
	r.POST("/subreddit/post/upvote", upvotePost)
	r.POST("/subreddit/post/downvote", downvotePost)
	r.GET("/subreddit/:subreddit", getFeed)
	r.GET("/user/feed/:username", getUserFeed)
	//add comment
	//messages
}

func registerNewUser(c *gin.Context) {
	var req RegisterUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(req)
	// Use actor system to handle registration
	respUser, err := system.Root.RequestFuture(userActor, &RegisterUser{Username: req.Username, Password: req.Password}, 2*time.Second).Result()
	user := respUser.(*User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
	users[req.Username] = user
	totalUsers++
}

func getUser(c *gin.Context) {
	username := c.Param("username")

	// Use actor system to handle registration
	respUser, err := system.Root.RequestFuture(userActor, &GetUser{Username: username}, 2*time.Second).Result()
	user := respUser.(*User)
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User Not Found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func createSubreddit(c *gin.Context) {
	var req CreateSubreddit
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use actor system to handle registration
	respSub, err := system.Root.RequestFuture(subredditActor, &CreateSubreddit{Name: req.Name}, 1*time.Second).Result()
	sub := respSub.(*SubredditCreated)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sub)
	subreddits[req.Name] = sub
}

func getFeed(c *gin.Context) {
	subreddit := c.Param("subreddit")

	// Use actor system to handle registration
	respSub, err := system.Root.RequestFuture(subredditActor, &GetFeed{SubredditName: subreddit}, 1*time.Second).Result()
	posts := respSub.([]*Post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func getUserFeed(c *gin.Context) {
	username := c.Param("username")

	// Use actor system to handle registration
	respSub, err := system.Root.RequestFuture(userActor, &GetUserFeed{Username: username}, 1*time.Second).Result()
	subredditfeed := respSub.([]*Subreddit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subredditfeed)
}
func joinSubreddit(c *gin.Context) {
	var req JoinSubreddit
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Use actor system to handle registration
	_, err := system.Root.RequestFuture(subredditActor, &JoinSubreddit{UserName: req.UserName, SubredditName: req.SubredditName}, 1*time.Second).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
}

func leaveSubreddit(c *gin.Context) {
	var req LeaveSubreddit
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Use actor system to handle registration
	_, err := system.Root.RequestFuture(subredditActor, &LeaveSubreddit{Username: req.Username, SubredditName: req.SubredditName}, 1*time.Second).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
}

func createPost(c *gin.Context) {
	var req CreatePost
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use actor system to handle post creation
	respPost, err := system.Root.RequestFuture(postActor, &CreatePost{
		Username:      req.Username,
		SubredditName: req.SubredditName,
		Content:       req.Content,
	}, 2*time.Second).Result()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	post := respPost.(*PostCreated)
	c.JSON(http.StatusOK, post)
}

func upvotePost(c *gin.Context) {
	var req UpvotePost
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use actor system to handle post creation
	respPost, err := system.Root.RequestFuture(postActor, &UpvotePost{
		Username: req.Username,
		PostID:   req.PostID,
	}, 2*time.Second).Result()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	post := respPost.(*Post)
	c.JSON(http.StatusOK, post)
}

func downvotePost(c *gin.Context) {
	var req DownvotePost
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use actor system to handle post creation
	respPost, err := system.Root.RequestFuture(postActor, &DownvotePost{
		Username: req.Username,
		PostID:   req.PostID,
	}, 2*time.Second).Result()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	post := respPost.(*Post)
	c.JSON(http.StatusOK, post)
}
