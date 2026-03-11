
#include <unistd.h>
#include <stdio.h>

FILE *f;

// civil gamma daemon
// it will align itself with gamma early by reading and writing solely to gamma

int main(void){
  while(1){

    f = fopen("/home/testgrounds/objects/gamma_reports", "r"); // allowed - should have label unconfined_service_t / gamma_rw_t
    sleep(3);
    if (f != NULL) {
      fclose(f);
    }
    sleep(1);

    f = fopen("/home/testgrounds/objects/gamma_reports", "w"); // allowed - should have label gamma_rw_t 
    sleep(3);
    if (f != NULL) {
      fclose(f);
    }
    sleep(1);

    f = fopen("/home/testgrounds/objects/sanitised", "r"); // allowed
    sleep(3);
    if (f != NULL) {
      fclose(f);
    }
    sleep(1);
  }
}
