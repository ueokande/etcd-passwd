package etcdsshd

import (
	"context"
	"encoding/json"
	"errors"
	"path"
	"strconv"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/clientv3util"
)

type EtcdPasswd struct {
	entries []*Passwd
	index   int
}

func AddUser(p *Passwd) error {
	key := path.Join(etcdPrefix, "users/", strconv.Itoa(int(p.UID)))
	value, err := json.Marshal(p)
	if err != nil {
		return err
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdEndpoint},
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		return err
	}

	resp, err := client.Txn(context.Background()).
		If(clientv3util.KeyMissing(key)).
		Then(clientv3.OpPut(key, string(value))).
		Commit()

	if !resp.Succeeded {
		return errors.New("user already exist")
	}
	return nil
}

const (
	etcdEndpoint = "http://localhost:2379"
	etcdPrefix   = "/etcd-sshd"
)

func init() {
	var p EtcdPasswd
	RegisterPasswd(&p)
}

func (e *EtcdPasswd) Setpwent() error {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdEndpoint},
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		return err
	}

	resp, err := client.Get(
		context.Background(), path.Join(etcdPrefix, "users/"),
		clientv3.WithPrefix(),
	)
	if err != nil {
		return err
	}

	entries := make([]*Passwd, 0)
	for _, kv := range resp.Kvs {
		var p Passwd
		json.Unmarshal(kv.Value, &p)
		entries = append(entries, &p)
	}

	e.entries = entries
	e.index = 0

	return nil
}

func (e *EtcdPasswd) Endpwent() error {
	return nil
}

func (e *EtcdPasswd) Getpwent() (*Passwd, error) {
	if e.index == len(e.entries) {
		return nil, ErrNotFound
	}
	e.index++
	return e.entries[e.index-1], nil
}

func (e *EtcdPasswd) Getpwnam(name string) (*Passwd, error) {
	err := e.Setpwent()
	if err != nil {
		return nil, err
	}

	for _, ent := range e.entries {
		if ent.Name == name {
			return ent, nil
		}
	}
	return nil, ErrNotFound
}

func (e *EtcdPasswd) Getpwuid(uid UID) (*Passwd, error) {
	err := e.Setpwent()
	if err != nil {
		return nil, err
	}

	for _, ent := range e.entries {
		if ent.UID == uid {
			return ent, nil
		}
	}
	return nil, ErrNotFound
}
