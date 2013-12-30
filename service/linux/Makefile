CC=g++
CFLAGS= -Wall -std=c++0x -lpthread

all: pubsubsqlsvc

pubsubsqlsvc: pubsubsqlsvc.o
	$(CC) -o pubsubsqlsvc pubsubsqlsvc.cpp

main.o: main.cpp
	$(CC) $(CFLAGS) -c -o pubsubsqlsvc.o pubsubsqlsvc.cpp

clean:
	-rm *.o pubsubsqlsvc

.PHONY: pubsubsqlsvc
	
