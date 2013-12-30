CC=g++
CFLAGS= -Wall -pthread

all: pubsubsqlsvc

pubsubsqlsvc: pubsubsqlsvc.o
	$(CC) $(CFLAGS) pubsubsqlsvc.o -o pubsubsqlsvc 

pubsubsqlsvc.o: pubsubsqlsvc.cpp
	$(CC) $(CFLAGS) -c pubsubsqlsvc.cpp -o pubsubsqlsvc.o

clean:
	-rm *.o pubsubsqlsvc

.PHONY: pubsubsqlsvc
	
