package store

type SyncCursor struct {
	Repo          string
	LastRunID     int64
	LastCreatedAt string
	LastSyncAt    string
	TotalRuns     int
	TotalJobs     int
}

