output "vm_public_ips" {
  description = "Public IPs of all VMs"
  value = {
    for name, vm in azurerm_linux_virtual_machine.vm : name => vm.public_ip_address
  }
}

output "ssh_commands" {
  description = "SSH commands for each VM"
  value = {
    for name, vm in azurerm_linux_virtual_machine.vm :
    name => "ssh ${var.admin_username}@${vm.public_ip_address}"
  }
}

output "ansible_inventory" {
  description = "Ansible inventory in INI format"
  value = join("\n", concat(
    ["[app]"],
    [for name, vm in azurerm_linux_virtual_machine.vm : 
      "${name} ansible_host=${vm.public_ip_address} ansible_user=${var.admin_username}"
      if name == "app"],
    ["", "[monitoring]"],
    [for name, vm in azurerm_linux_virtual_machine.vm : 
      "${name} ansible_host=${vm.public_ip_address} ansible_user=${var.admin_username}"
      if name == "monitoring"],
  ))
}