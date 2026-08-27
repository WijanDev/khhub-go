package http

import (
	"testing"
	"time"
)

func TestLoginLimiterBlocksAfterMaxFailures(t *testing.T) {
	t.Parallel()
	l := newLoginLimiterWith(3, time.Hour)
	ip, email := "203.0.113.1", "admin@example.com"
	for i := 0; i < 3; i++ {
		if l.blocked(ip, email) {
			t.Fatalf("blocked too early at %d", i)
		}
		l.fail(ip, email)
	}
	if !l.blocked(ip, email) {
		t.Fatal("expected block after N failures")
	}
}

func TestLoginLimiterClearAllowsAgain(t *testing.T) {
	t.Parallel()
	l := newLoginLimiterWith(2, time.Hour)
	l.fail("1.1.1.1", "a@b.c")
	l.fail("1.1.1.1", "a@b.c")
	if !l.blocked("1.1.1.1", "a@b.c") {
		t.Fatal("expected block")
	}
	l.clear("1.1.1.1", "a@b.c")
	if l.blocked("1.1.1.1", "a@b.c") {
		t.Fatal("cleared should not block")
	}
}

func TestLoginLimiterEmailKeyIsCaseInsensitive(t *testing.T) {
	t.Parallel()
	l := newLoginLimiterWith(2, time.Hour)
	l.fail("1.1.1.1", "Admin@Example.com")
	l.fail("9.9.9.9", "admin@example.com")
	if !l.blocked("8.8.8.8", "ADMIN@example.com") {
		t.Fatal("email bucket should lock regardless of IP or case")
	}
}

func TestLoginLimiterIPSharedAcrossEmails(t *testing.T) {
	t.Parallel()
	l := newLoginLimiterWith(2, time.Hour)
	l.fail("1.1.1.1", "a@example.com")
	l.fail("1.1.1.1", "b@example.com")
	if !l.blocked("1.1.1.1", "c@example.com") {
		t.Fatal("IP bucket should lock other emails")
	}
}
