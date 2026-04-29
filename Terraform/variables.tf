variable "subscription_id" {
  description = "Azure subscription ID"
  type        = string
}

variable "vm_name" {
  description = "The name of the virtual machine"
  type        = string
  default     = "main-vm"
}