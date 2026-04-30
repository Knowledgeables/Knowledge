variable "resource_group_name" {
  description = "The name of the resource group"
  type        = string
}

variable "location" {
  description = "The Azure location where the resources will be deployed"
  type        = string
}

variable "subscription_id" {
  description = "The subscription ID for the Azure account"
  type        = string

}

variable "vm_size" {
  description = "The size of the virtual machine"
  type        = string
  default     = "Standard_B2ats_v2"
}

variable "admin_username" {
  description = "The admin username for the virtual machine"
  type        = string
  default     = "azureuser"
}

variable "os_disk_caching" {
  description = "The caching setting for the OS disk"
  type        = string
  default     = "ReadWrite"
}

variable "os_disk_storage_account_type" {
  description = "The storage account type for the OS disk"
  type        = string
  default     = "Standard_LRS"
}

variable "source_image_publisher" {
  description = "The publisher of the source image"
  type        = string
  default     = "Canonical"
}

variable "source_image_offer" {
  description = "The offer of the source image"
  type        = string
  default     = "0001-com-ubuntu-server-jammy"
}

variable "source_image_sku" {
  description = "The SKU of the source image"
  type        = string
  default     = "22_04-lts-gen2"
}

variable "source_image_version" {
  description = "The version of the source image"
  type        = string
  default     = "latest"
}

variable "ssh_public_key" {
  description = "The path to the SSH public key"
  type        = string
  default     = "~/.ssh/id_rsa.pub"
}

variable "allowed_ssh_cidrs" {
  type        = list(string)
  default     = ["0.0.0.0/0"]
  description = "CIDRs allowed to reach SSH. Open to world by default for student project (dynamic home IP, fail2ban + key-only auth provide protection). Restrict in production."
}

variable "vms" {
  description = "Map of VMs to create"
  type = map(object({
    role = string
  }))
  default = {
    app = {
      role = "application"
    }
    monitoring = {
      role = "monitoring"
    }
  }
}