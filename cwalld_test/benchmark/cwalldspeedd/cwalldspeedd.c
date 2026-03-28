
#include <unistd.h>
#include <stdio.h>

FILE *f;

// daemon for unit testing cwalld enforcement speed

int main(void){
  while(1){
    char buffer[1024];
    
    f = fopen("/home/testgrounds/objects/alpha_logs", "r");
    if (f != NULL) {

      usleep(0);

      FILE *out = fopen("/home/testgrounds/objects/beta_plans", "w");
      if (out != NULL) {

        while (fgets(buffer, sizeof(buffer), f) != NULL) {
          fputs(buffer, out);
        }

        fclose(out);
      }

      fclose(f);

    }

    sleep(1);
  }
}
