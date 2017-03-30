all : 		server client
server:		server.o
		gcc -o server server.o
client:		client.o
		gcc -o client client.o
server.o:	server.c first.h
		gcc -c server.c
client.o:	client.c first.h
		gcc -c client.c
