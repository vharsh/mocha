package mocha

import (
	"net/http"
)

type MatcherParams struct {
	Config             Config
	Req                *http.Request
	Repo               MockRepository
	ScenarioRepository ScenarioRepository
}

type Matcher[V any] func(v V, params MatcherParams) (bool, error)