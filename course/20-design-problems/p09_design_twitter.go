package design

import (
	"container/heap"
	"sync"
)

// Problem: Design Twitter (simplified)
// Difficulty: Medium
//
// Users post tweets and follow each other. GetNewsFeed(user) returns the
// ids of the 10 most recent tweets posted by the user or anyone they follow,
// newest first. Following yourself is implicit; Unfollow of someone you
// don't follow is a no-op.
//
// Example:
//
//	PostTweet(1, 5); GetNewsFeed(1) -> [5]
//	Follow(1, 2); PostTweet(2, 6); GetNewsFeed(1) -> [6 5]
//	Unfollow(1, 2); GetNewsFeed(1) -> [5]
//
// Reference: https://leetcode.com/problems/design-twitter/
//
// Pattern: k-way merge with a heap ("fan-out on read"). Each user's tweets
// are already sorted by time (appended in order). To build a feed, push the
// newest tweet of each followee into a max-heap keyed by time, pop the top,
// then push that user's next-older tweet. Stop after 10 pops.
//
//	alice: t1 t4 t9        heap holds one cursor per user
//	bob:   t2 t7      ->   pop t9 (alice) push t4; pop t7 (bob) push t2; …
//
// A global counter supplies timestamps so ordering is total.
//
// Time: Post/Follow/Unfollow O(1); GetNewsFeed O(f + k log f) for f
// followees and k = 10. Space: O(users + tweets + follows).
//
// Safe for concurrent use. GetNewsFeed only reads, so it takes a read lock
// and feeds can be built in parallel.
//
// For the fan-out-on-write alternative (precomputed timelines) and its
// trade-offs see:
// https://github.com/neerajsinghi/system-design-patterns-go (patterns/architecture/fan-out-feed)

// FeedSize is how many tweets GetNewsFeed returns at most.
const FeedSize = 10

type tweet struct {
	id, time int
}

// Twitter is a tiny in-memory social network.
type Twitter struct {
	mu        sync.RWMutex
	clock     int
	tweets    map[int][]tweet          // user → tweets, oldest first
	following map[int]map[int]struct{} // follower → set of followees
}

// NewTwitter returns an empty network.
func NewTwitter() *Twitter {
	return &Twitter{
		tweets:    make(map[int][]tweet),
		following: make(map[int]map[int]struct{}),
	}
}

// PostTweet records that user posted tweetID now.
func (t *Twitter) PostTweet(user, tweetID int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.clock++
	t.tweets[user] = append(t.tweets[user], tweet{tweetID, t.clock})
}

// Follow makes follower see followee's tweets.
func (t *Twitter) Follow(follower, followee int) {
	if follower == followee {
		return // self is always included
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	set := t.following[follower]
	if set == nil {
		set = make(map[int]struct{})
		t.following[follower] = set
	}
	set[followee] = struct{}{}
}

// Unfollow stops follower seeing followee's tweets.
func (t *Twitter) Unfollow(follower, followee int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.following[follower], followee) // delete on a nil map is a no-op
}

// feedCursor points at one user's tweet list, at index pos.
type feedCursor struct {
	list []tweet
	pos  int
}

// cursorHeap is a max-heap by the time of each cursor's current tweet.
type cursorHeap []feedCursor

func (h cursorHeap) Len() int { return len(h) }
func (h cursorHeap) Less(i, j int) bool {
	return h[i].list[h[i].pos].time > h[j].list[h[j].pos].time
}
func (h cursorHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *cursorHeap) Push(x any)   { *h = append(*h, x.(feedCursor)) }
func (h *cursorHeap) Pop() any {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

// GetNewsFeed returns up to 10 tweet ids, newest first.
func (t *Twitter) GetNewsFeed(user int) []int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	h := cursorHeap{}
	addUser := func(u int) {
		if list := t.tweets[u]; len(list) > 0 {
			h = append(h, feedCursor{list, len(list) - 1}) // start at newest
		}
	}
	addUser(user)
	for followee := range t.following[user] {
		addUser(followee)
	}
	heap.Init(&h) // O(f) heapify instead of f pushes

	feed := make([]int, 0, FeedSize)
	for h.Len() > 0 && len(feed) < FeedSize {
		top := &h[0]
		feed = append(feed, top.list[top.pos].id)
		if top.pos == 0 {
			heap.Pop(&h) // this user has no older tweets
		} else {
			top.pos--
			heap.Fix(&h, 0) // its key got smaller: sift down
		}
	}
	return feed
}
