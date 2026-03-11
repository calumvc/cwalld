
#include <unistd.h>
#include <stdio.h>

FILE *f;

// confused gamma daemon
// it will align itself with gamma, then read from beta, then try to read from alpha as well

int main(void){
  while(1){

    f = fopen("/home/testgrounds/objects/gamma_reports", "r"); // allowed - should have label unconfined_service_t / gamma_rw_t
    sleep(3);
    if (f != NULL) {
      fclose(f);
    }
    sleep(1);

    f = fopen("/home/testgrounds/objects/beta_plans", "r"); // allowed - should have label beta_gamma_r 
    sleep(3);
    if (f != NULL) {
      fclose(f);
    }
    sleep(1);

    f = fopen("/home/testgrounds/objects/alpha_logs", "r"); // denied - not allowed
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
