package tries

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestSuggestedProducts(t *testing.T) {
	tests := []struct {
		name     string
		products []string
		search   string
		want     [][]string
	}{
		{
			"classic",
			[]string{"mobile", "mouse", "moneypot", "monitor", "mousepad"},
			"mouse",
			[][]string{
				{"mobile", "moneypot", "monitor"},
				{"mobile", "moneypot", "monitor"},
				{"mouse", "mousepad"},
				{"mouse", "mousepad"},
				{"mouse", "mousepad"},
			},
		},
		{"single product", []string{"havana"}, "havana", [][]string{
			{"havana"}, {"havana"}, {"havana"}, {"havana"}, {"havana"}, {"havana"},
		}},
		{"falls off early", []string{"bags", "baggage", "banner", "box", "cloths"}, "bz", [][]string{
			{"baggage", "bags", "banner"}, nil,
		}},
		{"no match at all", []string{"abc"}, "x", [][]string{nil}},
		{"empty search", []string{"abc"}, "", [][]string{}},
	}
	for _, tt := range tests {
		for name, fn := range map[string]func([]string, string) [][]string{
			"trie": SuggestedProducts, "sorted": SuggestedProductsSorted,
		} {
			if got := fn(tt.products, tt.search); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s/%s = %q, want %q", tt.name, name, got, tt.want)
			}
		}
	}
}

func TestSuggestedProductsRandom(t *testing.T) {
	r := rand.New(rand.NewSource(15))
	word := func() string {
		b := make([]byte, 1+r.Intn(4))
		for i := range b {
			b[i] = byte('a' + r.Intn(3))
		}
		return string(b)
	}
	for iter := 0; iter < 300; iter++ {
		seen := map[string]bool{}
		var products []string
		for k := r.Intn(8); k > 0; k-- {
			if w := word(); !seen[w] { // the problem promises distinct products
				seen[w] = true
				products = append(products, w)
			}
		}
		search := word()
		got, want := SuggestedProducts(products, search), SuggestedProductsSorted(products, search)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("products %q search %q: trie=%q sorted=%q", products, search, got, want)
		}
	}
}
