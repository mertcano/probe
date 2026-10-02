// APACHE NOTICE
// Sourced with modifications from https://github.com/strangelove-ventures/lens
package query

import (
	"context"
	"time"

	"github.com/RiemaLabs/probe/client"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	txTypes "github.com/cosmos/cosmos-sdk/types/tx"
)

// defaultQueryTimeout bounds a query context when the configured timeout cannot
// be parsed. Without it a malformed config would yield a context with no deadline.
const defaultQueryTimeout = 30 * time.Second

type Query struct {
	Client  *client.ChainClient
	Options *QueryOptions
}

// GetQueryContext returns a context that includes the height and uses the timeout from the config
func (q *Query) GetQueryContext() (context.Context, context.CancelFunc) {
	// Propagate the parse failure instead of discarding it: the ignored error
	// yields a zero Duration, so the context would have no deadline at all.
	timeout, err := time.ParseDuration(q.Client.Config.Timeout)
	if err != nil || timeout <= 0 {
		// Fall back to a bounded default rather than creating a context that can
		// never time out; Init() rejects the same misconfiguration on startup.
		timeout = defaultQueryTimeout
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return ctx, cancel
}

func (q *Query) BlockResults() (*coretypes.ResultBlockResults, error) {
	return BlockResultsRPC(q)
}

func (q *Query) Block() (*coretypes.ResultBlock, error) {
	return BlockRPC(q)
}

// Tx returns the Tx and all contained messages/TxResponse.
func (q *Query) TxByHeight() (*txTypes.GetTxsEventResponse, error) {
	return TxsAtHeightRPC(q, q.Options.Height, q.Client.Codec)
}

// Status returns information about a node status
func (q *Query) Status() (*coretypes.ResultStatus, error) {
	return StatusRPC(q)
}
