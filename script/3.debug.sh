#/bin/bash
NAME=$(kubectl get pod -n keti-system | grep -E 'cluster-api-server' | awk '{print $1}')

#echo "Exec Into '"$NAME"'"

#kubectl exec -it $NAME -n $NS /bin/sh
kubectl logs -f $NAME -n keti-system