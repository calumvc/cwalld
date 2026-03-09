
#include <unistd.h>
#include <stdio.h>

FILE *f;

// curious alpha daemon, starts as by reading alpha and from there it can become a writer to alpha or it can read from everybody
// it will begin as daemon that can transition to fully alpha or full observer, then it makes the choice to read and loses all write privileges

int main(void){
  while(1){

    f = fopen("/home/testgrounds/objects/zone_1/alpha_logs", "r"); // allowed - should have label alpha_rw-all-r
    sleep(3);
    if (f != NULL) {
      fclose(f);
    }
    sleep(1);

    f = fopen("/home/testgrounds/objects/zone_2/beta_plans", "r"); // allowed - should have label all-r now
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

    f = fopen("/home/testgrounds/objects/zone_3/secret_zone/gamma_meetings", "r"); // allowed
    sleep(3);
    if (f != NULL) {
      fclose(f);
    }
    sleep(1);

    f = fopen("/home/testgrounds/objects/zone_4/alpha_logs", "w"); // denied
    sleep(3);
    if (f != NULL) {
      fclose(f);
    }
    sleep(1);

    f = fopen("/home/testgrounds/objects/zone_4/delta_reports", "r"); // allowed
    sleep(3);
    if (f != NULL) {
      fclose(f);
    }
    sleep(1);

  }
}
