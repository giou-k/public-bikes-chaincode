//ΓΕΩΡΓΙΟΥ ΚΩΝΣΤΑΝΤΙΝΟΣ 5204
//ΒΑΣΙΛΗΣ ΛΙΝΑΡΔΟΣ 5016

#include "first2.h"   //voithitiko header gia synarthseis
#include <time.h>    //header gia synarthseis time
#include <string.h>  //header gia synarthseis string
#include <pthread.h> //header gia synarthseis thread
#include <sys/time.h>//header gia synarthseis xronou

//αρχικοποίηση σταθερών 
const int Nthl = 10; 	
const int Nbank = 4;
const int t_seatfind = 6;
void *pmytimer(int);
const int t_cardcheck = 2;
const int twait = 10;
const int ttransfer = 30;
void theater_seats(int);
void check_creditcard(int);
void timer_handler(int);
int list[20]={0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0};
int anamonh=0; //katastash anamonhs

pthread_mutex_t mutex = PTHREAD_MUTEX_INITIALIZER; //gia diavasma krathshs apo socket
pthread_mutex_t thl_free[10] = PTHREAD_MUTEX_INITIALIZER; //mutex gia thlefwnhtes
pthread_mutex_t bank_free[4] = PTHREAD_MUTEX_INITIALIZER; //mutex gia trapezas

void sig_chld( int signo );    //συνάρτηση για εξαφάνηση φαντασμάτων παιδιών 
void kill_server();     //συνάρτηση για χειρισμό σημάτων SIGINT 

char buff[5];			//gia diavasma krathshs apo socket
int krathsh=0;	//arithmos krathshs


void threads( int p)  //synarthsh gia etoimasia krathshs 
			       //apo nhma ths main
{
	
    printf( "PARAGGELIA No: %d!\n", p );
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
	

 	int  tmp[20];		//gia join nhmatwn mthread[20]
     pthread_t mthread[20];	//threads gia parallhles leitourgies
     pthread_attr_t attr;	//attribute sthn arxikpoihsh thread
  
     pthread_attr_init( &attr );//arxikopoihsh tou attribute
     pthread_attr_setdetachstate( &attr, PTHREAD_CREATE_JOINABLE );

	//create thread for phone company with specific time
	tmp[0] = pthread_create( &mthread[0],&attr, theater_seats, t_seatfind);
	//se periptosi apotyxias
     	if ( tmp[0] != 0 )
             {
                     fprintf(stderr,"Creating thread 0 failed!");
                     return 1;
             }
	//create thread for bank with specific time
	tmp[1] = pthread_create( &mthread[1],&attr, check_creditcard, t_cardcheck);
	//se periptosi apotyxias
     	if ( tmp[1] != 0 )
             {
                     fprintf(stderr,"Creating thread 1 failed!");
                     return 1;
             }
	
	//join
	tmp[0] = pthread_join( mthread[0], NULL );
	tmp[1] = pthread_join( mthread[1], NULL );


	printf("EGINE H KRATHSH: %d ,ANAMONH:%d \n",p,anamonh);
	printf("RESERVATION: %d\n",p);
	list[p] =anamonh;// stin thesi i apothikeuei thn kathysterisi ths krathshs i
        pthread_exit( NULL );//telos synarthshs dhmiourgeias thread apo main
}
pthread_attr_t attr;//attribute gia nhma
pthread_mutex_t mutex = PTHREAD_MUTEX_INITIALIZER;//arxikpoihsh attribute
               



