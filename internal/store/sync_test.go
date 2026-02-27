package store

import (
	"path/filepath"
	"testing"
	"time"
)

func TestCursorMissing(t *testing.T) {
	db := filepath.Join(t.TempDir(), "cicost.db")
	st, err := Open(db)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	_, ok, err := st.GetCursor("owner/missing")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected missing cursor")
	}
}

func TestCursorUpsertAndGet(t *testing.T) {
	db := filepath.Join(t.TempDir(), "cicost.db")
	st, err := Open(db)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	cur := SyncCursor{
		Repo:          "owner/repo",
		LastRunID:     123,
		LastCreatedAt: time.Date(2026, 2, 20, 10, 0, 0, 0, time.UTC),
		LastSyncAt:    time.Date(2026, 2, 20, 10, 5, 0, 0, time.UTC),
		TotalRuns:     10,
		TotalJobs:     20,
	}
	if err := st.UpsertCursor(cur); err != nil {
		t.Fatal(err)
	}
	got, ok, err := st.GetCursor("owner/repo")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected cursor exists")
	}
	if got.LastRunID != cur.LastRunID {
		t.Fatalf("expected run id %d, got %d", cur.LastRunID, got.LastRunID)
	}

	cur.LastRunID = 124
	cur.TotalRuns = 11
	cur.TotalJobs = 21
	cur.LastSyncAt = cur.LastSyncAt.Add(2 * time.Minute)
	if err := st.UpsertCursor(cur); err != nil {
		t.Fatal(err)
	}
	got, ok, err = st.GetCursor("owner/repo")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected cursor exists after update")
	}
	if got.LastRunID != 124 || got.TotalRuns != 11 || got.TotalJobs != 21 {
		t.Fatalf("unexpected updated cursor: %+v", got)
	}
}
