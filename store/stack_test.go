package store

import (
	"fmt"
	"os"
	"testing"
)

func mks(s ...Store) []Store {
	return s
}

func makeStack(t *testing.T, stacks []Store, fail bool, flags ...bool) *StackedStore {
	st := &StackedStore{}
	st.Open(nil)
	var err error
	for i, sub := range stacks {
		bCNO := false
		bCO := false
		flagI := i << 1
		if flagI < len(flags) {
			bCNO = flags[flagI]
			if flagI+1 < len(flags) {
				bCO = flags[flagI+1]
			}
		}
		if err != nil {
			continue
		}
		if sub == nil {
			sub, _ = Open("memory://")
		}
		err = st.Push(sub, bCNO, bCO)
	}
	if err == nil {
		if fail {
			t.Errorf("Expected to fail creating stack, but didn't")
		}
		return st
	}
	if fail {
		t.Logf("Got expected error creating stack: %v", err)
	} else {
		t.Errorf("Expected to create stack, but failed with %v", err)
	}
	return nil
}

func checkErr(t *testing.T, ref, err error) {
	refType, errType := fmt.Sprintf("%T", ref), fmt.Sprintf("%T", err)
	if refType == errType {
		t.Logf("Got expected error %s: %v", refType, err)
	} else {
		t.Errorf("Unexpected error %s: %v", refType, err)
	}
}

func TestStackSimple(t *testing.T) {
	tobj := struct{ Foo, Bar string }{"foo", "bar"}
	var tgt interface{}
	s2, _ := Open("memory://")
	s2.Save("foo", &tobj)
	sub, _ := s2.MakeSub("sub1")
	s2.MakeSub("sub2")
	sub.Save("foo", &tobj)
	s3, _ := Open("memory://")
	s3.Save("bar", &tobj)
	s3.MakeSub("sub2")
	st := makeStack(t, mks(nil, s2, s3), false,
		false, true,
		true, true,
		false, false)
	if st == nil {
		return
	}
	checkErr(t, StackCannotOverride(""), st.Save("bar", &tobj))
	checkErr(t, StackCannotBeOverridden(""), st.Save("foo", &tobj))
	checkErr(t, nil, st.Save("baz", &tobj))
	checkErr(t, nil, st.Remove("baz"))
	checkErr(t, os.ErrNotExist, st.Remove("baz"))
	checkErr(t, UnWritable(""), st.Remove("foo"))
	checkErr(t, UnWritable(""), st.Remove("bar"))
	sub, err := st.MakeSub("sub1")
	checkErr(t, nil, err)
	checkErr(t, nil, sub.Load("foo", &tgt))
	checkErr(t, UnWritable(""), sub.Remove("foo"))
	_, err = st.MakeSub("sub3")
	checkErr(t, nil, err)
	st.Close()
}

func TestStackCannotBeOverridden(t *testing.T) {
	tobj := struct{ Foo, Bar string }{"foo", "bar"}
	s1, _ := Open("memory://")
	s1.Save("foo", &tobj)
	s1.MakeSub("sub1")
	s2, _ := Open("memory://")
	s2.Save("foo", &tobj)
	s2.MakeSub("sub1")
	st := makeStack(t, mks(s1, s2), true,
		false, false,
		true, false)
	if st != nil {
		t.Errorf("Expected stack creation to fail, it passed!")
	} else {
		t.Logf("Stack creation failed, as expected.")
	}
}

func TestSubStackCannotBeOverridden(t *testing.T) {
	tobj := struct{ Foo, Bar string }{"foo", "bar"}
	s1, _ := Open("memory://")
	sub1, _ := s1.MakeSub("sub1")
	sub1.Save("foo", &tobj)
	s2, _ := Open("memory://")
	sub2, _ := s2.MakeSub("sub1")
	sub2.Save("foo", &tobj)
	st := makeStack(t, mks(s1, s2), true,
		false, false,
		true, false)
	if st != nil {
		t.Errorf("Expected stack creation to fail, it passed!")
	} else {
		t.Logf("Stack creation failed, as expected.")
	}
}

func TestStackCannotOverride(t *testing.T) {
	tobj := struct{ Foo, Bar string }{"foo", "bar"}
	s1, _ := Open("memory://")
	s1.Save("foo", &tobj)
	s1.MakeSub("sub1")
	s2, _ := Open("memory://")
	s2.Save("foo", &tobj)
	s2.MakeSub("sub1")
	st := makeStack(t, mks(s1, s2), true,
		false, true,
		false, false)
	if st != nil {
		t.Errorf("Expected stack creation to fail, it passed!")
	} else {
		t.Logf("Stack creation failed, as expected.")
	}
}
func TestSubStackCannotOverride(t *testing.T) {
	tobj := struct{ Foo, Bar string }{"foo", "bar"}
	s1, _ := Open("memory://")
	sub1, _ := s1.MakeSub("sub1")
	sub1.Save("foo", &tobj)
	s2, _ := Open("memory://")
	sub2, _ := s2.MakeSub("sub1")
	sub2.Save("foo", &tobj)
	st := makeStack(t, mks(s1, s2), true,
		false, true,
		false, false)
	if st != nil {
		t.Errorf("Expected stack creation to fail, it passed!")
	} else {
		t.Logf("Stack creation failed, as expected.")
	}
}
