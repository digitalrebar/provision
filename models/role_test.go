package models

import (
	"log"
	"testing"
)

func twoClaims(frags ...string) []*Claim {
	if len(frags) != 6 {
		panic("Need 6 args")
	}
	log.Printf("Making %v and %v", frags[0:3], frags[3:6])
	return makeClaims(frags...)
}

func claimContains(t *testing.T, c []*Claim) {
	t.Helper()
	a, b := c[0], c[1]
	if a.Contains(b) {
		t.Logf("Claim '%s' contains '%s'", a, b)
	} else if b.Contains(a) {
		t.Errorf("ERROR: Claim '%s' does not contain '%s'", a, b)
	} else {
		t.Errorf("ERROR: Claims '%s' and '%s' are disjoint", a, b)
	}
}

func claimDoesNotContain(t *testing.T, c []*Claim) {
	t.Helper()
	a, b := c[0], c[1]
	if a.Contains(b) {
		t.Errorf("ERROR Claim '%s' contains '%s'", a, b)
	} else {
		t.Logf("Claim '%s' does not contain '%s'", a, b)
	}
}

func claimsAreDisjoint(t *testing.T, c []*Claim) {
	t.Helper()
	a, b := c[0], c[1]
	if a.Contains(b) {
		t.Errorf("ERROR: Claim '%s' contains '%s'", a, b)
	} else if b.Contains(a) {
		t.Errorf("ERROR: Claim '%s' contains '%s'", b, a)
	} else {
		t.Logf("Claims '%s' and '%s' are disjoint", a, b)
	}
}

func claimsAreOrdered(t *testing.T, c []*Claim) {
	t.Helper()
	a, b := c[0], c[1]
	if a.Contains(b) && !b.Contains(a) {
		t.Logf("Claims '%s' and '%s' are ordered", a, b)
	} else {
		t.Errorf("ERROR: Claims '%s' and '%s' are not ordered", a, b)
	}
}

func claimsAreEqual(t *testing.T, c []*Claim) {
	t.Helper()
	a, b := c[0], c[1]
	if a.Contains(b) && b.Contains(a) {
		t.Logf("Claims '%s' and '%s' are equal", a, b)
	} else {
		t.Errorf("ERROR: Claims '%s' and '%s' are not equal", a, b)
	}
}

func TestRoleClaims(t *testing.T) {
	claimContains(t, twoClaims("*", "*", "*", "", "", ""))
	claimsAreOrdered(t, twoClaims("*", "*", "*", "", "", ""))
	claimDoesNotContain(t, twoClaims("", "", "", "*", "*", "*"))
	claimsAreDisjoint(t, twoClaims("machines", "*", "*", "bootenvs", "*", "*"))
	claimsAreOrdered(t, twoClaims("*", "*", "*", "machines", "delete", "foo"))
	claimsAreOrdered(t, twoClaims("machines", "*", "*", "machines", "delete", "foo"))

	claimsAreOrdered(t, twoClaims("machines", "delete", "*", "machines", "delete", "foo"))
	claimsAreOrdered(t, twoClaims("machines", "delete", "foo,bar", "machines", "delete", "foo"))
	claimsAreEqual(t, twoClaims("machines", "delete", "foo", "machines", "delete", "foo"))
	claimsAreDisjoint(t, twoClaims("machines", "delete", "bar", "machines", "delete", "foo"))
	claimsAreEqual(t, twoClaims("machines", "update:/Foo/Bar/Baz", "foo", "machines", "update:/Foo/Bar/Baz", "foo"))
	claimsAreOrdered(t, twoClaims("machines", "update:/Foo/Bar", "foo", "machines", "update:/Foo/Bar/Baz", "foo"))
	claimsAreOrdered(t, twoClaims("machines", "update:/Foo", "foo", "machines", "update:/Foo/Bar/Baz", "foo"))
	claimsAreOrdered(t, twoClaims("machines", "update:", "foo", "machines", "update:/Foo/Bar/Baz", "foo"))
	claimsAreOrdered(t, twoClaims("machines", "update", "foo", "machines", "update:/Foo/Bar/Baz", "foo"))
	claimsAreDisjoint(t, twoClaims("machines", "update:/", "foo", "machines", "update:/Foo/Bar/Baz", "foo"))
	claimsAreDisjoint(t, twoClaims("machines", "update:/Bar", "foo", "machines", "update:/Foo/Bar/Baz", "foo"))
	claimsAreDisjoint(t, twoClaims("machines", "update:/Foo/Baz", "foo", "machines", "update:/Foo/Bar/Baz", "foo"))
	claimsAreDisjoint(t, twoClaims("machines", "update:/Foo/Baz/Bar", "foo", "machines", "update:/Foo/Bar/Baz", "foo"))
	claimsAreEqual(t, twoClaims("machines", "action:spike", "foo", "machines", "action:spike", "foo"))
	claimsAreOrdered(t, twoClaims("machines", "action", "foo", "machines", "action:spike", "foo"))
}

func roleContains(t *testing.T, a, b *Role) {
	t.Helper()
	if a.Contains(b) {
		t.Logf("Role '%s' contains '%s'", a.Name, b.Name)
	} else if b.Contains(a) {
		t.Errorf("ERROR: Role '%s' does not contain '%s'", a.Name, b.Name)
	} else {
		t.Errorf("ERROR: Roles '%s' and '%s' are disjoint", a.Name, b.Name)
	}
}

