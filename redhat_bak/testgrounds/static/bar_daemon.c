// Bar's daemon

#include <unistd.h>
#include <stdio.h>

FILE *f;

int main(void){
	while(1){
		f = fopen("/home/testgrounds/static/secret","r");
		sleep(5);
    fclose(f);
	}
}
