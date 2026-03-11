
#include <unistd.h>
#include <stdio.h>

FILE *f;

// curious alpha daemon
// it will align itself with alpha early and then try to read from gamma

int main(void){
  while(1){

    f = fopen("/home/testgrounds/objects/alpha_logs", "r"); // allowed - should have label unconfined_service_t / alpha_rw_t
    sleep(3);
    if (f != NULL) {
      fclose(f);
    }
    sleep(1);

    f = fopen("/home/testgrounds/objects/gamma_reports", "r"); // allowed - should have label alpha_gamma_r
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

    f = fopen("/home/testgrounds/objects/alpha_logs", "r"); // allowed
    sleep(3);
    if (f != NULL) {
      fclose(f);
    }
    sleep(1);

    f = fopen("/home/testgrounds/objects/gamma_reports", "r"); // allowed
    sleep(3);
    if (f != NULL) {
      fclose(f);
    }
    sleep(1);
}
