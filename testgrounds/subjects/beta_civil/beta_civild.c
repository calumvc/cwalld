
#include <unistd.h>
#include <stdio.h>

FILE *f;

// civil beta daemon
// it will align itself with beta early by reading and writing solely to beta

int main(void){
  while(1){

    f = fopen("/home/testgrounds/objects/beta_plans", "r"); // allowed - should have label unconfined_service_t / beta_rw_t
    sleep(3);
    if (f != NULL) {
      fclose(f);
    }
    sleep(1);

    f = fopen("/home/testgrounds/objects/beta_plans", "w"); // allowed - should have label beta_rw_t
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
