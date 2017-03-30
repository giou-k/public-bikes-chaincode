//ΓΕΩΡΓΙΟΥ ΚΩΝΣΤΑΝΤΙΝΟΣ 5204
//ΒΑΣΙΛΗΣ ΛΙΝΑΡΔΟΣ 5016

//simperilipsi header file
#include "first2.h"
char buffer[5];//buffer για αποθήκευση της εισόδου από το χρήστη

    int main( int argc,char *argv[])
    {

	printf("######### Theatro EMPROS #########\n");

   
        int s, t, len,connection_established;
        struct sockaddr_un remote;     //δήλωση διεύθυνσης του server

        //αρχικοποίηση socket discriptor και έλεγχος για τη δημιουργία του
        if ((s = socket(AF_UNIX, SOCK_STREAM, 0)) == -1) {
            perror("socket");
            exit(1);
        }


        //ορισμός οικογένειας socket που χρησιμοποιείται για την ίδια μηχανή 
        remote.sun_family = AF_UNIX;
        //αντιγραφή της διεύθυνσης του socket στη διεύθυνση του server 
        strcpy(remote.sun_path, SOCK_PATH);
        //καθορισμός του συνολικού μήκους διεύθυνσης 
        len = strlen(remote.sun_path) + sizeof(remote.sun_family);
       
        //με τη κλήση του συστήματος connect() συνδέεται o client με τον server 
        //διαφορετικά εκτυπώνεται μήνυμα λάθους 
		connection_established = connect(s, (struct sockaddr *)&remote, len);
        if (connection_established ==-1) {
            perror("connect");
            exit(1);
        }
        else{
		
            printf("Connected.\n");
	   
	  
            printf("create new customer\n");
	    int choose;
			int pososto;
			printf("epelexe 1 gia dikes sou times h' 0 gia tyxaies \n") ;
			scanf("%d",&choose);
			getchar();
		
			if (choose == 1){
		
					printf("create new customer\n");
					int i,count2;//για for
					char buffer[5];//buffer για αποθήκευση της εισόδου από το χρήστη
					
					if (argc == 1)//για ένα όρισμα εισόδου 
					{
						pelatis p1;//δημιουργία πελάτη             
						for (i=0;i<4;i++ ) p1.data[i]=0;//αρχικοποίηση πίνακα 
										
						
						reservation(&p1);//κλήση συνάρτησης για δημιουργία κράτησης
						
						//casting toy p1.zwnh[] apo struct pelatis sto data[]
						buffer[0]=(char)p1.data[0];
						buffer[1]=(char)p1.data[1];
						buffer[2]=(char)p1.data[2];
						buffer[3]=(char)p1.data[3];
						buffer[4]=(char)p1.card_id;
						
					}	
				
				
					if (argc > 1)//σε περίπτωση πολλών ορισμάτων
					{
					
						int w = 0; //δείχνει θέση στο buffer 
						for(count2=1; count2<argc; count2++)
						{
							//αποθήκευση στο buffer των ορισμάτων εισόδου
							if (strcmp(argv[count2],"0")==0) buffer[w]=0;
							if (strcmp(argv[count2],"1")==0) buffer[w]=1;
							if (strcmp(argv[count2],"2")==0) buffer[w]=2;
							if (strcmp(argv[count2],"3")==0) buffer[w]=3;
							if (strcmp(argv[count2],"4")==0) buffer[w]=4;
							
							w++;//αύξηση κατά 1
						}
					
					}//if argc>1
			}//choose == 1	
			
			
			else if (choose == 0)//random τιμές
					{
						
						int b = rand() % (4)+1 ; //ο buffer θα πάρε τυχαία μία τιμή απο το 1 έως το 4 
						buffer[0] = (char)b; 
						pososto = rand() % (100)+1; //η μεταβλητή pososto θα πάρει μια τιμή από το 1-100 
						if (pososto <= 10 && pososto >= 1)
							buffer[1]='A';
						else if (pososto<=30 && pososto>=11)
							buffer[1] ='B';
						else if (pososto<=60 && pososto>=31)
							buffer[1]='C';
						else if (pososto<=100 && pososto>=61)
							buffer[1]='D';
					}
			
			//γράψιμο στο server της κράτησης
			if (write(s, buffer, sizeof(buffer)) == -1) {
                perror("write");
                exit(1);
            }

        }

		close(s);//κλείσιμο σύνδεσης 
		
	}
