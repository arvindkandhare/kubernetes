#Expected: A passwd_file file with emccode:<secret key> in the current directory
# Update the values in the controller file

set -v
set -e
for hostname in `gcloud compute instances list | grep minion|cut -f 1 -d " "`; do  
 echo $hostname;
 zonename=`gcloud compute instances list $hostname | grep minion|cut -f 2 -d " "`

 gcloud compute ssh $hostname "sudo rm -rf /data/passwd_file" --zone $zonename
 gcloud compute copy-files ./examples/ecs_deploy/filesvc/passwd_file $hostname:/tmp/ --zone $zonename
 gcloud compute ssh $hostname "sudo mkdir -p /data;sudo cp /tmp/passwd_file /data; sudo chmod 400 /data/passwd_file" --zone $zonename
done

kubectl.sh create -f examples/ecs_deploy/mosaicme/mosaicweb-service.yaml
kubectl.sh create -f examples/ecs_deploy/mosaicme/mosaicweb-controller.yaml 

