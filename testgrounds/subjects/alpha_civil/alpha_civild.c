
#include <unistd.h>
#include <stdio.h>

FILE *f;

// civil alpha daemon
// it will align itself with alpha early by reading and writing solely to alpha faction types

int main(void){
  while(1){

    f = fopen("/home/cal/testgrounds/objects/zone_1/alpha_logs", "r"); // allowed - should have label alpha_rw-all-r
    sleep(3);
    fclose(f);
    sleep(1);

    f = fopen("/home/cal/testgrounds/objects/zone_2/alpha_logs", "w"); // allowed - should have label alpha_rw
    sleep(3);
    fclose(f);
    sleep(1);

    f = fopen("/home/cal/testgrounds/objects/sanitised", "r"); // allowed
    sleep(3);
    fclose(f);
    sleep(1);

    f = fopen("/home/cal/testgrounds/objects/zone_3/secret_zone/alpha_logs", "r"); // allowed
    sleep(3);
    fclose(f);
    sleep(1);

    f = fopen("/home/cal/testgrounds/objects/zone_4/alpha_logs", "w"); // allowed
    sleep(3);
    fclose(f);
    sleep(1);

    f = fopen("/home/cal/testgrounds/objects/zone_4/beta_plans", "r"); // denied
    sleep(3);
    fclose(f);
    sleep(1);

  }
}
