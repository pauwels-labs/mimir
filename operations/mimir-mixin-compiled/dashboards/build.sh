for FILE in *.json; do
    kubectl create configmap "${FILE%%.json}" -oyaml --dry-run=client -n i-grafana --from-file "$FILE" | yq w - 'metadata.labels"pauwelslabs.com/dashboards"' "true" --style=double | yq d - 'metadata.creationTimestamp' > yaml/configmap-${FILE%%.json}.yaml
done
