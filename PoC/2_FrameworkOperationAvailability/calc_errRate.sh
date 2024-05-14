#!/bin/bash

SET=$(seq 1 5)
for i in $SET
do
  ../Utils/calc_reqTime ./KetiFaaS/jtl/request_40000/serviceReqTime_$i.jtl 40000 300
done

index=1
sum=0

while [ $index -le 5 ]
do
  err=`../Utils/calc_reqTime ./KetiFaaS/jtl/request_40000/serviceReqTime_$index.jtl 40000 300 | awk '{print $3}'`
  res=${err:7:5}
  sum=`echo $sum + $res | bc`
  index=`expr $index + 1`
done

totalErr=`echo "scale=2;$sum / 5" | bc -l`
successRate=`echo "scale=2;100 - $totalErr" | bc -l`

echo "----------------------------------------------------------------"
echo "|                                                              |"
echo "|                  Total of error rate is $sum %                |"
echo "|                 Average of error rate is $totalErr %               |"
echo "|    Framework Operation Availability Success Rate is $successRate %  |"
echo "|                                                              |"
echo "----------------------------------------------------------------"

