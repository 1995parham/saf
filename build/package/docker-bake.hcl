group "default" {
  targets = [
    "koochooloo"
  ]
}

variable "TAG" {
  default = "latest"
}

target "koochooloo" {
  dockerfile = "build/package/Dockerfile"
  tags = [
    "ghcr.io/1995parham/saf:${TAG}"
  ]
}
