resource "digitalocean_kubernetes_cluster" "cluster" {
    name = "cluster"
    region = "nyc3"
    version = "1.18.6-do.0"
    vpc_uuid = digitalocean_vpc.public.id

    node_pool {
        name = "backend-pool"
        size = "s-2vcpu-4gb"
        node_count = 1
    }
}