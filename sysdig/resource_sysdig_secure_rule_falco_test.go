package sysdig_test

import (
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"os"
	"testing"
)

func TestAccRuleFalco(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SYSDIG_SECURE_API_TOKEN"); v == "" {
				t.Fatal("SYSDIG_SECURE_API_TOKEN must be set for acceptance tests")
			}
		},
		Providers: map[string]terraform.ResourceProvider{
			"sysdig": sysdig.Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config: ruleFalcoTerminalShell(rText()),
			},
			{
				Config: ruleFalcoKubeAudit(rText()),
			},
		},
	})
}

func ruleFalcoTerminalShell(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_falco" "terminal_shell" {
  name = "TERRAFORM TEST %s - Terminal Shell"
  description = "TERRAFORM TEST %s"
  tags = ["container", "shell", "mitre_execution"]

  condition = "spawned_process and container and shell_procs and proc.tty != 0 and container_entrypoint"
  output = "A shell was spawned in a container with an attached terminal (user=%%user.name %%container.info shell=%%proc.name parent=%%proc.pname cmdline=%%proc.cmdline terminal=%%proc.tty container_id=%%container.id image=%%container.image.repository)"
  priority = "notice"
  source = "syscall" // syscall or k8s_audit
}`, name, name)
}

func ruleFalcoKubeAudit(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_rule_falco" "kube_audit" {
  name = "TERRAFORM TEST %s - KubeAudit"
  description = "TERRAFORM TEST %s"
  tags = ["k8s"]

  condition = "kall"
  output = "K8s Audit Event received (user=%%ka.user.name verb=%%ka.verb uri=%%ka.uri obj=%%jevt.obj)"
  priority = "debug"
  source = "k8s_audit" // syscall or k8s_audit
}`, name, name)
}
