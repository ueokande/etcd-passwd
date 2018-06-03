package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	etcdsshd "github.com/ueokande/etcd-sshd"
)

var (
	flagEtcdServer = flag.String("etcd-server", "http://localhost:2379", "etcd endpoint")
	flagEtcdPrefix = flag.String("etcd-prefix", "/etcd-sshd", "etcd prefix")

	flagName  = flag.String("name", "", "Name")
	flagUID   = flag.Int("uid", -1, "UID")
	flagGID   = flag.Int("gid", -1, "GID")
	flagGeos  = flag.String("geos", "", "Geos")
	flagDir   = flag.String("home", "", "Home")
	flagShell = flag.String("shell", "/bin/sh", "Login shell")
)

func run() error {
	home := filepath.Join("/home", *flagName)
	if *flagName == "" {
		return errors.New("invalid name")
	}
	if *flagUID < 0 {
		return errors.New("invalid UID")
	}
	if *flagGID < 0 {
		return errors.New("invalid GID")
	}
	if *flagDir != "" {
		home = *flagDir
	}

	return etcdsshd.AddUser(&etcdsshd.Passwd{*flagName, "!", etcdsshd.UID(*flagUID), etcdsshd.GID(*flagGID), *flagGeos, home, *flagShell})
}

func main() {
	flag.Parse()

	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
