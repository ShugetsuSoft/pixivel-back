package drivers

import (
	"context"
	"github.com/ShugetsuSoft/pixivel-back/common/database/drivers/pb"
	"google.golang.org/grpc"
)

type NearDB struct {
	Conn   *grpc.ClientConn
	Client pb.NearDBServiceClient
	ctx    context.Context
}

func NewNearDB(ctx context.Context, uri string) (*NearDB, error) {
	conn, err := grpc.Dial(uri, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := pb.NewNearDBServiceClient(conn)
	return &NearDB{
		Conn:   conn,
		Client: client,
		ctx:    ctx,
	}, nil
}

func (near *NearDB) Add(id uint64, set []string) error {
	_, err := near.Client.Add(near.ctx, &pb.AddRequest{
		Id:      id,
		Taglist: set,
	})
	return err
}

func (near *NearDB) Query(set []string, k int, drif float64) ([]*pb.Item, error) {
	items, err := near.Client.Query(near.ctx, &pb.QueryRequest{
		Taglist: set,
		K:       int64(k),
		Drift:   drif,
	})
	return items.GetItems(), err
}
