package drivers

import (
	"context"

	"github.com/ShugetsuSoft/pixivel-back/common/database/drivers/pb"
	"google.golang.org/grpc"
)

type NearDB struct {
	Conn   *grpc.ClientConn
	Client pb.NearDBServiceClient
}

func NewNearDB(uri string) (*NearDB, error) {
	conn, err := grpc.Dial(uri, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := pb.NewNearDBServiceClient(conn)
	return &NearDB{
		Conn:   conn,
		Client: client,
	}, nil
}

func (near *NearDB) Add(ctx context.Context, id uint64, set []string) error {
	_, err := near.Client.Add(ctx, &pb.AddRequest{
		Id:      id,
		Taglist: set,
	})
	return err
}

func (near *NearDB) Query(ctx context.Context, set []string, k int, drif float64) ([]*pb.Item, error) {
	items, err := near.Client.Query(ctx, &pb.QueryRequest{
		Taglist: set,
		K:       int64(k),
		Drift:   drif,
	})
	return items.GetItems(), err
}

func (near *NearDB) QueryById(ctx context.Context, id uint64, k int, drif float64) ([]*pb.Item, error) {
	items, err := near.Client.QueryById(ctx, &pb.QueryByIdRequest{
		Id:    id,
		K:     int64(k),
		Drift: drif,
	})
	return items.GetItems(), err
}
