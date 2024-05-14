#!/bin/bash

SET=$(seq 1 5)
for i in $SET
do
  ../Utils/calc_reqTime ./KetiFaaS/jtl/request_40000/serviceReqTime_$i.jtl 20000 300
done

index=1
sum=0

while [ $index -le 5 ]
do
  avg=`../Utils/calc_reqTime ./KetiFaaS/jtl/request_40000/serviceReqTime_$index.jtl 20000 300 | awk '{print $6}'`
  res=${avg:1:6}
  sum=`echo $sum + $res | bc`
  index=`expr $index + 1`
done

totalAvg=`echo "scale=2;$sum / 5" | bc -l`

echo "----------------------------------------------------------------"
echo "|                                                              |"
echo "|         Total of service response times is $sum ms          |"      
echo "|        Average of service response times is $totalAvg ms        |"
echo "|                                                              |"
echo "----------------------------------------------------------------"
