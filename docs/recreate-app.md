# Infrastructure & Deployment Runbook

This document describes how to provision, configure, and deploy the Knowledge application from scratch using Infrastructure as Code (Terraform) and Configuration Management (Ansible).

## Architecture Overview

The system uses a three-layer separation of concerns:

| Layer | Tool | Purpose | Frequency |
|---|---|---|---|
| Infrastructure | Terraform | VMs, networking, IPs | Rare (manual) |
| Configuration | Ansible | OS-level setup (Docker, fail2ban) | Occasional (CI on `ansible/**` changes) |
| Application | Docker Compose + existing CD pipeline | App deployment | Continuous (every push) |

Two VMs are provisioned in Azure: `app-vm` and `monitoring-vm`. They share networking (one VNet, one subnet, one NSG with all required ports) but are configured for distinct roles via the application layer.

## Prerequisites

### Tools
- Terraform 1.9+
- Azure CLI (`az`)
- Ansible 2.16+ (Linux/WSL only — does not run natively on Windows)
- Active Azure subscription with B-series VM quota
- SSH keypair at `~/.ssh/id_rsa` (public key registered with Azure)

### Repository setup
- Clone the monorepo
- Copy `terraform/terraform.tfvars.example` to `terraform/terraform.tfvars` and fill in subscription ID
- Copy `ansible/group_vars/all.yml.example` to `ansible/group_vars/all.yml` and fill in `ansible_user`

### GitHub secrets (for CI)
- `SSH_PRIVATE_KEY` — contents of SSH private key (existing)
- `SSH_USER` — VM admin username, e.g. `azureuser` (existing)
- (Plus existing app-deployment secrets: POSTGRES_*, DATABASE_URL, etc.)

## Recreating the system from scratch

### 1. Provision infrastructure (Terraform)

```bash
cd terraform/
terraform init
terraform plan      # review what will be created
terraform apply     # confirm with 'yes'
```

Expected duration: ~2-3 minutes. Outputs include:
- `vm_public_ips` — map of VM names to public IPs
- `ssh_commands` — pre-built SSH commands for each VM
- `ansible_inventory` — ready-to-use Ansible inventory

### 2. Update Ansible inventory

```bash
terraform output -raw ansible_inventory > ../ansible/inventory.ini
```

This regenerates the inventory with current VM IPs. Commit this file to git so CI can use it.

### 3. Configure VMs (Ansible)

**Local run:**
```bash
cd ../ansible/
ansible all -i inventory.ini -m ping       # sanity check
ansible-playbook -i inventory.ini playbook.yml
```

Expected duration: ~3-5 minutes per VM. The playbook installs:
- System packages (apt update/upgrade)
- Docker + docker-compose
- Fail2ban

**CI alternative:** push the inventory change to `main` — `configure.yml` workflow runs automatically.

### 4. Update DNS

Update the A record for `<app-domain>` to point at the app VM's public IP.

### 5. Update GitHub deployment secrets

The existing CD pipeline uses `SSH_HOST` to deploy. Update it to the new app VM IP (Settings → Secrets → Actions → SSH_HOST).

### 6. Bootstrap Let's Encrypt (one-time, app VM only)

```bash
ssh azureuser@<app-vm-ip>
cd ~/app
./init-letsencrypt.sh
```

This generates real TLS certs. Required only on a fresh VM.

### 7. Trigger application deployment

Either push a commit to `main` (triggers existing CD pipeline) or manually trigger the deployment workflow.

### 8. Verify

```bash
curl -I https://<app-domain>           # should return 200
ssh azureuser@<app-vm-ip> 'docker compose ps'
ssh azureuser@<monitoring-vm-ip> 'docker --version && systemctl is-active fail2ban'
```

## Operational workflows

### Configure (`configure.yml`)
Runs Ansible playbook against all VMs. Triggers:
- Manual via GitHub Actions UI
- Push to `main` touching `ansible/**`

### Drift Check (`drift-check.yml`)
Detects configuration drift via Ansible dry-run. Triggers:
- Schedule: 3 AM UTC daily
- Manual via GitHub Actions UI

If drift is detected, opens a GitHub issue tagged `drift, infra` with the diff and remediation links.

### Reconciling drift
1. Read the auto-generated issue
2. Decide: is the change intentional?
   - **Yes** → update `ansible/playbook.yml` to match new desired state, push
   - **No** → manually trigger `configure.yml` to reconcile VMs back to declared state
3. Close the issue

## Teardown

```bash
cd terraform/
terraform destroy
```

Destroys all infrastructure. Application data in Postgres volumes is lost (data is the application's responsibility, not infrastructure). Resource group `knowledge-resource-group` is removed entirely.

## Trade-offs and known limitations

### Terraform runs locally, not in CI
Student Azure subscriptions in our school's tenant don't permit Azure AD application registration, which is required for unattended Terraform CI authentication via service principal. We accepted this constraint and kept Terraform as a manual local operation. Infrastructure changes are deliberate enough that human-driven `terraform apply` is acceptable. Configuration drift is continuous enough that automation pays off — hence Ansible runs in CI.

### Inventory is manually committed
After `terraform apply` issues new IPs, `inventory.ini` must be regenerated and committed for CI to use them. Could be automated if Terraform ran in CI (see above).

### Single SSH key shared across pipelines
Both your existing CD pipeline and the new Ansible workflows reuse `SSH_PRIVATE_KEY`. Convenient, but couples key rotation across all workflows. Real production setups would have separate identities per workflow.

### State stored locally
Terraform state lives on the developer's laptop, not in remote storage. Single-developer scale tolerates this; multi-developer would require migration to Azure Storage backend.

## Related Files

- `terraform/` — infrastructure definitions
- `ansible/` — configuration management
- `.github/workflows/configure.yml` — CI Ansible runner
- `.github/workflows/drift-check.yml` — nightly drift detection
- `.github/ISSUE_TEMPLATE/drift.md` — drift issue template