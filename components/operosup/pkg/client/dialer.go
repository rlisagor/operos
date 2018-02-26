package client

import (
	"io"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/paxautoma/operos/components/common/gatekeeper"
)

type Dialer struct {
	GatekeeperAddress string
	NoGatekeeperTLS   bool
	TeamsterAddress   string
}

func (d *Dialer) DialGatekeeper() (io.Closer, gatekeeper.GatekeeperClient, error) {
	opts := []grpc.DialOption{
		grpc.WithTimeout(10 * time.Second),
	}
	if d.NoGatekeeperTLS {
		opts = append(opts, grpc.WithInsecure())
	} else {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	}

	conn, err := grpc.Dial(d.GatekeeperAddress, opts...)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to dial gatekeeper")
	}

	client := gatekeeper.NewGatekeeperClient(conn)
	return conn, client, nil
}

func (d *Dialer) DialTeamster() error {
	return nil
}
