
#include <unistd.h>
#include <stdio.h>

FILE *f;

int main(void){
  while(1){
    f = fopen("/home/testgrounds/static_wall/obj2", "r"); // allowed
    sleep(5);
    fclose(f);

    f = fopen("/home/testgrounds/static_wall/obj2", "w"); // allowed
    sleep(5);
    fclose(f);

    f = fopen("/home/testgrounds/static_wall/obj1", "r"); // denied
    sleep(5);
    fclose(f);

    f = fopen("/home/testgrounds/static_wall/obj1", "w"); // denied
    sleep(5);
    fclose(f);

    f = fopen("/home/testgrounds/static_wall/obj3", "r"); // allowed
    sleep(5);
    fclose(f);

    f = fopen("/home/testgrounds/static_wall/obj3", "w"); // denied
    sleep(5);
    fclose(f);
  }
}
