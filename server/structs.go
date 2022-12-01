package server

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// JSONData is the main struct for posts, it holds all the necessary data about the user, the post and the comments
type JSONData struct {
	Post_id       int                  `json:"post_id"`
	User_id       int                  `json:"user_id"`
	Heading       string               `json:"heading"`
	Body          string               `json:"body"`
	Closed_user   int                  `json:"closed_user"`
	Closed_admin  int                  `json:"closed_admin"`
	Closed_date   string               `json:"closed_date"`
	Insert_time   string               `json:"insert_time"`
	Update_time   string               `json:"update_time"`
	Image         string               `json:"image"`
	Comments      map[int]JSONComments `json:"comments"`
	Categories    []string             `json:"categories"`
	Reactions     []map[int]string     `json:"reactions"`
	UserReaction  string               `json:"user_reaction"`
	Username      string               `json:"username"`
	Profile_image string               `json:"profile_image"`
}

// JSONComments holds the necessary data about comments
type JSONComments struct {
	CommentID     int              `json:"comment_id"`
	Post_id       int              `json:"post_id"`
	User_id       int              `json:"user_id"`
	Body          string           `json:"body"`
	Insert_time   string           `json:"insert_time"`
	Reactions     []map[int]string `json:"reactions"`
	UserReaction  string           `json:"user_reaction"`
	Username      string           `json:"username"`
	Profile_image string           `json:"profile_image"`
}

type newPosts struct {
	Post_id    int      `json:"post_id"`
	Body       string   `json:"postBody"`
	Heading    string   `json:"postHeading"`
	Categories []string `json:"postCats"`
}

type LastPost struct {
	ID int `json:"lastPostID"`
}

type Comments struct {
	Body    string `json:"postComment"`
	Post_id int    `json:"postID"`
}

type SessionManager struct {
	sessions map[string]*SessionData
}

type SessionData struct {
	UId  string
	Misc map[string]interface{}
}

// IPRateLimiter is a struct that holds the connected IP addresses, a RWMutex, a rate.Limit on requests/second and a limit on package bursts
type IPRateLimiter struct {
	ips         map[string]*rate.Limiter
	lastCleaned time.Time
	mu          *sync.RWMutex
	r           rate.Limit
	b           int
}

// NewIPRateLimiter returns a pointer to a new IPRateLimiter object
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	i := &IPRateLimiter{
		ips:         make(map[string]*rate.Limiter),
		lastCleaned: time.Now(),
		mu:          &sync.RWMutex{},
		r:           r,
		b:           b,
	}
	return i
}

type Info struct {
	Image    string `json:"Image"`
	Username string `json:"Username"`
}

// ErrorMessage holds the basic error codes for the ErrorPage function
type ErrorMessage struct {
	ErrorCode int
	Message   string
}

type JSONCategories struct {
	Categories string `json:"categories"`
}

type NewUser struct {
	User_id    int    `json:"user_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	ConfirmPwd string `json:"confirmPassword"`
}
