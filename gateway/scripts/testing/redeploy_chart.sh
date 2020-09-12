# Do not use except for testing
echo "Removing chart"
helm delete vscode 
# echo "Deleting Nginx Ingress"
# helm delete nginx-ingress
# echo "Installing Nginx Ingress"
# helm install nginx-ingress nginx-stable/nginx-ingress

echo "Installing chart"
helm install vscode ../.././charts/vscode