int main(void)
{
 	pthread_attr_init( &attr );// arxikopoihsh 
    pthread_attr_setdetachstate( &attr, PTHREAD_CREATE_JOINABLE );
	pthread_t mythread;//elegxos kathusterhshs
	int mytmp = pthread_create( &mythread,&attr, pmytimer, t_seatfind);//dhmiourgeia thread

	//arxikopoihsh synarthshs xeirismou shmatos SIGINT 
	signal( SIGINT, kill_server );

        //dilosi file descriptors pou epistrefontai apo tin klisi 
        //tis sinartisis socket kai accept antistoixa
        int listenfd, connfd;
        //dilosi mikon dieuthinseon tou client kai server antistoixa
        int clilen, len;
        //dilosei dieuthiseon server kai client antistoixa
        struct sockaddr_un servaddr, cliaddr;
        //dimiourgia tou end point tou server
        if ((listenfd = socket(AF_UNIX, SOCK_STREAM, 0)) == -1) {
            perror("socket");
            exit(1);
        }

        servaddr.sun_family = AF_UNIX;            //kathorismos tou tipou tou socket se local (Unix Domain)
        strcpy(servaddr.sun_path, SOCK_PATH);     //kathorismos tou onomatos aftou tou socket
        unlink(servaddr.sun_path);                //svisimo opoioudipote proigoumenou socket me to idio filename
        //sinoliko mikos dieuthinseos server
        len = strlen(servaddr.sun_path) + sizeof(servaddr.sun_family);

        //elegxos sindesis socket descriptor me ena local port kai ektiposi
        //minimatos sfalmatos se periptosi sfalmatos
        if (bind(listenfd, (struct sockaddr *)&servaddr, len) == -1) {
            perror("bind");
            exit(1);
        }

        //dimiourgia mias listas aitiseon gia tous clients me mikos LISTENQ
        if (listen(listenfd, LISTENQ) == -1) {
            perror("listen");
            exit(1);
        }
        /*atermonos vrogxos pou periexei ton kodika me ton opoio ginetai i sindesi 
          me ton client kai eksipiretite.*/ 
 
        for(;;) {

            printf("Waiting for a connection...\n");
            //kathorismos megethous dieuthinsis tou client
            clilen = sizeof(cliaddr);

            //antigrafi tis epomenis aitisis apo tin oura aitiseon 
            //sti metavliti connfd kai diagrafi tis apo tin oura
            connfd = accept(listenfd, (struct sockaddr*)&cliaddr, &clilen); 
		if (connfd == -1){
                perror("accept");
                exit(1);
            	}

		int  tmp; //gia arxikpoihsh nhmatos kai join
        	pthread_t thread;//nhma gia etoimasia krathshs

                // creating thread gia etoimasia krathshs
       	tmp = pthread_create( &thread, &attr,threads, krathsh);
	krathsh++;//auksisi tou metrhth pelatwn
	if ( tmp != 0 )
                {
                        fprintf(stderr,"Creating thread  failed!");
                        return 1;
                }
	


}//end for

}//end main


void *pmytimer(int a)
{

struct sigaction sa;
 struct itimerval timer;

 /* Install timer_handler as the signal handler for SIGVTALRM. */
 memset (&sa, 0, sizeof (sa));
 sa.sa_handler = &timer_handler;
 sigaction (SIGVTALRM, &sa, NULL);

//arxikopoihsh timwn
 timer.it_value.tv_sec = 0;
 timer.it_value.tv_usec = 250000;

 timer.it_interval.tv_sec = 0 ;
 timer.it_interval.tv_usec =250000 ;
 /* Start a virtual timer. It counts down whenever this process is
   executing. */
 setitimer (ITIMER_VIRTUAL, &timer, NULL);
while(1);



}

void timer_handler (int signum)
{
int i;
static int end=0;
for (i=0;i<20;i++)
{
//elegxos kathysterhshs 
if (list[i]>=twait) 
printf("\nKRATHSH No:%d Mas sugxwreite gia thn kathusterhsh",i);
}
end = i;
sleep(20);
fflush(stdout);
}


void kill_server() 
{


	signal( SIGINT, kill_server );
	printf("** ! ** SERVER PROCESS WAS TERMINATED!\n");
	exit(1);
}



void check_creditcard(int t)
{
int anamonh_check;	//gia kathysterhsh
int i;		//gia for
int tryl;	//gia try_lock
for(i=0;i<4;i++)
{
	//try_lock gia elegxo diathesimothtas trapezas
	tryl = pthread_mutex_trylock(&bank_free[i]);
	if(tryl==0)//trapeza eleftheros
	{
	printf("BANK %d WORKING\n",i);
	sleep(t);
	pthread_mutex_unlock(&bank_free[i]);
	printf("BANK %d FINISHED\n",i);
	i=0;
	break;
	}
	else if(tryl==EBUSY)//trapeza mh eleftherh
	{
	//an exei ftasei sthn teleytaia trapeza kai einai 
	//oloi oi alles desmeymenes
	if(i==3) {i=0;sleep(1);anamonh_check++;}//perimenei ena deyterolepto kai 
	//prospelaunei ap'thn arxh tis trapezes
	}
}//end for
anamonh = anamonh_check+t;//prosthesi tis kathysterhshs sto check_creditcard sth synolikh kathysterhsh
}//end check_creditcard

void theater_seats(int t)
{
int anamonh_check;	//gia kathysterhsh
int i;		//gia for
int tryl;	//gia try_lock
for(i=0;i<10;i++)
{
	//try_lock gia elegxo diathesimothtas thlefwnhth
	tryl = pthread_mutex_trylock(&thl_free[i]);
	if(tryl==0)//thlefwnhths eleftheros
	{
	printf("LINE %d WORKING\n",i);
	sleep(t);
	pthread_mutex_unlock(&thl_free[i]);
	printf("LINE %d CLOSED\n",i);
	i=0;
	break;
	}
	else if(tryl==EBUSY)//thlefwnhthS mh eleftherh
	{
	//an exei ftasei ston teleytaio thlefwnhth kai einai 
	//oloi oi alloi desmeymenoi
	if(i==9) {i=0;sleep(1);anamonh_check++;}//perimenei ena deyterolepto kai 
	//prospelaunei ap'thn arxh tous thlefwnhtes
	}
}//end for
anamonh = anamonh_check+t;//prosthesi tis kathysterhshs sto check_creditcard sth synolikh kathysterhsh
}//end check_creditcard


