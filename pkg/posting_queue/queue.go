package posting_queue

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/goriiin/kotyari-bots_backend/internal/model"
	"github.com/goriiin/kotyari-bots_backend/pkg/otvet"
)

// QueuedPost represents a post waiting to be published
type QueuedPost struct {
	ID                 uuid.UUID
	Post               *model.Post
	Candidate          model.Candidate
	Request            PostRequest
	RequiresModeration bool
	Approved           bool
	CreatedAt          time.Time
}

// PostRequest contains information needed to publish a post
type PostRequest struct {
	Platform  model.PlatformType
	PostType  model.PostType
	TopicType int
	Spaces    []otvet.Space
	// ModerationRequired allows per-request override from bot configuration
	ModerationRequired bool
}

// Account represents an account for posting
type Account struct {
	ID        string
	AuthToken string
	Client    *otvet.OtvetClient
	LastPost  time.Time
}

// Queue is an in-memory queue for posts
type Queue struct {
	mu                 sync.RWMutex
	posts              []*QueuedPost
	accounts           map[string]*Account
	postingInterval    time.Duration
	processingInterval time.Duration
	moderationRequired bool
	stopChan           chan struct{}
}

// NewQueue creates a new posting queue
func NewQueue(postingInterval, processingInterval time.Duration, moderationRequired bool) *Queue {
	return &Queue{
		posts:              make([]*QueuedPost, 0),
		accounts:           make(map[string]*Account),
		postingInterval:    postingInterval,
		processingInterval: processingInterval,
		moderationRequired: moderationRequired,
		stopChan:           make(chan struct{}),
	}
}

// AddAccount adds an account to the queue
func (q *Queue) AddAccount(id string, authToken string, client *otvet.OtvetClient) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.accounts[id] = &Account{
		ID:        id,
		AuthToken: authToken,
		Client:    client,
		LastPost:  time.Time{},
	}
}

// Enqueue adds a post to the queue
func (q *Queue) Enqueue(post *model.Post, candidate model.Candidate, req PostRequest) *QueuedPost {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Determine if moderation is required either by queue default or by per-request flag
	requiresModeration := q.moderationRequired || req.ModerationRequired

	queuedPost := &QueuedPost{
		ID:                 uuid.New(),
		Post:               post,
		Candidate:          candidate,
		Request:            req,
		RequiresModeration: requiresModeration,
		Approved:           !requiresModeration,
		CreatedAt:          time.Now(),
	}

	q.posts = append(q.posts, queuedPost)
	return queuedPost
}

// ApprovePost approves a post for publishing by post ID from database
func (q *Queue) ApprovePost(postID uuid.UUID) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	for _, post := range q.posts {
		if post.Post != nil && post.Post.ID == postID {
			post.Approved = true
			return nil
		}
	}
	return ErrPostNotFound
}

// GetPostByID returns a post by post ID from database
func (q *Queue) GetPostByID(postID uuid.UUID) (*QueuedPost, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	for _, post := range q.posts {
		if post.Post != nil && post.Post.ID == postID {
			return post, nil
		}
	}
	return nil, ErrPostNotFound
}

// StartProcessing starts the queue processing loop
func (q *Queue) StartProcessing(ctx context.Context, publishFunc func(ctx context.Context, account *Account, post *QueuedPost) error) {
	ticker := time.NewTicker(q.processingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-q.stopChan:
			return
		case <-ticker.C:
			q.processQueue(ctx, publishFunc)
		}
	}
}

// Stop stops the queue processing
func (q *Queue) Stop() {
	close(q.stopChan)
}

// processQueue processes the queue and publishes approved posts
func (q *Queue) processQueue(ctx context.Context, publishFunc func(ctx context.Context, account *Account, post *QueuedPost) error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.accounts) == 0 {
		return
	}

	// Get random account
	account := q.getRandomAccount()
	if account == nil {
		return
	}

	// Check if account can post (respecting posting interval)
	if !q.canPost(account) {
		return
	}

	// Find first approved post (either doesn't require moderation or is approved)
	var postToPublish *QueuedPost
	postIndex := -1

	for i, post := range q.posts {
		if post.Approved {
			postToPublish = post
			postIndex = i
			break
		}
	}

	if postToPublish == nil {
		return
	}

	// Publish post
	if err := publishFunc(ctx, account, postToPublish); err != nil {
		return
	}

	// Update account last post time
	account.LastPost = time.Now()

	// Remove post from queue
	q.posts = append(q.posts[:postIndex], q.posts[postIndex+1:]...)
}

// getRandomAccount returns a random account from the map
func (q *Queue) getRandomAccount() *Account {
	if len(q.accounts) == 0 {
		return nil
	}

	accounts := make([]*Account, 0, len(q.accounts))
	for _, acc := range q.accounts {
		accounts = append(accounts, acc)
	}

	return accounts[rand.Intn(len(accounts))]
}

// canPost checks if account can post (respecting posting interval)
func (q *Queue) canPost(account *Account) bool {
	if account.LastPost.IsZero() {
		return true
	}
	return time.Since(account.LastPost) >= q.postingInterval
}

// GetQueueSize returns the current queue size
func (q *Queue) GetQueueSize() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.posts)
}

// GetPendingModerationCount returns count of posts waiting for moderation
func (q *Queue) GetPendingModerationCount() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	count := 0
	for _, post := range q.posts {
		if post.RequiresModeration && !post.Approved {
			count++
		}
	}
	return count
}
