set -v
set -e
PATH=$PATH:$PWD/cluster
export NUM_MINIONS=3
kube-up.sh
gcloud compute instances list | grep minion | awk "{ IPS= IPS \",\" \$5;} END{ print IPS}" > /tmp/tp
IPLIST=`cat /tmp/tp`
IPLIST=${IPLIST:1}

### 1. Mount and format disks etc and ....

echo "IPAddresses are:" $IPLIST

DISKID=1
for hostname in `gcloud compute instances list | grep minion|cut -f 1 -d " "`; do  
 echo $hostname;
 zonename=`gcloud compute instances list $hostname | grep minion|cut -f 2 -d " "`


 diskpresent=`gcloud compute disks list data-disk-$DISKID | grep data-disk-$DISKID |wc -l`

 if [ $diskpresent == "0" ] 
  then
    gcloud compute disks create "data-disk-$DISKID" --size "512" --zone "us-central1-b" --type "pd-standard"
 fi

 gcloud compute instances attach-disk $hostname --disk data-disk-$DISKID --zone $zonename --device-name data-disk-$DISKID


 gcloud compute copy-files ./examples/ecs_deploy/com.emc.vipr.fabric.hal.conf $hostname:/tmp/ --zone $zonename
 gcloud compute copy-files ./examples/ecs_deploy/dbus_service.py $hostname:/tmp/ --zone $zonename

 gcloud compute copy-files ./examples/ecs_deploy/gce_hostprep.sh $hostname:/tmp/ --zone $zonename

 devicename=`gcloud compute ssh  $hostname  "sudo ls -ltr /dev/disk/by-id/" --zone $zonename | grep google-data-disk-$DISKID | cut -d \/ -f 3`

 gcloud compute ssh  $hostname  "sudo /tmp/gce_hostprep.sh $IPLIST $devicename" --zone $zonename 

 DISKID=$((DISKID+1))

done

echo "DONE SETTING UP first round!!"




### 2. bring up ECS with some junk config ...


done=0
kubectl.sh create -f examples/ecs_deploy/ecs-controller.yaml 
kubectl.sh create -f examples/ecs_deploy/ecs-service.yaml 

echo "Sleeping for 5 mins"

sleep 300


#should actually loop till the pods are stable
running=`kubectl.sh get pods | grep ecs | grep -i "running" | wc -l`
echo $running


### 3. Get the actual IP addresses of the containers

IPLIST=
for hostname in `gcloud compute instances list | grep minion|cut -f 1 -d " "`; do  
 echo $hostname;
 zonename=`gcloud compute instances list $hostname | grep minion|cut -f 2 -d " "`
 
 gcloud compute ssh $hostname "ps -ef  | grep boot.sh | grep -v grep" --zone $zonename | awk "{procid=\$2;} END{print procid;}" > /tmp/tp
 PROCID=`cat /tmp/tp`
 
 gcloud compute ssh $hostname "sudo nsenter -t $PROCID -m -u -i -n -p ifconfig eth0" --zone $zonename |grep 'inet addr:' | cut -d: -f2 | awk '{ print $1}' > /tmp/tp

 IP=`cat /tmp/tp`

 HN=`gcloud compute ssh $hostname "sudo nsenter -t $PROCID -m -u -i -n -p hostname" --zone $zonename`

 printf '{"private_interface_name":"eth0","public_interface_name":"eth0","hostname":"%s","public_ip":"%s"}' $HN $IP > network.json

 cat network.json
 gcloud compute copy-files ./network.json $hostname:/tmp/ --zone $zonename
# gcloud compute copy-files ./examples/ecs_deploy/binplace $hostname:/tmp/ --zone $zonename
# gcloud compute ssh $hostname "sudo cp -rf /tmp/binplace /host/files" --zone $zonename
# gcloud compute ssh $hostname "sudo nsenter -t $PROCID -m -u -i -n -p /host/files/binplace/cpbinplace.sh" --zone $zonename
 
 gcloud compute ssh $hostname "sudo cp -f /tmp/network.json /host/data/network.json" --zone $zonename

 IPLIST=$IPLIST,$IP

done
 IPLIST=${IPLIST:1}
 echo $IPLIST > seeds

echo "DONE SETTING UP first round!!"
echo "Restarting individual boot.sh s again!!"


for hostname in `gcloud compute instances list | grep minion|cut -f 1 -d " "`; do  
 echo $hostname;
 zonename=`gcloud compute instances list $hostname | grep minion|cut -f 2 -d " "`
 gcloud compute copy-files ./seeds $hostname:/tmp/ --zone $zonename
 gcloud compute ssh $hostname "sudo cp -f /tmp/seeds /host/files/seeds" --zone $zonename
 gcloud compute ssh $hostname "ps -ef  | grep boot.sh | grep -v grep" --zone $zonename | awk "{procid=\$2;} END{print procid;}" > /tmp/tp
 PROCID=`cat /tmp/tp`
 gcloud compute ssh $hostname "sudo nsenter -t $PROCID -m -u -i -n -p  service storageos-dataservice restart" --zone $zonename
done

 echo "Sleeping for 5 minutes till the situation stabilizes"

 sleep 300
 gcloud compute copy-files examples/ecs_deploy/object_prep.py $hostname:/tmp/ --zone $zonename
 gcloud compute ssh $hostname "sudo cp -f /tmp/object_prep.py /host/files/" --zone $zonename


 gcloud compute copy-files examples/ecs_deploy/license.xml $hostname:/tmp/ --zone $zonename
 gcloud compute ssh $hostname "sudo cp -f /tmp/license.xml /host/files/" --zone $zonename
 
OBJ_PREP_CMD="python object_prep.py --ECSNodes=$IPLIST --Namespace=ns1 --ObjectVArray=ova1 --ObjectVPool=ovp1 --UserName=emccode --DataStoreName=ds1 --VDCName=vdc1 --MethodName="
 echo "Prepping object deployments"

echo "cd /host/files" > setup.sh
echo $OBJ_PREP_CMD >> setup.sh
chmod a+x setup.sh


gcloud compute copy-files setup.sh $hostname:/tmp/ --zone $zonename
gcloud compute ssh $hostname "sudo cp -f /tmp/setup.sh /host/files/" --zone $zonename
 
#gcloud compute ssh $hostname "sudo /tmp/setup.sh" --zone $zonename
