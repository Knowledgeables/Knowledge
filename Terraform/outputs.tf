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
    ["[knowledge]"],
    [for name, vm in azurerm_linux_virtual_machine.vm : 
      "${vm.public_ip_address} ansible_user=${var.admin_username}"
      if name == "knowledge"],
    ["", "[kmonitor]"],
    [for name, vm in azurerm_linux_virtual_machine.vm : 
      "${vm.public_ip_address} ansible_user=${var.admin_username}"
      if name == "kmonitor"],
  )) 
}

output "ssh_config" {
  description = "SSH config for ~/.ssh/config"
  value = join("\n", [
    for name, ip in azurerm_public_ip.vm : <<-EOF
    Host ${name}
        HostName ${ip.ip_address}
        User ${var.admin_username}
        IdentityFile ~/.ssh/id_rsa
    EOF
  ])
}