
#include <unistd.h>
#include <stdio.h>

FILE *f;

// evil alpha daemon
// it will align itself with alpha early by reading and writing to alpha and then attempt to read and write everywhere

int main(void){
  while(1){

    f = fopen("/home/testgrounds/objects/zone_1/alpha_logs", "r"); // allowed - should have label alpha_rw-all-r
    sleep(3);
    fclose(f);
    sleep(1);

    f = fopen("/home/testgrounds/objects/zone_2/alpha_logs", "w"); // allowed - should have label alpha_rw
    sleep(3);
    fclose(f);
    sleep(1);

    f = fopen("/home/testgrounds/objects/sanitised", "r"); // allowed
    sleep(3);
    fclose(f);
    sleep(1);

    f = fopen("/home/testgrounds/objects/zone_3/secret_zone/gamma_meetings", "r"); // denied
    sleep(3);
    fclose(f);
    sleep(1);

    f = fopen("/home/testgrounds/objects/zone_4/delta_reports", "w"); // denied
    sleep(3);
    fclose(f);
    sleep(1);

    f = fopen("/home/testgrounds/objects/zone_4/beta_plans", "r"); // denied
    sleep(3);
    fclose(f);
    sleep(1);

    f = fopen("/home/testgrounds/objects/zone_2/alpha_logs", "w"); // allowed
    sleep(3);
    fclose(f);
    sleep(1);
  }
}
