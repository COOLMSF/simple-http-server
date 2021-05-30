#!/bin/bash

sum = 0
for ((i=0;i<=100;i+=2))
do
  echo $i
  let sum=$sum+$i
done

echo $sum
