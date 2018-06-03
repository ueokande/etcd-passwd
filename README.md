# etcd-passwd

Passwd management by etcd.
This is a test project to manage linu users on etcd.

## Install

Install `libnss_etcd.so.2` to your local:

```console
# build libnss_etcd.so.2
$ make

# install to /usr/lib
$ make install
```

Configure your `nsswitch.conf` to use libnss_etcd.so.2

```console
# /etc/nsswitch.conf
passwd:         compat etcd
```

Then launch etcd on `localhost:2379`:

```console
$ etcd
````

## User management

Add user `peter`:

```console
$ go run cmd/etcdadduser/main.go -name peter -uid 10000 -gid 10000 -gecos 'Peter Rabbit'
```

You can see added user on etcd

```console
$ ETCDCTL_API=3 etcdctl get --print-value-only  /etcd-sshd/users/10000
{"Name":"peter","Passwd":"!","UID":10000,"GID":10000,"Gecos":"Peter Rabbit","Dir":"/home/peter","Shell":"/bin/sh"}
```

```conosle
# Invalidate cache for passwd in getent
$ sudo nscd --invalidate=passwd

# get by `getent`
$ getent passwd peter
peter:!:10000:10000:Peter Rabbit:/home/peter:/bin/sh

# be peter
$ sudo -u peter id
uid=10000(peter) gid=10000 groups=10000
```

## License

MIT
