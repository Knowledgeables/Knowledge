# Infrastructure & Deployment Runbook

This document describes how to provision, configure, and deploy the Knowledge application from scratch using Infrastructure as Code (Terraform) and Configuration Management (Ansible).

## Architecture Overview

The system uses a three-layer separation of concerns:

| Layer | Tool | Purpose | Frequency |
|---|---|---|---|
| Infrastructure | Terraform | VMs, networking, IPs | Rare (manual) |
| Configuration | Ansible | OS-level setup (Docker, fail2ban) | Occasional (CI on `ansible/**` changes) |
| Application | Docker Compose + existing CD pipeline | App deployment | Continuous (every push) |

Two VMs are provisioned in Azure: knowledge (application server, runs the Go app + Postgres + nginx with TLS) and kmonitor (monitoring server, runs Grafana + Loki behind nginx). They share networking (one VNet, one subnet, one NSG) but are configured for distinct roles via the application layer.

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
# macOS / Linux / WSL
terraform output -raw ansible_inventory > ../ansible/inventory.ini
```
```bash
# Windows PowerShell — do NOT use `>` redirection (writes UTF-16 with BOM, breaks Ansible parsing)
terraform output -raw ansible_inventory | Out-File -Encoding ascii ../ansible/inventory.ini
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

### 6. Trigger application deployment

Push a commit to main (triggers existing CD pipeline) or manually trigger the deployment workflow. This transfers docker-compose.yml, nginx configs, and .env to ~/app on the knowledge VM, then runs docker compose pull/up.
The first deploy will fail HTTPS health checks because Let's Encrypt certs don't exist yet. That's expected — proceed to step 7.

### 7. Bootstrap Let's Encrypt (one-time, knowledge VM only)

```bash
ssh azureuser@<knowledge-vm-ip>
cd ~/app
./init-letsencrypt.sh
```

This generates real TLS certs. After it completes, restart the stack to pick up the new certs:

```bash
docker compose down && docker compose up -d
```
Required only on a fresh VM. Subsequent deploys reuse the existing certs.


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

### Verifying drift detection works
The drift-check workflow itself can drift — a typo in a regex, a placeholder left in a command, a mishandled exit code — and silently stop catching real problems. A green workflow run only proves the workflow didn't crash, not that it actually checks anything. Periodically validate the full loop end-to-end:

Introduce deliberate drift on the monitoring VM:
```bash
ssh azureuser@<kmonitor-vm-ip>
sudo systemctl stop fail2ban
```
Manually trigger drift-check.yml from the GitHub Actions UI.

Confirm an issue is auto-opened in the repo, tagged drift, infra, with the diff.

Trigger configure.yml to reconcile, or run the playbook locally.

SSH back to the monitoring VM and confirm systemctl status fail2ban reports active.


Trigger drift-check.yml once more and confirm no new issue is opened.

Close the original issue with a reference to the reconciling run.

Run this validation after any change to drift-check.yml, playbook.yml, or related task files.
## Teardown

```bash
cd terraform/
terraform destroy
```

Destroys everything in knowledge-resource-group: VMs, disks, public IPs, networking, NSG. Application data in Postgres volumes is lost (data is the application's responsibility, not infrastructure).

Shadow infrastructure (one-time cleanup)
Resources that exist outside Terraform's state — manually-created VMs from before IaC was introduced, leftover resource groups from earlier iterations — are not touched by terraform destroy and must be deleted separately:

```bash
az group list --query "[?starts_with(name, 'KNOWLEDGE') || starts_with(name, 'MONITOR')].name" -o tsv
# for each result returned:
az group delete --name <name> --yes
```

## Trade-offs and known limitations

### Terraform runs locally, not in CI
Student Azure subscriptions in our school's tenant don't permit Azure AD application registration, which is required for unattended Terraform CI authentication via service principal. We accepted this constraint and kept Terraform as a manual local operation. Infrastructure changes are deliberate enough that human-driven `terraform apply` is acceptable. Configuration drift is continuous enough that automation pays off — hence Ansible runs in CI.

### Inventory is manually committed
After `terraform apply` issues new IPs, `inventory.ini` must be regenerated and committed for CI to use them. Could be automated if Terraform ran in CI (see above).

### Single SSH key shared across pipelines
Both your existing CD pipeline and the new Ansible workflows reuse `SSH_PRIVATE_KEY`. Convenient, but couples key rotation across all workflows. Real production setups would have separate identities per workflow.

### State stored locally
Terraform state lives on the developer's laptop, not in remote storage. Single-developer scale tolerates this; multi-developer would require migration to Azure Storage backend.

### Shadow infrastructure from pre-IaC era

The original Knowledge and KMonitor VMs were created manually via the Azure Portal before Terraform was introduced, in separate resource groups (RG.KNOWLEDGE, MONITORKNOWLEDGE). Bringing them retroactively into Terraform state would have required reverse-engineering their exact configuration via terraform import, which is fragile when originals were created with portal defaults that don't match declared values.
The chosen path was: keep both running until the project is finished, then tear everything down (Terraform-managed and shadow alike) and rebuild fresh from declarations. This validates the rebuild path end-to-end and eliminates the shadow infrastructure as a side effect. Lesson recorded for future projects: start with IaC from day one, or accept that pre-IaC resources will eventually need a destroy-and-rebuild migration rather than a clean retrofit.

## Related Files

- `terraform/` — infrastructure definitions
- `ansible/` — configuration management
- `.github/workflows/configure.yml` — CI Ansible runner
- `.github/workflows/drift-check.yml` — nightly drift detection
- `.github/ISSUE_TEMPLATE/drift.md` — drift issue template