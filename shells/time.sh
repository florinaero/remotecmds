#! /bin/bash

for((i=1;i<=100;i++));
    do 
    curl localhost:8080/time;
    sleep 1;
done