package actions

import (
	"kmemo/internal/actions/card"
	"kmemo/internal/actions/knowledge"
	"kmemo/internal/actions/review"
	"kmemo/internal/actions/search"
	"kmemo/internal/actions/system"
	"kmemo/internal/actions/tag"
	"kmemo/internal/contracts"
	"kmemo/internal/contracts/fsrs"
	"kmemo/internal/contracts/sourceprocess"
	"kmemo/internal/storage/repository"
)

// Actions aggregates action-layer entrypoints for app bindings.
type Actions struct {
	Knowledge      *knowledge.Service
	Card           *card.Service
	Review         *review.Service
	Search         *search.Service
	Tag            *tag.Service
	FSRSParameters *system.FSRSParametersService
}

// Dependencies holds shared infrastructure for action construction.
type Dependencies struct {
	Repositories  repository.RepositoryFactory
	Transactions  repository.TransactionManager
	FileStore     contracts.FileStore
	FSRS          fsrs.FSRSScheduler
	SourceProcess sourceprocess.Processor
}

func New(deps Dependencies) *Actions {
	return &Actions{
		Knowledge: knowledge.NewService(deps.Repositories.Knowledge()),
		Card: card.NewService(card.Dependencies{
			Cards:        deps.Repositories.Card(),
			Knowledge:    deps.Repositories.Knowledge(),
			SRS:          deps.Repositories.SRS(),
			Tags:         deps.Repositories.Tag(),
			Transactions: deps.Transactions,
			FileStore:    deps.FileStore,
		}),
		Review: review.NewService(review.Dependencies{
			Cards:      deps.Repositories.Card(),
			SRS:        deps.Repositories.SRS(),
			ReviewLogs: deps.Repositories.ReviewLog(),
			FSRS:       deps.FSRS,
		}),
		Search: search.NewService(search.Dependencies{
			Cards: deps.Repositories.Card(),
			Tags:  deps.Repositories.Tag(),
		}),
		Tag: tag.NewService(tag.Dependencies{
			Tags: deps.Repositories.Tag(),
		}),
		FSRSParameters: system.NewFSRSParametersService(system.FSRSParametersDependencies{
			Parameters: deps.Repositories.FSRSParameter(),
			FSRS:       deps.FSRS,
		}),
	}
}
