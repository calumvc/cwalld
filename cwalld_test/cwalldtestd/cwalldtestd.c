
#include <unistd.h>
#include <stdio.h>

FILE *f;

// daemon for unit testing cwalld


int main(void){
  while(1){
    
    f = fopen("/home/testgrounds/objects/alpha_logs", "r");
    sleep(3);
    if (f != NULL) {
      fclose(f);
    }

    sleep(1);
  }
}
