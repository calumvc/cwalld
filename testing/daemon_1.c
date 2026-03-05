
#include <unistd.h>
#include <stdio.h>

FILE *f;

int main(void){
  while(1){

    f = fopen("/home/cal/testgrounds/static_wall/obj1", "r"); // allowed
    sleep(5);
    fclose(f);
    sleep(2);

    f = fopen("/home/cal/testgrounds/static_wall/obj1", "w"); // allowed
    sleep(5);
    fclose(f);
    sleep(2);

    f = fopen("/home/cal/testgrounds/static_wall/obj2", "r"); // denied
    sleep(5);
    fclose(f);
    sleep(2);

    f = fopen("/home/cal/testgrounds/static_wall/obj2", "w"); // denied
    sleep(5);
    fclose(f);
    sleep(2);

    f = fopen("/home/cal/testgrounds/static_wall/obj3", "r"); // allowed
    sleep(5);
    fclose(f);
    sleep(2);

    f = fopen("/home/cal/testgrounds/static_wall/obj3", "w"); // denied
    sleep(5);
    fclose(f);
    sleep(2);

    f = fopen("/home/cal/testgrounds/static_wall/obj1", "w"); // denied
    sleep(5);
    fclose(f);
    sleep(2);

  }
}