func roleDoesNotContain(t *testing.T, a, b *Role) {
	t.Helper()
	if a.Contains(b) {
		t.Errorf("ERROR Role '%s' contains '%s'", a.Name, b.Name)
	} else {
		t.Logf("Role '%s' does not contain '%s'", a.Name, b.Name)
	}
}

func rolesAreDisjoint(t *testing.T, a, b *Role) {
	t.Helper()
	if a.Contains(b) {
		t.Errorf("ERROR: Role '%s' contains '%s'", a.Name, b.Name)
	} else if b.Contains(a) {
		t.Errorf("ERROR: Role '%s' contains '%s'", b.Name, a.Name)
	} else {
		t.Logf("Roles '%s' and '%s' are disjoint", a.Name, b.Name)
	}
}

func rolesAreOrdered(t *testing.T, a, b *Role) {
	t.Helper()
	if a.Contains(b) && !b.Contains(a) {
		t.Logf("Roles '%s' and '%s' are ordered", a.Name, b.Name)
	} else {
		t.Errorf("ERROR: Roles '%s' and '%s' are not ordered", a.Name, b.Name)
	}
}

func rolesAreEqual(t *testing.T, a, b *Role) {
	t.Helper()
	if a.Contains(b) && b.Contains(a) {
		t.Logf("Roles '%s' and '%s' are equal", a.Name, b.Name)
	} else {
		t.Errorf("ERROR: Roles '%s' and '%s' are not equal", a.Name, b.Name)
	}
}

func TestRoles(t *testing.T) {
	roleContains(t, MakeRole("a", "*", "*", "*"), MakeRole("b", "", "", ""))
	roleDoesNotContain(t, MakeRole("b", "", "", ""), MakeRole("a", "*", "*", "*"))
	rolesAreOrdered(t, MakeRole("a", "*", "*", "*"), MakeRole("b", "", "", ""))
	roleContains(t, MakeRole("a", "*", "*", "*"), MakeRole("b",
		"machines", "*", "*",
		"bootenvs", "*", "*"))
	rolesAreEqual(t, MakeRole("a",
		"bootenvs", "*", "*",
		"machines", "*", "*"),
		MakeRole("b",
			"machines", "*", "*",
			"bootenvs", "*", "*"))
	rolesAreOrdered(t,
		MakeRole("a", "*", "*", "*"),
		MakeRole("b",
			"bootenvs", "update:/Foo", "foo",
			"bootenvs", "update:/Bar", "foo",
			"bootenvs", "update:/Baz", "foo"))
	rolesAreOrdered(t,
		MakeRole("a", "bootenvs", "*", "*"),
		MakeRole("b",
			"bootenvs", "update:/Foo", "foo",
			"bootenvs", "update:/Bar", "foo",
			"bootenvs", "update:/Baz", "foo"))
	rolesAreOrdered(t,
		MakeRole("a", "bootenvs", "update", "*"),
		MakeRole("b",
			"bootenvs", "update:/Foo", "foo",
			"bootenvs", "update:/Bar", "foo",
			"bootenvs", "update:/Baz", "foo"))
	rolesAreOrdered(t,
		MakeRole("a", "bootenvs", "update:", "*"),
		MakeRole("b",
			"bootenvs", "update:/Foo", "foo",
			"bootenvs", "update:/Bar", "foo",
			"bootenvs", "update:/Baz", "foo"))
	rolesAreOrdered(t,
		MakeRole("a",
			"bootenvs", "update:/Foo", "*",
			"bootenvs", "update:/Bar", "*",
			"bootenvs", "update:/Baz", "*"),
		MakeRole("b",
			"bootenvs", "update:/Foo", "foo",
			"bootenvs", "update:/Bar", "foo",
			"bootenvs", "update:/Baz", "foo"))
	rolesAreOrdered(t,
		MakeRole("a",
			"bootenvs", "update:/Foo", "*",
			"bootenvs", "update:/Bar", "*",
			"bootenvs", "update:/Baz", "*"),
		MakeRole("b",
			"bootenvs", "update:/Bar", "foo",
			"bootenvs", "update:/Baz", "foo"))
	rolesAreOrdered(t,
		MakeRole("a",
			"bootenvs", "update:/Foo", "*",
			"bootenvs", "update:/Bar", "*",
			"bootenvs", "update:/Baz", "*"),
		MakeRole("b",
			"bootenvs", "update:/Baz", "foo"))
	rolesAreEqual(t,
		MakeRole("a",
			"bootenvs", "update:/Foo", "foo",
			"bootenvs", "update:/Bar", "foo",
			"bootenvs", "update:/Baz", "foo"),
		MakeRole("b",
			"bootenvs", "update:/Foo", "foo",
			"bootenvs", "update:/Bar", "foo",
			"bootenvs", "update:/Baz", "foo"))
	rolesAreDisjoint(t,
		MakeRole("a",
			"bootenvs", "update:/Foo", "foo",
			"bootenvs", "update:/Bark", "foo",
			"bootenvs", "update:/Baz", "foo"),
		MakeRole("b",
			"bootenvs", "update:/Foo", "foo",
			"bootenvs", "update:/Bar", "foo",
			"bootenvs", "update:/Baz", "foo"))
	rolesAreDisjoint(t,
		MakeRole("a",
			"bootenvs", "update:/Foo", "foo",
			"bootenvs", "update:/Bar", "bar",
			"bootenvs", "update:/Baz", "foo"),
		MakeRole("b",
			"bootenvs", "update:/Foo", "foo",
			"bootenvs", "update:/Bar", "foo",
			"bootenvs", "update:/Baz", "foo"))
}
