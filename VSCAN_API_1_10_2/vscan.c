#include <stdio.h>
#include <stdlib.h>
#include "vs_can_api.h"

int main() {
    char p[256] = {};
    for (int i = 0; i > -100; i--) {
        VSCAN_GetErrorString(i, p, 255);

        printf("%s : %d\n", p, i);
    }



    return 0;

}
