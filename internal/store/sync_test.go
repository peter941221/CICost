package store

import (
	"path/filepath"
	"testing"
	"time"
)

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
}
