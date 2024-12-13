package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/gin-gonic/gin"
)

const (
	numUsers        = 1
	numSubreddits   = 1
	maxPostsPerSub  = 1
	simulationSteps = 1 // Number of simulation events
)

var (
	system         *actor.ActorSystem
	userActor      *actor.PID
	subredditActor *actor.PID
	postActor      *actor.PID
	users          map[string]*User

	subreddits map[string]*SubredditCreated
	totalUsers int
)

func main() {
	system = actor.NewActorSystem()

	// Create actors

	userProps := actor.PropsFromProducer(func() actor.Actor { return &UserActor{users: make(map[string]*User)} })
	userActor = system.Root.Spawn(userProps)

	subredditProps := actor.PropsFromProducer(func() actor.Actor { return &SubredditActor{subreddits: make(map[string]*Subreddit)} })
	subredditActor = system.Root.Spawn(subredditProps)

	postProps := actor.PropsFromProducer(func() actor.Actor { return &PostActor{posts: make(map[string]*Post)} })
	postActor = system.Root.Spawn(postProps)

	// commentProps := actor.PropsFromProducer(func() actor.Actor { return &CommentActor{comments: make(map[string]*Comment)} })
	// commentActor := system.Root.Spawn(commentProps)

	// dmProps := actor.PropsFromProducer(func() actor.Actor { return &DirectMessageActor{directMessages: []DirectMessage{}} })
	// dmActor := system.Root.Spawn(dmProps)

	users = make(map[string]*User)
	subreddits = make(map[string]*SubredditCreated)
	totalUsers = 0
	r := gin.Default()
	setupRoutes(r)

	//register users
	for i := 0; i < numUsers; i++ {
		respUser, _ := system.Root.RequestFuture(userActor, &RegisterUser{Username: fmt.Sprintf("user-%d", i)}, 1*time.Second).Result()
		user := respUser.(*User)
		users["t"] = user
	}
	//create subreddits
	for i := 0; i < numSubreddits; i++ {
		respSub, _ := system.Root.RequestFuture(subredditActor, &CreateSubreddit{Name: fmt.Sprintf("subreddit-%d", i)}, 1*time.Second).Result()
		sub := respSub.(*SubredditCreated)
		subreddits[""] = sub
	}

	// Assign users to subreddits based on Zipf distribution

	// distribution := generateZipfDistribution(numSubreddits, 1.1, 1)

	// for i, sub := range subreddits {
	// 	numMembers := distribution[i]
	// 	for j := 0; j < numMembers && j < numUsers; j++ {
	// 		system.Root.Send(subredditActor, &JoinSubreddit{UserID: users[j].ID, SubredditID: sub.ID})

	// 	}
	// }

	time.Sleep(1 * time.Second)

	// userConnectionManager := NewUserConnectionManager(users)
	// SimulateUserConnections(userConnectionManager, users, 10*time.Microsecond) // Toggle connection every 10 seconds

	// start := time.Now()

	// for step := 0; step < simulationSteps; step++ {
	// 	userIndex := rand.Intn(numUsers)
	// 	subIndex := rand.Intn(numSubreddits)

	// 	//randomize user activity
	// 	action := rand.Intn(6)
	// 	switch action {
	// 	case 0: // Create post
	// 		// check if user is part of the subreddit
	// 		if subreddits[subIndex].Users[users[userIndex].ID] {
	// 			SimulateSubredditActivity(system, postActor, subreddits[subIndex], users[userIndex], userConnectionManager)
	// 		}
	// 	case 1: // Upvote post
	// 		system.Root.Send(postActor, &UpvotePost{
	// 			UserID: users[userIndex].ID,
	// 			PostID: fmt.Sprintf("post-%d", rand.Intn(maxPostsPerSub)+1),
	// 		})
	// 	case 2: // Downvote post
	// 		system.Root.Send(postActor, &DownvotePost{
	// 			UserID: users[userIndex].ID,
	// 			PostID: fmt.Sprintf("post-%d", rand.Intn(maxPostsPerSub)+1),
	// 		})
	// 	case 3: // Leave subreddit
	// 		//check if user is part of subreddit
	// 		if subreddits[subIndex].Users[users[userIndex].ID] {
	// 			system.Root.Send(subredditActor, &LeaveSubreddit{
	// 				UserID:      users[userIndex].ID,
	// 				SubredditID: subreddits[subIndex].ID,
	// 			})
	// 		}
	// 	case 4: // Add comment
	// 		system.Root.Send(commentActor, &AddComment{
	// 			User:    users[userIndex],
	// 			PostID:  fmt.Sprintf("post-%d", rand.Intn(maxPostsPerSub)+1),
	// 			Parent:  nil,
	// 			Content: "This is a comment",
	// 		})
	// 	case 5: //Send Message
	// 		randomUser2 := rand.Intn(numUsers)
	// 		system.Root.Send(dmActor, &SendDirectMessage{
	// 			Sender:   users[userIndex],
	// 			Receiver: users[randomUser2],
	// 			Content:  "Hello, how are you?",
	// 		})
	// 	}
	// 	time.Sleep(100 * time.Microsecond)
	// }
	// elapsed := time.Since(start)

	// fmt.Printf("Execution time: %s\n", elapsed)
	r.Run(":8080")
}

