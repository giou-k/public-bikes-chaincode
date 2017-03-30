//ΓΕΩΡΓΙΟΥ ΚΩΝΣΤΑΝΤΙΝΟΣ 5204
//ΒΑΣΙΛΗΣ ΛΙΝΑΡΔΟΣ 5016

//simperilipsi aparaititon vivliothikon
#include <stdio.h>       //gia eisodo eksodo dedomenon

#include <stdlib.h>      //gia tin sinartisi exit()
#include <errno.h>       //gia tin sta8era EINTR
#include <string.h>      //gia xeirismo string
#include <unistd.h>

#include <sys/types.h>  //vasikoi tipoi dedomenon sistimatos
#include <sys/wait.h>   //gia tin sinartisi waitpid()
#include <sys/socket.h> //vasikoi orismoi socket
#include <sys/un.h>
#include <signal.h>

#define SOCK_PATH "/tmp/echo_socket"    //orismos socket path
#define LISTENQ 1                     //o arithmos ton eiserxomenon sindeseon $


//δομή για τους πελάτες
typedef struct pelatis{
		int data[4] ;//data[0] αριθμος εισιτηρίων για ζώνη Α, data[1] αριθμος εισιτηρίων για ζώνη Β, data[2] αριθμος εισιτηρίων για ζώνη Γ, data[3] αριθμος εισιτηρίων για ζώνη Δ
		int card_id; //κωδικός card	
		//int t_all; //ο συνολικός χρόνος εξυπηρέτησης του πελάτη

	}pelatis;


void reservation(pelatis* p) //για δημιουργία παραγγελίας
	{
	      //ch: επιλογή ζώνης ,count2: μετρητής εισιτηρίων 
              int ch,count2;
	      //num:αριθμός εισιτηρίων,card->card_id
	      int num,card;
              int pososto_card; 
			   printf("make data :0-A , 1-B , 2-C, 3-D ,-1 finished:");
               count2 = 0;//αρχικοποίηση count2
               scanf("%d",&ch);
	       //είσοδος δεδομένων μέχρις ότου ο χρήστης πατήσει -1
	      
            while (ch != -1)
               {
            
				printf("\nhow many seats ");
				match(ch);//αντιστοίχηση ακεραίου εισόδου χρήστη σε αλφαριθμητικό για τη ζώνη
				scanf("%d",&num);//πόσες θέσεις θέλει
				//πρόσθεση του αριθμού που θέλει στη ζώνη στο συνολικό μετρητή
				if (ch == 0){
					count2 += num;
					if (count2 == 100 || count2>100) printf("zone A is full");
					else {
							if (num > 4)
							{
								printf("you cannot take more seats than 4.try again");
								do//έλεγχος ότι δεν μπορεί να κρατήσει πάνω από 4 θέσεις
								{          count2 -=num;
											scanf("%d",&num);
								}while (num <= 4); 
								count2+=num;
							}
							p->data[ch] += num ;//αντιστοίχηση της εισόδου του χρήστη στο struct
							printf ("count2 = %d \n",count2);
                            printf("OK ! \n");
                            
								
							
						}
						
				}
			   
				if (ch == 1){
					count2 += num;
					if (count2 == 130 || count2>130) printf("zone B is full");
					else {
							if (num > 4)
							{
								printf("you cannot take more seats than 4.try again");
								do//έλεγχος ότι δεν μπορεί να κρατήσει πάνω από 4 θέσεις
								{          count2 -=num;
											scanf("%d",&num);
								}while (num <= 4); 
								count2+=num;
							}
							p->data[ch] += num ;//αντιστοίχηση της εισόδου του χρήστη στο struct
							printf ("count2 = %d \n",count2);
                            printf("OK ! \n");
                            
								
							
						}
						
				}
			   
				if (ch == 2){
					count2 += num;
					if (count2 == 180 || count2>180) printf("zone C is full");
					else {
							if (num > 4)
							{    
								printf("you cannot take more seats than 4.try again");
								do//έλεγχος ότι δεν μπορεί να κρατήσει πάνω από 4 θέσεις
								{          count2 -=num;
											scanf("%d",&num);
								}while (num <= 4); 
								count2+=num;
							}
							p->data[ch] += num ;//αντιστοίχηση της εισόδου του χρήστη στο struct
							printf ("count2 = %d \n",count2);
                            printf("OK ! \n");
                            
								
							
						}
						
				}
			   
				if (ch == 3){
					count2 += num;
					if (count2 == 230 || count2>230) printf("zone D is full");
					else {
							if (num > 4)
							{
								printf("you cannot take more seats than 4.try again");
								do//έλεγχος ότι δεν μπορεί να κρατήσει πάνω από 4 θέσεις
								{          count2 -=num;
											scanf("%d",&num);
								}while (num <= 4); 
								count2+=num;
							}
							p->data[ch] += num ;//αντιστοίχηση της εισόδου του χρήστη στο struct
							printf ("count2 = %d \n",count2);
                            printf("OK ! \n");
                            
								
							
						}
						
				}
				}//while
				printf("Doste ton arithmo ths pistotikhs sas kartas");
				scanf("%d",&card);
				do{
				pososto_card = rand() % (100) +1;//η rand θα δώσει έναν τυχαίο αριθμό απο το 1 μέχρι το 100 
				if(pososto_card >= 1 && pososto_card <= 10 )   
	                   {  
						printf("h pistotikh karta den einai egyrh\n");
					    printf("Parakalw prospathiste ksana\n");
						scanf("%d",&card);
						}
				else
					{
						printf("h pistotikh karta einai eggyrh\n");
					}	
				}while (pososto_card > 10);
				p->card_id=card;
	}
		
    
	
	void match(int s)//αντιστοίχηση ακεραίου στη ζώνη
	{
    
	if (s == 0) printf ("A ");//0-
     if (s == 1) printf ("B ");//1-
     if (s == 2) printf ("C ");//2-
	 if (s == 3) printf ("D ");//2-
     printf ("would you like:");
     return;
	 
	}

