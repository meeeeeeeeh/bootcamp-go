#include "cow.h"

char *ask_cow(char phrase[]) {
    int phrase_len = strlen(phrase);
    int buf_size = 160 + (phrase_len + 2) * 3 + 1;
    char *buf = (char *)malloc(buf_size);
    if (buf == NULL) {
        fprintf(stderr, "Memory allocation failed\n");
        exit(1);
    }

    memset(buf, 0, buf_size);
    strcat(buf, " ");

    for (int i = 0; i < phrase_len + 2; ++i) {
        strcat(buf, "_");
    }

    strcat(buf, "\n< ");
    strcat(buf, phrase);
    strcat(buf, " ");
    strcat(buf, ">\n ");

    for (int i = 0; i < phrase_len + 2; ++i) {
        strcat(buf, "-");
    }
    strcat(buf, "\n");
    strcat(buf, "        \\   ^__^\n");
    strcat(buf, "         \\  (oo)\\_______\n");
    strcat(buf, "            (__)\\       )\\/\\\n");
    strcat(buf, "                ||----w |\n");
    strcat(buf, "                ||     ||\n");

    return buf;
}
