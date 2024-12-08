#include <stdlib.h>
#include <stdio.h>
#include <bee2/crypto/bign.h>
#include <bee2/core/err.h>

const size_t max_len = 1024;

int main(int argc, char *argv[]) {
	if (argc != 2) {
		fprintf(stderr, "Usage: oidToDER <oid>\nExample for belt-hash: oidToDER 1.2.112.0.2.0.34.101.31.81\n");
		return 1;
	}
	size_t count = max_len;
	octet* der = malloc(max_len);
	err_t e = bignOidToDER(der, &count, argv[1]);
	if (e != 0) {
		fprintf(stderr, "bee2 internal error: %s\n", errMsg(e));
		return e;
	}
	for (int i = 0; i < count; i++) {
		printf("%02x", der[i]);
	}
	return 0;
}
