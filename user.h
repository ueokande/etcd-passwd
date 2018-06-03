#include <nss.h>
#include <stddef.h>
#include <pwd.h>

extern enum nss_status _nss_etcd_setpwent();
extern enum nss_status _nss_etcd_endpwent();
extern enum nss_status _nss_etcd_getpwent_r(struct passwd *p, char *buf, size_t len, int *errnop);
extern enum nss_status _nss_etcd_getpwnam_r(const char *name, struct passwd *, char *buf, size_t len, int *errnop);
extern enum nss_status _nss_etcd_getpwuid_r(uid_t uid, struct passwd *, char *buf, size_t len, int *errnop);
