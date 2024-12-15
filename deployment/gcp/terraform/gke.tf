resource "google_service_account" "gke_access" {
  account_id   = "gke-access"
  display_name = "GKE access"
}

resource "google_container_cluster" "gke" {
  name     = "gke-cluster"
  location = var.region

  initial_node_count = 1
  remove_default_node_pool = true
}

resource "google_container_node_pool" "gke_node_pool" {
  name       = "gke-node-pool"
  location   = var.region
  cluster    = google_container_cluster.gke.name
  node_count = 1

  autoscaling {
    total_min_node_count = 1
    total_max_node_count = 3
  }

  node_config {
    preemptible  = true
    machine_type = "e2-standard-2"
    disk_size_gb = 20

    service_account = google_service_account.gke_access.email
    oauth_scopes    = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
  }
}

resource "google_compute_address" "load_balancer_1" { // gitlab
  name          = "load-balancer-1"
  ip_version    = "IPV4"
  address_type  = "EXTERNAL"
  region        = var.region
}

resource "google_compute_address" "load_balancer_2" { // mongo
  name          = "load-balancer-2"
  ip_version    = "IPV4"
  address_type  = "EXTERNAL"
  region        = var.region
}

resource "google_compute_address" "load_balancer_4" { // ingress
  name          = "load-balancer-4"
  ip_version    = "IPV4"
  address_type  = "EXTERNAL"
  region        = var.region
}

resource "google_compute_address" "load_balancer_5" {
  // mongo_2
  name         = "load-balancer-5"
  ip_version   = "IPV4"
  address_type = "EXTERNAL"
  region       = var.region
}
