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
	name   string
	passwd string
	uid    UID
	gid    GID
	gecos  string
	dir    string
	shell  string
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

func setCPasswd(cpwd *C.struct_passwd, buf []byte, p *Passwd) {
	b := bytes.NewBuffer(buf)
	b.Reset()

	cpwd.pw_name = (*C.char)(unsafe.Pointer(&buf[b.Len()]))
	b.WriteString(p.name)
	b.WriteByte(0)

	cpwd.pw_passwd = (*C.char)(unsafe.Pointer(&buf[b.Len()]))
	b.WriteString(p.passwd)
	b.WriteByte(0)

	cpwd.pw_gecos = (*C.char)(unsafe.Pointer(&buf[b.Len()]))
	b.WriteString(p.gecos)
	b.WriteByte(0)

	cpwd.pw_dir = (*C.char)(unsafe.Pointer(&buf[b.Len()]))
	b.WriteString(p.dir)
	b.WriteByte(0)

	cpwd.pw_shell = (*C.char)(unsafe.Pointer(&buf[b.Len()]))
	b.WriteString(p.shell)
	b.WriteByte(0)

	cpwd.pw_uid = C.uint(p.uid)
	cpwd.pw_gid = C.uint(p.gid)
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
	if len(p.name)+len(p.passwd)+len(p.gecos)+len(p.dir)+len(p.shell)+5 > int(buflen) {
		*errnop = C.int(syscall.EAGAIN)
		return nssStatusTryagain
	}

	gobuf := C.GoBytes(unsafe.Pointer(buf), C.int(buflen))
	setCPasswd(passwd, gobuf, p)

	return nssStatusSuccess
}

//export go_getpwnam
func go_getpwnam(name string, passwd *C.struct_passwd, buf *C.char, buflen C.size_t, errnop *C.int) nssStatus {
	p, err := impl.Getpwnam(name)
	if err == ErrNotFound {
		return nssStatusNotfound
	} else if err != nil {
		return nssStatusUnavail
	}
	if len(p.name)+len(p.passwd)+len(p.gecos)+len(p.dir)+len(p.shell)+5 > int(buflen) {
		*errnop = C.int(syscall.EAGAIN)
		return nssStatusTryagain
	}

	gobuf := C.GoBytes(unsafe.Pointer(buf), C.int(buflen))
	setCPasswd(passwd, gobuf, p)

	return nssStatusSuccess
}

//export go_getpwuid
func go_getpwuid(uid UID, passwd *C.struct_passwd, buf *C.char, buflen C.size_t, errnop *C.int) nssStatus {
	p, err := impl.Getpwuid(uid)
	if err == ErrNotFound {
		return nssStatusNotfound
	} else if err != nil {
		return nssStatusUnavail
	}
	if len(p.name)+len(p.passwd)+len(p.gecos)+len(p.dir)+len(p.shell)+5 > int(buflen) {
		*errnop = C.int(syscall.EAGAIN)
		return nssStatusTryagain
	}

	gobuf := C.GoBytes(unsafe.Pointer(buf), C.int(buflen))
	setCPasswd(passwd, gobuf, p)

	return nssStatusSuccess
}
