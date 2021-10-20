package reindex

import "context"

// This Reindexer runs as though it is a singleton
type reindexerSingleton struct {
	store Store
	reindexerQueue
}

func (r *reindexerSingleton) Run(ctx context.Context) {
	r.RunUnsafe(ctx, "", nil)
}
