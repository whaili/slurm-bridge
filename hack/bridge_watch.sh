#!/bin/sh

watch -n1 "\
    echo 'SLURM PODS'; \
    kubectl get pods -o wide -n slurm -l app.kubernetes.io/name=slurmd; echo; \
    echo 'SLURM BRIDGE PODS'; \
    kubectl get pods -o wide -n slurm-bridge; echo; \
    echo 'PODGROUP STATUS'; \
    kubectl get podgroup -n slurm-bridge; echo; \
    echo 'JOB STATUS'; \
    kubectl get jobs -n slurm-bridge; echo; \
    echo 'JOBSET STATUS'; \
    kubectl get jobset -n slurm-bridge; echo; \
    echo 'SINFO'; \
    kubectl exec -n slurm statefulset/slurm-controller -- sinfo; echo; \
    echo 'SQUEUE PENDING'; \
    kubectl exec -n slurm statefulset/slurm-controller -- squeue --states=pending; echo; \
    echo 'SQUEUE RUNNING'; \
    kubectl exec -n slurm statefulset/slurm-controller -- squeue --states=running ; echo; \
    echo 'SQUEUE COMPLETE'; \
    kubectl exec -n slurm statefulset/slurm-controller -- squeue  --states=BF,CA,CD,CF,CG,DL,F,NF,OOM,PR,RD,RF,RH,RQ,RS,RV,SI,SE,SO,S,TO"
