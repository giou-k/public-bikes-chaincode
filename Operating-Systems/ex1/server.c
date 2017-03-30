//ΓΕΩΡΓΙΟΥ ΚΩΝΣΤΑΝΤΙΝΟΣ 5204

#include "first.h"
#include <time.h>
#include <string.h>

//αρχικοποίηση σταθερών 
const int Nthl = 10; 	
const int Nbank = 4;
const int t_seatfind = 6;
const int t_cardcheck = 2;
const int t_wait = 10;
const int ttransfer = 30;

int thl_sem[10];		//semaphores για τηλεφωνητές
int bank_sem[4];		//semaphores για τράπεζες


int *bank_free;  //δεικτης για τραπεζα
int *thl_free;  //δεικτης για τηλεφωνητές
int *count1;	//δείκτης για shared mem.

#define SHM_SIZE 5
#define SHMSIZE 10
#define SHMSIZ 4
#define SHM_KEY IPC_PRIVATE
#define SOCK_PATH "/tmp/echo_socket"

void sig_chld( int signo );    //συνάρτηση για εξαφάνιση φαντασμάτων παιδιών 
void kill_server();     //συνάρτηση για χειρισμό σημάτων SIGINT 




int main(void)
{


//**********************SHAREDMEMORY*******************


//για χρήση στη for
int i;


//shared memory για τηλεφωνητές


key_t thl_key = 9876;//κλειδί για shared memory

//Δημιουργία shared memory 
int thl_shm = shmget(thl_key, SHMSIZE, 0600 | IPC_CREAT);


//έλεγχος για την περίπτωση αποτυχίας δημιουργίας shared memory, και έξοδος
if ( thl_shm < 0 ) 
{
        printf("---> Could not create shared memory!\n");
        exit(1);
}


//αποθήκευση shared memory στη πρώτη εύκαιρη θέση του thl_free
thl_free = (int *)shmat( thl_shm, NULL, 0 );

for(i=0; i<10; i++) * (thl_free + i) = 0;//αρχικοποίηση με 0 για να είναι όλοι ελεύθεροι


		
//shared memory για τράπεζα


key_t bank_key = 8888;//κλειδί για shared memory

//Δημιουργία shared memory 
int bank_shm = shmget(bank_key, SHMSIZ, 0600 | IPC_CREAT);


/*έλεγχος για την περίπτωση αποτυχίας δημιουργίας shared memory, και έξοδος*/
if ( bank_shm < 0 ) 
{
        printf("---> Could not create shared memory!\n");
        exit(1);
}
		
//αποθήκευση shared memory στη πρώτη εύκαιρη θέση του bank_free
bank_free = (int *)shmat( bank_shm, NULL, 0 );

for(i=0; i<10; i++) * (bank_free + i) = 0;//αρχικοποίηση με 0 για να είναι όλες ελεύθερες



key_t my_key = 1234;//κλειδί για shared memory

//Δημιουργία shared memory 
int my_shm = shmget(my_key, SHM_SIZE, 0600 | IPC_CREAT);

//έλεγχος για την περίπτωση αποτυχίας δημιουργίας shared memory, και έξοδος
if ( my_shm < 0 ) 
{
	printf("---> Could not create shared memory!\n");
	exit(1);
}
//αποθήκευση της shared memory στη 1η θέση του count1 
count1 = (int *)shmat( my_shm, NULL, 0 );
*(count1+0) = 0;

signal( SIGINT, kill_server );


int j = 0;
int count2 = 0;
int pid;

 //δήλωση file descriptors που επιστρέφονται από τη κλήση 
 //της συνάρτησης socket και accept αντίστοιχα 
 int listenfd, connfd;
 
 //δήλωση μηκών διευθύνσεων του client και server αντίστοιχα 
 int clientlen, serverlen;
 
 //δήλωση διευθύνσεων server και client αντίστοιχα 							
 struct sockaddr_un serveraddr, clientaddr;
    
//δήλωση buffer όπου αποθηκεύονται χαρακτήρες για χρήση
//στις συναρτήσεις συστήματος read και write
char buff[5];

//δήλωση id διεργασίας της child process
pid_t childpid;

//δημιουργία semaphore για κάθε τηλεφωνητή 
for(j=0; j<Nthl; j++)
	{
	if((thl_sem[j] = semget(IPC_PRIVATE,1,PERMS | IPC_CREAT)) == -1){
		printf("\n can't create mutex semaphore %d",j);
		exit(1);
	}
	sem_create(thl_sem[j],1);
	
	}
	
//δημιουργία semaphore για κάθε τράπεζα 
for(j=0; j<4; j++)
    {
    if((bank_sem[j] = semget(IPC_PRIVATE,1,PERMS | IPC_CREAT)) == -1)
    {
    printf("\n can't create mutex semaphore %d",j);
    exit(1);
    }
    sem_create(bank_sem[j],1);

}


//ειδοποίηση στη διεργασία πατέρας ότι η διεργασία παιδί εχει σταματήσει ή έχει κανει exit       
signal( SIGCHLD, sig_chld );    


//δημιουργία του end point του server 
if ((listenfd = socket(AF_UNIX, SOCK_STREAM, 0)) == -1) {
    perror("socket");
    exit(1);
}


//καθορισμός του τύπου του socket σε local (unix domain)
serveraddr.sun_family = AF_UNIX;    

//καθορισμός του ονόματος αφτού του socket        
strcpy(serveraddr.sun_path,SOCK_PATH);     

//σβήσιμο οποιουδήποτε προηγούμενου socket με το ίδιο όνομα 
        unlink(serveraddr.sun_path);                


//συνολικό μήκος διεύθυνσης server 
serverlen = strlen(serveraddr.sun_path) + sizeof(serveraddr.sun_family);

//έλεγχος σύνδεσης socket descriptor με έν local port και εκτύπωση μηνύματος σφάλματος σε περίπτωση σφάλματος
if (bind(listenfd, (struct sockaddr *)&serveraddr, serverlen) == -1) {
    perror("bind");
    exit(1);
}

//δημιουργία μια λίστας αιτήσεων για τους clients με μήκος LISTENQ
if (listen(listenfd, LISTENQ) == -1) {
    perror("listen");
    exit(1);
}




        //*ατέρμονος βρόγχος που περιέχει τον κώδικα με τον οποίο γίνεται η σύνδεση με το client για την εξυπηρέτηση.*// 


     
for(;;) {


    //σφάλμα κατά τη διαδικασία ανάγνωσης και γραφής από και στον client
    int block;
    int str_length;

    printf("Waiting for a connection...\n");
    //καθορισμός μεγέθους διεύθυνσης του client 
    clientlen = sizeof(clientaddr);

    //αντιγραφή της επόμενης αίτησης από την ουρά αιτήσεων στη μεταβλητή connfd και διαγραφή της από την ουρά 
    connfd = accept(listenfd, (struct sockaddr*)&clientaddr, &clientlen); 
	if (connfd == -1){
	        perror("accept");
	        exit(1);
	}
	

	//δημιουργία 2 διεργασιών, μια για εξυπηρέτηση και μία για κλείσιμο σύνδεσης 
	pid = fork();


	//εκτέλεση child 
	if(pid == 0)
	{
		
		
		   printf("Connected PELATIS:%d--\n",*(count1+0));
		//close(listenfd);///////

		//μέγεθος μηνύματος client 
		int str_len = read(connfd, buff, sizeof(buff));
	
		printf("arithmos eishthriwn gia zwnh A %d \n",buff[0]);
		printf("arithmos eishthriwn gia zwnh B %d \n",buff[1]);
		printf("arithmos eishthriwn gia zwnh C %d \n",buff[2]);
		printf("arithmos eishthriwn gia zwnh D %d \n",buff[3]);
		printf("arithmos credit card %d \n",buff[4]);
		//casting buff[0] buff[1] buff[2] buff[3] buff[4] σε p0 p1 p2 p3 p4 αντίστοιχα
		char p0 = buff[0];
		   char p1 = buff[1];
		   char p2 = buff[2];
		char p3 = buff[3];
		char p4 = buff[4];
		
		int status;//gia wait pid
		//δηλώσεις μεταβλητών για fork   
		pid_t my_f,endID;
		pid_t thl_f;
		pid_t bank_f;
		my_f = fork();
		if(my_f == 0) 
		{
		
			//έλεγχος creditcard
			bank_f = fork();
			if (bank_f == 0)
			{	
				check_creditcard (t_cardcheck,count2);
			}
			
			thl_f = fork();
			if (thl_f == 0)
			{
				theater_seats(t_seatfind,count2);
			}	
		
		}
		
		else{
			
			for(j = 0; j < 20; j++) {
		   		endID = waitpid(my_f, &status, WNOHANG|WUNTRACED);

		   		if (endID == -1) {            /* error calling waitpid  */
		     		 perror("waitpid error");
		      		exit(EXIT_FAILURE);
		   		}//endif

		   		else if (endID == 0) sleep(1);

				else break;
			}//endfor
			
		}
		//κλείσιμο σύνδεσης και αύξηση του count
		close(connfd);
		count2++;
	}
}//END OF FOR(;;)


//κώδικας sig_child
void sig_chld( int signo )
{
       pid_t pid;
       int stat;
       while ( ( pid = waitpid( -1, &stat, WNOHANG ) ) > 0 ) {
              printf( "Child %d terminated.\n", pid );
       }
}


//κώδικας creditcard
void check_creditcard(int t,int count2)
{
			int i;
			int flag = 0;
			printf("ELEGXOS PISTOTIKHS KARTAS \n");

			//for για να δεσμευσουμε τραπεζικό σύστημα
	for(i = 0 ;i<4 ;i++)
  	{
		//έλεγχος για το αν είναι ελεύθερο καποιο τραπεζικό συστημα
		if (*(bank_free+i) == 0)
		{
			flag =1;
			*(bank_free+i) = 1;//Δεν ειναι ελεύθερο
			printf("bank[%d] reserved\n",i);
			//sem_wait για να μην μπορεί να χρησιμοποιηθεί από άλλη διεργασία
			sem_wait(bank_sem[i]);
			sleep(t);//t_cardcheck
			sem_signal(bank_sem[i]);//ελευθέρωση semaphore
			*(bank_free+i) = 0;//ελευθέρωση και από τον δείκτη
			printf("bank[%d] free\n",i);
			return;
		}
		else {
			if (i==3 && flag ==0) 
			{
				//printf("No free bank system\n");
				sleep(1);
				
				i=0;
			}
				}
	}//endfor

}


void theater_seats(int t,int count2)
{


	int i;
	int flag = 0;
 	printf("ELEGXOS ELEUTHERWN THESEWN \n");
	//for για να δεσμευσουμε τηλεφωνητές
    for(i = 0; i<10; i++)
    {
		//έλεγχος για το αν είναι ελεύθερο καποιος τηλεφωνητης
		if (*(thl_free+i) == 0)
		{
			flag =1 ;
				
			*(thl_free+i) = 1;  //Δεν ειναι ελεύθερο
			printf("thlefwnhths[%d] reserved\n",i);
			//sem_wait για να μην μπορεί να χρησιμοποιηθεί από άλλη διεργασία
			sem_wait(thl_sem[i]);
			
			sleep(t);  //t_seatfind
			//ελευθέρωση semaphore
			sem_signal(thl_sem[i]);
			//ελευθέρωση και από τον δείκτη
			*(thl_free+i) = 0;
			printf("thelfwnhths[%d] free\n",i);
			return;
		}
		else {
			if (i==9 && flag ==0) 
			{
				//   printf("KANENAS ELEUTHEROS THLEFWNHTHS\n");
				sleep(1);
		
				i=0;
			}
		}
	}




}

void kill_server() /*συνάρτηση εξουδετέρωσης πατέρα*/
{
	signal( SIGINT, kill_server );

	shmctl(bank_shm, IPC_RMID, NULL); /*διαγραφή κοινής μνήμης*/
	printf("** ! ** SERVER PROCESS WAS TERMINATED!\n");

	/*εξουδετέρωση διεργασίας*/
	exit(1);
}

}


                                              
