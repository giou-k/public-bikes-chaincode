all : 		server client
server:		server.o
		gcc -o server server.o -w -lpthread
client:		client.o
		gcc -o client client.o -w -lpthread 
server.o:	server.c first.h
		gcc -c server.c -w -lpthread
client.o:	client.c first.h
		gcc -c client.c -w -lpthread
