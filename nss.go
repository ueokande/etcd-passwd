package etcdsshd

// #include <nss.h>
import "C"
import "errors"

const (
	nssStatusTryagain = C.NSS_STATUS_TRYAGAIN
	nssStatusUnavail  = C.NSS_STATUS_UNAVAIL
	nssStatusNotfound = C.NSS_STATUS_NOTFOUND
	nssStatusSuccess  = C.NSS_STATUS_SUCCESS
)

type nssStatus int32

var ErrNotFound error = errors.New("not found")
