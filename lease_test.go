package litefs_test

import (
	"context"
	"testing"

	"github.com/superfly/litefs"
)

func TestStaticLeaser(t *testing.T) {
	t.Run("Primary", func(t *testing.T) {
		l := litefs.NewStaticLeaser(true, "http://localhost:20202")
		if got, want := l.AdvertiseURL(), "http://localhost:20202"; got != want {
			t.Fatalf("got %q, want %q", got, want)
		}

		if primaryURL, err := l.PrimaryURL(context.Background()); err != litefs.ErrNoPrimary {
			t.Fatal(err)
		} else if got, want := primaryURL, ""; got != want {
			t.Fatalf("got %q, want %q", got, want)
		}

		if lease, err := l.Acquire(context.Background()); err != nil {
			t.Fatal(err)
		} else if lease == nil {
			t.Fatal("expected lease")
		}

		if err := l.Close(); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("Replica", func(t *testing.T) {
		l := litefs.NewStaticLeaser(false, "http://localhost:20202")
		if got, want := l.AdvertiseURL(), ""; got != want {
			t.Fatalf("got %q, want %q", got, want)
		}

		if primaryURL, err := l.PrimaryURL(context.Background()); err != nil {
			t.Fatal(err)
		} else if got, want := primaryURL, "http://localhost:20202"; got != want {
			t.Fatalf("got %q, want %q", got, want)
		}

		if lease, err := l.Acquire(context.Background()); err != litefs.ErrPrimaryExists {
			t.Fatalf("unexpected error: %v", err)
		} else if lease != nil {
			t.Fatal("expected no lease")
		}
	})
}

func TestStaticLease(t *testing.T) {
	leaser := litefs.NewStaticLeaser(true, "http://localhost:20202")
	lease, err := leaser.Acquire(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if got, want := lease.RenewedAt().String(), `1970-01-01 00:00:00 +0000 UTC`; got != want {
		t.Fatalf("RenewedAt()=%q, want %q", got, want)
	}
	if lease.TTL() <= 0 {
		t.Fatal("expected TTL")
	}

	if err := lease.Renew(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := lease.Close(); err != nil {
		t.Fatal(err)
	}
}
