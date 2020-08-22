resource "digitalocean_database_cluster" "redis" {
    name = "redis"
    engine = "redis"
    version = "5"
    size = "db-s-1vcpu-1gb"
    region = "nyc3"
    node_count =1
    private_network_uuid = digitalocean_vpc.public.id

}