func generateZipfDistribution(n int, s, v float64) []int {
	zipf := rand.NewZipf(rand.New(rand.NewSource(time.Now().Unix())), s, v, uint64(n-1))
	distribution := make([]int, n)
	for i := 0; i < n; i++ {
		distribution[i] = int(zipf.Uint64()) + 1 // Add 1 to avoid zeros
	}
	return distribution
}

type UserConnectionManager struct {
	connections map[string]bool
	mu          sync.Mutex
}

func NewUserConnectionManager(users []*User) *UserConnectionManager {
	connections := make(map[string]bool)
	for _, user := range users {
		connections[user.ID] = true // Initially, all users are connected
	}
	return &UserConnectionManager{connections: connections}
}

func (u *UserConnectionManager) ToggleConnection(userID string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.connections[userID] = !u.connections[userID]
}

func (u *UserConnectionManager) IsConnected(userID string) bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.connections[userID]
}

// Simulate user connection toggling
func SimulateUserConnections(manager *UserConnectionManager, users []*User, interval time.Duration) {
	for _, user := range users {
		manager.ToggleConnection(user.ID)
	}
	time.Sleep(interval)
	go func() {
		for {
			for _, user := range users {
				manager.ToggleConnection(user.ID)
			}
			time.Sleep(interval)
		}
	}()
}

// Increase post frequency for subreddits with more subscribers and simulate reposts
// func SimulateSubredditActivity(system *actor.ActorSystem, postActor *actor.PID, subreddit *Subreddit, user *User, userConnectionManager *UserConnectionManager) {
// 	numPosts := len(subreddit.Users) / 10 // Increase post frequency based on subscribers
// 	for i := 0; i < numPosts; i++ {
// 		// Choose a random connected user to make a post
// 		if userConnectionManager.IsConnected(user.ID) {
// 			if rand.Float64() < 0.2 { // 20% chance to repost
// 				// Select a random post from the subreddit
// 				var originalPost *Post
// 				for _, post := range subreddit.Posts {
// 					originalPost = post
// 					break
// 				}
// 				if originalPost != nil {
// 					content := "Repost: " + originalPost.Content
// 					system.Root.Send(postActor, &CreatePost{
// 						User:        user,
// 						SubredditID: subreddit.ID,
// 						Content:     content,
// 					})
// 				}
// 			} else {
// 				// Create a new post
// 				content := fmt.Sprintf("Post by %s in %s", user.Username, subreddit.Name)
// 				system.Root.Send(postActor, &CreatePost{
// 					User:        user,
// 					SubredditID: subreddit.ID,
// 					Content:     content,
// 				})
// 			}
// 		}
// 	}

// }
