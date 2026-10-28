package design

import (
	"math/rand"
	"slices"
	"sync"
	"testing"
)

func TestTwitter(t *testing.T) {
	tw := NewTwitter()
	tw.PostTweet(1, 5)
	if got := tw.GetNewsFeed(1); !slices.Equal(got, []int{5}) {
		t.Fatalf("feed = %v", got)
	}
	tw.Follow(1, 2)
	tw.PostTweet(2, 6)
	if got := tw.GetNewsFeed(1); !slices.Equal(got, []int{6, 5}) {
		t.Fatalf("feed = %v", got)
	}
	tw.Unfollow(1, 2)
	if got := tw.GetNewsFeed(1); !slices.Equal(got, []int{5}) {
		t.Fatalf("feed = %v", got)
	}
	tw.Unfollow(1, 1) // self: no-op
	tw.Follow(1, 1)
	if got := tw.GetNewsFeed(1); !slices.Equal(got, []int{5}) {
		t.Fatalf("self follow broke feed: %v", got)
	}
	if got := tw.GetNewsFeed(99); len(got) != 0 {
		t.Fatalf("unknown user feed = %v", got)
	}
}

func TestTwitterFeedLimit(t *testing.T) {
	tw := NewTwitter()
	for i := 1; i <= 25; i++ {
		tw.PostTweet(i%3, i) // users 0, 1, 2
	}
	tw.Follow(0, 1)
	tw.Follow(0, 2)
	want := []int{25, 24, 23, 22, 21, 20, 19, 18, 17, 16}
	if got := tw.GetNewsFeed(0); !slices.Equal(got, want) {
		t.Fatalf("feed = %v, want %v", got, want)
	}
}

// Model: keep every tweet globally; filter and sort on read.
func TestTwitterRandom(t *testing.T) {
	r := rand.New(rand.NewSource(9))
	tw := NewTwitter()
	type post struct{ user, id int }
	var all []post // in posting order
	follows := map[[2]int]bool{}
	for op := 0; op < 2000; op++ {
		a, b := r.Intn(5), r.Intn(5)
		switch r.Intn(4) {
		case 0:
			tw.PostTweet(a, op)
			all = append(all, post{a, op})
		case 1:
			tw.Follow(a, b)
			follows[[2]int{a, b}] = true
		case 2:
			tw.Unfollow(a, b)
			delete(follows, [2]int{a, b})
		default:
			var want []int
			for i := len(all) - 1; i >= 0 && len(want) < FeedSize; i-- {
				if p := all[i]; p.user == a || follows[[2]int{a, p.user}] {
					want = append(want, p.id)
				}
			}
			if got := tw.GetNewsFeed(a); !slices.Equal(got, want) && !(len(got) == 0 && len(want) == 0) {
				t.Fatalf("feed(%d) = %v, want %v", a, got, want)
			}
		}
	}
}

func TestTwitterConcurrent(t *testing.T) {
	tw := NewTwitter()
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func(u int) {
			defer wg.Done()
			for i := 0; i < 200; i++ {
				tw.PostTweet(u, u*1000+i)
				tw.Follow(u, (u+1)%8)
				_ = tw.GetNewsFeed(u)
				if i%10 == 0 {
					tw.Unfollow(u, (u+1)%8)
				}
			}
		}(g)
	}
	wg.Wait()
	if got := tw.GetNewsFeed(0); len(got) != FeedSize {
		t.Fatalf("feed length = %d", len(got))
	}
}
