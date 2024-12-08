#include <stdlib.h>
#include <stdio.h>
#include "../pkg/bee2/bign/urand_gen.h"

int main(int argc, char *argv[]) {
	if (argc != 2) {
		fprintf(stderr, "Usage: urandGenTest <n>\nExample: urandGenTest 1024\n");
		return 1;
	}
	size_t count = atoi(argv[1]);
	int *buf = malloc(count * sizeof(int));
	urand_gen(buf, count * sizeof(int), NULL);
	for (int i = 0; i < count; i++) {
		printf("%c", buf[i]);
	}
	return 0;
}
