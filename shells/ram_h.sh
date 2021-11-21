#! /bin/bash

for((i=1;i<=1000*1000;i++));
    do 
    curl localhost:8080/ram_h;
    sleep 0.01;
done