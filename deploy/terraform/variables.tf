variable "DigitalOceanToken" {
  description = "Digital Ocean Secret. Do not pass via this file"
  type        = string
  sensitive   = true
}

variable "CloudflareToken" {
  description = "Cloudflare Api Token. Do not pass via this file"
  type        = string
  sensitive   = true
}
