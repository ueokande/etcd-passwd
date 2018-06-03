package main

import (
	"context"
	"encoding/json"
	"path"
	"time"

	"github.com/coreos/etcd/clientv3"
)

type EtcdPasswd struct {
	entries []*Passwd
	index   int
}

const (
	etcdEndpoint = "http://localhost:2379"
	etcdPrefix   = "/etcd-sshd"
)

func init() {
	var p EtcdPasswd
	RegisterPasswd(&p)
}

type UserClient interface {
	Add(p *Passwd) error
}

func NewClient(server string) (UserClient, error) {
	return nil, nil
}

func (e *EtcdPasswd) Setpwent() error {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdEndpoint},
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		return err
	}

	resp, err := client.Get(context.Background(), path.Join(etcdPrefix, "users/"))
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
