#! /bin/bash

for((i=1;i<=1000;i++));
    do 
    curl localhost:8080/ram;
    sleep 3;
done