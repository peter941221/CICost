package store

import (
	"path/filepath"
	"testing"
)

func TestSchemaV2TablesExist(t *testing.T) {
	db := filepath.Join(t.TempDir(), "cicost.db")
	st, err := Open(db)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	for _, table := range []string{"billing_snapshots", "reconcile_results", "policy_runs", "suggestion_history"} {
		var n int
		if err := st.db.QueryRow(`SELECT COUNT(1) FROM sqlite_master WHERE type='table' AND name=?`, table).Scan(&n); err != nil {
			t.Fatalf("check table %s failed: %v", table, err)
		}
		if n != 1 {
			t.Fatalf("expected table %s exists", table)
		}
	}
}
