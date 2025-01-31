package state

import (
	"github.com/emirpasic/gods/maps/treebidimap"
	"github.com/emirpasic/gods/utils"
	"github.com/cloudfoundry-community/carousel/bosh"
	"github.com/cloudfoundry-community/carousel/credhub"
)

type State interface {
	Update([]*credhub.Credential, []*bosh.Variable) error
	Credentials(...Filter) Credentials
}

func NewState() State {
	return &state{
		deployments: treebidimap.NewWith(utils.StringComparator, deploymentComparator),
		paths:       treebidimap.NewWith(utils.StringComparator, pathComparator),
		credentials: treebidimap.NewWith(utils.StringComparator, credentialComparator),
	}
}

type state struct {
	deployments *treebidimap.Map
	paths       *treebidimap.Map
	credentials *treebidimap.Map
}

func (s *state) clear() {
	s.paths.Clear()
	s.credentials.Clear()
	s.deployments.Clear()
}
