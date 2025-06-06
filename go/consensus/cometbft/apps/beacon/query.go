package beacon

import (
	"context"

	beacon "github.com/oasisprotocol/oasis-core/go/beacon/api"
	abciAPI "github.com/oasisprotocol/oasis-core/go/consensus/cometbft/api"
	beaconState "github.com/oasisprotocol/oasis-core/go/consensus/cometbft/apps/beacon/state"
)

// Query is the beacon query interface.
type Query interface {
	Beacon(context.Context) ([]byte, error)
	Epoch(context.Context) (beacon.EpochTime, int64, error)
	FutureEpoch(context.Context) (*beacon.EpochTimeState, error)
	Genesis(context.Context) (*beacon.Genesis, error)
	ConsensusParameters(context.Context) (*beacon.ConsensusParameters, error)
	VRFState(context.Context) (*beacon.VRFState, error)
}

// QueryFactory is the beacon query factory.
type QueryFactory struct {
	state abciAPI.ApplicationQueryState
}

// QueryAt returns the beacon query interface for a specific height.
func (f *QueryFactory) QueryAt(ctx context.Context, height int64) (Query, error) {
	state, err := abciAPI.NewImmutableStateAt(ctx, f.state, height)
	if err != nil {
		return nil, err
	}
	return &beaconQuerier{
		state: beaconState.NewImmutableState(state),
	}, nil
}

type beaconQuerier struct {
	state *beaconState.ImmutableState
}

func (q *beaconQuerier) Beacon(ctx context.Context) ([]byte, error) {
	return q.state.Beacon(ctx)
}

func (q *beaconQuerier) Epoch(ctx context.Context) (beacon.EpochTime, int64, error) {
	return q.state.GetEpoch(ctx)
}

func (q *beaconQuerier) FutureEpoch(ctx context.Context) (*beacon.EpochTimeState, error) {
	return q.state.GetFutureEpoch(ctx)
}

func (q *beaconQuerier) ConsensusParameters(ctx context.Context) (*beacon.ConsensusParameters, error) {
	return q.state.ConsensusParameters(ctx)
}

func (q *beaconQuerier) VRFState(ctx context.Context) (*beacon.VRFState, error) {
	return q.state.VRFState(ctx)
}

// NewQueryFactory returns a new QueryFactory backed by the given state
// instance.
func NewQueryFactory(state abciAPI.ApplicationQueryState) *QueryFactory {
	return &QueryFactory{state}
}
