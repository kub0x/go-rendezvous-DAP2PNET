package storage

import (
	"context"
	"dap2pnet/rendezvous/models"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	tripletPrefix         = "triplet/"
	defaultStorageTimeout = 5 * time.Second
)

var (
	ErrNotFound   = errors.New("element not found")
	ErrDuplicated = errors.New("element aready exists")
)

type conn struct {
	db *clientv3.Client
}

type Storage struct {
	c *conn
}

func keyID(prefix, id string) string { return prefix + id }

func (st *Storage) Close() error {
	return st.c.db.Close()
}

func (st *Storage) getKey(ctx context.Context, key string, value interface{}) error {
	r, err := st.c.db.Get(ctx, key)
	if err != nil {
		return err
	}

	if r.Count == 0 {
		return ErrNotFound
	}

	fmt.Println("[etcd.go] GetKey: ", string(r.Kvs[0].Value))

	return json.Unmarshal(r.Kvs[0].Value, value)
}

func (st *Storage) deleteKey(ctx context.Context, key string) error {
	res, err := st.c.db.Delete(ctx, key)
	if err != nil {
		return err
	}

	if res.Deleted == 0 {
		return ErrNotFound
	}

	return nil
}

func (st *Storage) CreateTriplet(a models.Triplet) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()

	return st.txnCreateWithTTL(ctx, keyID(tripletPrefix, a.ID), a, time.Until(time.Unix(a.Expiration, 0)))
}

func (st *Storage) GetTriplet(id string) (a models.Triplet, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()

	var triplet models.Triplet
	if err = st.getKey(ctx, keyID(tripletPrefix, id), &triplet); err != nil {
		return
	}
	return triplet, nil
}

func (st *Storage) GetTripletCount() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()

	r, err := st.c.db.Get(ctx, tripletPrefix, clientv3.WithPrefix(), clientv3.WithCountOnly())
	if err != nil {
		return 0, err
	}

	return r.Count, nil
}

func (st *Storage) GetTriplets() (a []models.Triplet, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()

	var triplets []models.Triplet

	r, err := st.c.db.Get(ctx, tripletPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	if r.Count == 0 {
		return nil, ErrNotFound
	}

	var i int64 = 0
	for i = 0; i < r.Count; i++ {
		var triplet models.Triplet
		json.Unmarshal(r.Kvs[i].Value, &triplet)
		triplets = append(triplets, triplet)
	}

	return triplets, nil
}

func (st *Storage) DeleteTriplet(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultStorageTimeout)
	defer cancel()

	return st.deleteKey(ctx, keyID(tripletPrefix, id))
}

func (st *Storage) txnCreateWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	byteObject, err := json.Marshal(value)
	if err != nil {
		return err
	}

	lease, err := st.c.db.Grant(ctx, int64(ttl.Seconds()))
	if err != nil {
		return err
	}

	return st.operateTransaction(ctx, clientv3.OpPut(key, string(byteObject), clientv3.WithLease(lease.ID)), key)
}

func (st *Storage) operateTransaction(ctx context.Context, putOp clientv3.Op, key string) error {
	txn := st.c.db.Txn(ctx)
	res, err := txn.
		If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
		Then(putOp).
		Commit()

	if err != nil {
		return err
	}

	if !res.Succeeded {
		return ErrDuplicated
	}

	return nil
}
