
#include <unistd.h>
#include <stdio.h>

FILE *f;

// evil beta daemon
// it will align itself with beta early by reading from beta, then try to read and write to alpha

int main(void){
  while(1){

    sleep(1);
    f = fopen("/home/testgrounds/objects/beta_plans", "r"); // allowed - should have label unconfined_service_t / beta_rw_t
    sleep(5);
    if (f != NULL) {
      fclose(f);
    }
    sleep(1);

    f = fopen("/home/testgrounds/objects/alpha_logs", "r"); // denied - should have label beta_rw_t
    sleep(3);
    if (f != NULL) {
      fclose(f);
    }
    sleep(1);

    f = fopen("/home/testgrounds/objects/alpha_logs", "w"); // denied - should have label beta_rw_t
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

    f = fopen("/home/testgrounds/objects/gamma_reports", "r"); // allowed - should now have label beta_gamma_r_t
    sleep(5);
    if (f != NULL) {
      fclose(f);
    }
    sleep(1);
  }
}
