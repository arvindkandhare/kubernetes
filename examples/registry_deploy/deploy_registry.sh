#Expected: A passwd_file file with emccode:<secret key> in the current directory

set -v
set -e

for hostname in `gcloud compute instances list | grep minion|cut -f 1 -d " "`; do  
 echo $hostname;
 zonename=`gcloud compute instances list $hostname | grep minion|cut -f 2 -d " "`

 gcloud compute copy-files ./certs $hostname:/tmp/ --zone $zonename
 gcloud compute ssh $hostname "sudo cp -rf /tmp/certs /host" --zone $zonename
done

kubectl.sh create -f examples/registry_deploy/registry-controller.yaml --validate=false
kubectl.sh create -f examples/registry_deploy/registry-service.yaml

REGISTRY_IP=

while [ -z "$REGISTRY_IP"]
do
 echo "Waiting for the LB to assign an IP"
 sleep 30
 REGISTRY_IP=`kubectl.sh describe services/registry-service | grep "LoadBalancer Ingress:" | cut -f 2 -d : |  tr -d '[[:space:]]' `
done

## Create the certificate
openssl req   -newkey rsa:4096 -nodes -sha256 -keyout certs/domain.key   -x509 -days 365 -out certs/domain.crt -batch -subj "/C=US/CN=Alon"

#Copy the new ceritificate

echo "Copying the new cert to all the nodes"
for hostname in `gcloud compute instances list | grep minion|cut -f 1 -d " "`; do
 echo $hostname;
 zonename=`gcloud compute instances list $hostname | grep minion|cut -f 2 -d " "`

 gcloud compute copy-files ./certs $hostname:/tmp/ --zone $zonename
 gcloud compute ssh $hostname "sudo cp -rf /tmp/certs /host" --zone $zonename
done


#Restart the pods
echo "Restarting registry pod"
SERVICE_POD=`kubectl.sh get pods | grep "registry-" | cut -f 1 -d " "`
kubectl.sh delete po $SERVICE_POD

 
