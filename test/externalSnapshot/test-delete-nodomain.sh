#!/bin/bash

virsh destroy test111
virsh undefine test111
kubectl apply -f ../vm/01-CreateVMFromISO.json

rm -rf /var/lib/libvirt/pooltest111/disktest111/snapshots
cat << EOF > /var/lib/libvirt/pooltest111/disktest111/config.json
{"current":"/var/lib/libvirt/pooltest111/disktest111/disktest111","dir":"/var/lib/libvirt/pooltest111/disktest111","name":"disktest111","pool":"pooltest111"}
EOF

go build -o /tmp/sdsctl /root/go_project/sdsctl/cmd/sdsctl/main.go
yes | cp  /tmp/sdsctl /usr/bin

kubectl delete vmdsn disktest111-snapshot2 disktest111-snapshot
sleep 5
kubectl apply -f 01-CreateExternalSnapshot.json
sleep 3
kubectl apply -f 01-CreateExternalSnapshot-2.json
sleep 2
kubectl get vmdsn