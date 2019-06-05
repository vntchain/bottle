#include "vntlib.h"

                 
    
                  
             
        
                      
             
        

typedef struct
{
    string ccc;
    mapping(string, string) aaa;
} S1;

KEY mapping(string, S1) b;


void key9maewdrh(){
AddKeyInfo( &b.value.ccc, 6, &b, 9, false);
AddKeyInfo( &b.value.ccc, 6, &b.key, 6, false);
AddKeyInfo( &b.value.ccc, 6, &b.value.ccc, 9, false);
AddKeyInfo( &b.value.aaa.value, 6, &b, 9, false);
AddKeyInfo( &b.value.aaa.value, 6, &b.key, 6, false);
AddKeyInfo( &b.value.aaa.value, 6, &b.value.aaa, 9, false);
AddKeyInfo( &b.value.aaa.value, 6, &b.value.aaa.key, 6, false);
}
constructor test() {
key9maewdrh();
InitializeVariables();}
