#!/bin/sh

SET=$(seq 1 5)
for i in $SET
do
  /root/PrivateCloudContainer/PCC-5th-PoC/Utils/JMeter/apache-jmeter-5.1/bin/jmeter -n -t ./KetiFaaS/jmx/KetiFaaS-total_request_40000.jmx -l ./KetiFaaS/jtl/request_40000/serviceReqTime_$i.jtl -j ./KetiFaaS/log/request_40000/serviceReqTime_$i.log
  sleep 1m
done

