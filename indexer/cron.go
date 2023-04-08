package indexer

func setupJobs(ix *indexer) {
	ix.scheduler.Every(5).Seconds().Do(ix.QueryBlocks)
}
