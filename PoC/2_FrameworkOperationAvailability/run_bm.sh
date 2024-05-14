#!/bin/sh

#SET=$(seq 1 5)
#for i in $SET
#do
#  /root/PrivateCloudContainer/PCC-3rd-PoC/Utils/JMeter/apache-jmeter-5.1/bin/jmeter -n -t ./KetiFaaS/jmx/KetiFaaS-total_request_40000.jmx -l ./KetiFaaS/jtl/request_40000/serviceReqTime_$i.jtl -j ./KetiFaaS/log/request_40000/serviceReqTime_$i.log
#done

/root/PrivateCloudContainer/PCC-5th-PoC/Utils/JMeter/apache-jmeter-5.1/bin/jmeter -n -t ./KetiFaaS/jmx/KetiFaaS-total_request_40000.jmx -l ./KetiFaaS/jtl/request_40000_orig/serviceReqTime_1.jtl -j ./KetiFaaS/log/request_40000_orig/serviceReqTime_1.log

