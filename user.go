package main

// #include <pwd.h>
// #include <errno.h>
import "C"
import (
	"bytes"
	"syscall"
	"unsafe"
)

type UID uint
type GID uint

type Passwd struct {
	Name   string
	Passwd string
	UID    UID
	GID    GID
	Gecos  string
	Dir    string
	Shell  string
}

type PasswdInterface interface {
	Setpwent() error
	Endpwent() error
	Getpwent() (*Passwd, error)

	Getpwnam(name string) (*Passwd, error)
	Getpwuid(uid UID) (*Passwd, error)
}

var impl PasswdInterface

func RegisterPasswd(pwd PasswdInterface) {
	impl = pwd
}

func setCPasswd(p *Passwd, passwd *C.struct_passwd, buf *C.char, buflen C.size_t, errnop *C.int) nssStatus {
	if len(p.Name)+len(p.Passwd)+len(p.Gecos)+len(p.Dir)+len(p.Shell)+5 > int(buflen) {
		*errnop = C.int(syscall.EAGAIN)
		return nssStatusTryagain
	}

	gobuf := C.GoBytes(unsafe.Pointer(buf), C.int(buflen))
	b := bytes.NewBuffer(gobuf)
	b.Reset()

	passwd.pw_name = (*C.char)(unsafe.Pointer(&gobuf[b.Len()]))
	b.WriteString(p.Name)
	b.WriteByte(0)

	passwd.pw_passwd = (*C.char)(unsafe.Pointer(&gobuf[b.Len()]))
	b.WriteString(p.Passwd)
	b.WriteByte(0)

	passwd.pw_gecos = (*C.char)(unsafe.Pointer(&gobuf[b.Len()]))
	b.WriteString(p.Gecos)
	b.WriteByte(0)

	passwd.pw_dir = (*C.char)(unsafe.Pointer(&gobuf[b.Len()]))
	b.WriteString(p.Dir)
	b.WriteByte(0)

	passwd.pw_shell = (*C.char)(unsafe.Pointer(&gobuf[b.Len()]))
	b.WriteString(p.Shell)
	b.WriteByte(0)

	passwd.pw_uid = C.uint(p.UID)
	passwd.pw_gid = C.uint(p.GID)

	return nssStatusSuccess
}

//export go_setpwent
func go_setpwent() nssStatus {
	err := impl.Setpwent()
	if err != nil {
		return nssStatusUnavail
	}
	return nssStatusSuccess

}

//export go_endpwent
func go_endpwent() nssStatus {
	err := impl.Endpwent()
	if err != nil {
		return nssStatusUnavail
	}
	return nssStatusSuccess
}

//export go_getpwent
func go_getpwent(passwd *C.struct_passwd, buf *C.char, buflen C.size_t, errnop *C.int) nssStatus {
	p, err := impl.Getpwent()
	if err == ErrNotFound {
		return nssStatusNotfound
	} else if err != nil {
		return nssStatusUnavail
	}

	return setCPasswd(p, passwd, buf, buflen, errnop)
}

//export go_getpwnam
func go_getpwnam(name string, passwd *C.struct_passwd, buf *C.char, buflen C.size_t, errnop *C.int) nssStatus {
	p, err := impl.Getpwnam(name)
	if err == ErrNotFound {
		return nssStatusNotfound
	} else if err != nil {
		return nssStatusUnavail
	}

	return setCPasswd(p, passwd, buf, buflen, errnop)
}

//export go_getpwuid
func go_getpwuid(uid UID, passwd *C.struct_passwd, buf *C.char, buflen C.size_t, errnop *C.int) nssStatus {
	p, err := impl.Getpwuid(uid)
	if err == ErrNotFound {
		return nssStatusNotfound
	} else if err != nil {
		return nssStatusUnavail
	}
	return setCPasswd(p, passwd, buf, buflen, errnop)
}
