# kubernetes-assign

# kubebuilder deploy
make deploy IMG=akankshakumari393/step-kubebuilder
kustomize build config/kustomize-sample | kubectl apply -f -

# kube-operator docker image
make deploy IMG=akankshakumari393/step-operator