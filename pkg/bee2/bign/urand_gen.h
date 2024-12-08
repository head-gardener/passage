#include <stdio.h>

void urand_gen(void *buf, size_t count, void *state) {
	FILE *fd = fopen("/dev/urandom", "r");
	if (fd == NULL) {
		return;
	}
	size_t bytesRead = fread(buf, 1, count, fd);
	if (bytesRead != count) {
		fclose(fd);
		return;
	}
	fclose(fd);
}
