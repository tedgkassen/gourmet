package lazy

import (
	"testing"
	"github.com/tedgkassen/gourmet/pkg/seq/eager"
)

func inc (i interface{}) interface{} {
	return i.(int) + 1
}

func dec (i interface{}) interface{} {
	return i.(int) - 1
}


func TestCycle(t *testing.T) {
	s := Cycle([]interface{}{1,2,3}...)
	d := Take(4, s)
	for _, v := range([]int{1,2,3,1}) {
		curr := <-d
		if v != curr {
			t.Fatalf("Cycle failed: Expected %d, got %d", v, curr)
		}
	}
}

func lessThan(n int) func(interface{})bool {
	return func(i interface{}) bool {
		return i.(int) < n
	}
}

func TestTake(t *testing.T) {
	s := Cycle([]interface{}{1,2,3}...)
	d := TakeEvery(2, s)
	f := eager.Collect(Take(3, d))
	for i, v := range([]int{1,3,2}) {
		if v != f[i] {
			t.Fatalf("Take failed: Expected %d, got %d", v, f[i])
		}
	}
	s = Seq(1,2,3,4,5)
	e := eager.Collect(TakeWhile(s, lessThan(4)))
	for i, v := range([]int{1,2,3}) {
		if v != e[i] {
			t.Fatalf("TakeWhile failed: Expected %d, got %d", v, e[i])
		}
	}
}

func TestIterate(t *testing.T) {
	inc := func(i interface{}) interface{} {return i.(int) + 1}
	s := Take(4, Iterate(1, inc))
	for _, v := range([]int{1,2,3,4}) {
		curr := <-s
		if v != curr {
			t.Fatalf("Iterate failed: Expected %d, got %d", v, curr)
		}
	}
}

func TestMap(t *testing.T) {
	s := eager.Collect(Map(inc, Seq(1,2,3,4)))
	for i, v := range([]int{2,3,4,5}) {
		if v != s[i] {
			t.Fatalf("Map failed: Expected %d, got %d", v, s[i])
		}
	}
}

func TestReduce(t *testing.T) {
	reducer := func(v interface{}, sum interface{}) interface{}{
		return v.(int) + sum.(int)
	}
	s := eager.Collect(Reduce(reducer, 0, Seq(1,2,3)))
	sum := s[len(s)-1]
	if sum != 6 {
		t.Fatalf("Reduce failed: Expected %d, got %d", 6, s[len(s)-1])
	}
}

func TestEach(t *testing.T) {
	s := Seq(1,2,3)
	r := []int{}
	e := func(i interface{}) {
		r = append(r, i.(int))
	}
	<-Each(e, s)
	for i, v := range([]int{1,2,3}) {
		if v != r[i] {
			t.Fatalf("Each failed: Expected %d, got %d", v, r[i])
		}
	}
}

func isEven(i interface{}) bool {
	return i.(int) % 2 == 0
}

func TestFilter(t *testing.T) {
	a := Seq(1,2,3,4)
	b := eager.Collect(Filter(isEven, a))
	for i, v := range([]int{2,4}) {
		if v != b[i] {
			t.Fatalf("Filter failed: Expected %d, got %d", v, b[i])
		}
	}
}

func TestInterleave(t *testing.T) {
	a := Seq(1,3,5)
	b := Seq(2,4,6)
	c := eager.Collect(Interleave(a,b))
	for i, v := range([]int{1,2,3,4,5,6}) {
		if v != c[i] {
			t.Fatalf("Interleave failed: Expected %d, got %d", v, c[i])
		}
	}
}

func TestFork(t *testing.T) {
	s := Seq(2,3,4,5)
	a, b := Fork(s)
	ar := eager.Collect(Map(inc, a))
	br := eager.Collect(Map(dec, b))
	for i, v := range([]int{2,3,4,5}) {
		if v+1 != ar[i] {
			t.Fatalf("Fork (possibly) failed: Expected %d from ar, got %d", v+1, ar[i])
		}
		if v-1 != br[i] {
			t.Fatalf("Fork (possibly) failed: Expected %d from br, got %d", v-1, br[i])
		}
	}
}
