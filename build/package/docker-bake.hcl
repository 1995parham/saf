group "default" {
  targets = [
    "saf"
  ]
}

variable "TAG" {
  default = "latest"
}

target "saf" {
  dockerfile = "build/package/Dockerfile"
  tags = [
    "ghcr.io/1995parham/saf:${TAG}"
  ]
}
