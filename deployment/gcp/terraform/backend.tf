terraform {
  backend "gcs" {
    bucket = "game-state"
    prefix = "terraform/state"
  }
}

provider "google" {
  project     = "884086106552"
  region      = "europe-west1"
}
