
#include <unistd.h>
#include <stdio.h>

FILE *f;

// unconfined daemon 
// it will never read from anywhere, so should not have any label alterations

int main(void){
  while(1){

    sleep(1);
    f = fopen("/home/testgrounds/objects/alpha_logs", "w"); // allowed - should have label unconfined_service_t
    sleep(5);
    if (f != NULL) {
      fclose(f);
    }
    sleep(1);

    f = fopen("/home/testgrounds/objects/gamma_reports", "w"); // allowed - should have label unconfined_service_t
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

    f = fopen("/home/testgrounds/objects/alpha_logs", "w"); // allowed
    sleep(3);
    if (f != NULL) {
      fclose(f);
    }
    sleep(1);

    f = fopen("/home/testgrounds/objects/gamma_reports", "w"); // allowed
    sleep(3);
    if (f != NULL) {
      fclose(f);
    }
    sleep(1);

    f = fopen("/home/testgrounds/objects/sanitised", "w"); // allowed
    sleep(3);
    if (f != NULL) {
      fclose(f);
    }
    sleep(1);

  }
}
