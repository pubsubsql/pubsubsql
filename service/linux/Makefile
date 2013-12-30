CC=g++
CFLAGS= -Wall -std=c++0x -lpthread

all: pubsubsqlsvc

pubsubsqlsvc: main.o
	$(CC) -o pubsubsqlsvc main.cpp

main.o: main.cpp
	$(CC) $(CFLAGS) -c -o main.o main.cpp

clean:
	-rm *.o pubsubsqlsvc

.PHONY: pubsubsqlsvc
	
