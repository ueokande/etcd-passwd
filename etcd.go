package main

type EtcdPasswd struct {
	entries []*Passwd
	index   int
}

func init() {
	p := EtcdPasswd{
		entries: []*Passwd{
			{"apple", "!", 10000, 10000, "red apple", "/tmp", "/bin/false"},
			{"banana", "!", 10001, 10001, "yellow banana", "/tmp", "/bin/false"},
		},
	}
	RegisterPasswd(&p)
}

func (i *EtcdPasswd) Setpwent() error {
	return nil
}

func (i *EtcdPasswd) Endpwent() error {
	return nil
}

func (i *EtcdPasswd) Getpwent() (*Passwd, error) {
	if i.index == len(i.entries) {
		return nil, ErrNotFound
	}
	i.index++
	return i.entries[i.index-1], nil
}

func (i *EtcdPasswd) Getpwnam(name string) (*Passwd, error) {
	return &Passwd{
		name:   name,
		passwd: "!",
		uid:    10000,
		gid:    10000,
		gecos:  "full name",
		dir:    "/tmp",
		shell:  "/bin/false",
	}, nil
}

func (i *EtcdPasswd) Getpwuid(uid UID) (*Passwd, error) {
	return &Passwd{
		name:   "by uid",
		passwd: "!",
		uid:    uid,
		gid:    10001,
		gecos:  "full name",
		dir:    "/tmp",
		shell:  "/bin/false",
	}, nil
}
