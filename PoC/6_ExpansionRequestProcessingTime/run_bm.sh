#!/bin/sh

/root/PrivateCloudContainer/PCC-5th-PoC/Utils/JMeter/apache-jmeter-5.1/bin/jmeter -n -t ./KetiFaaS/jmx/KetiFaaS-total_request_forever.jmx -l ./KetiFaaS/jtl/request_forever/result.jtl -j ./KetiFaaS/log/request_forever/result.log
