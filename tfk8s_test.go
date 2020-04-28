package main

import (
	"strings"
	"testing"
)

// TODO add diffing for when the tests fail

func TestToHCLSingle(t *testing.T) {
	yaml := `---
apiVersion: v1
kind: ConfigMap
metadata:
  name: test
data:
  TEST: test`

	r := strings.NewReader(yaml)
	output, err := ToHCL(r)

	if err != nil {
		t.Fatal("Converting to HCL failed:", err)
	}

	expected := `
resource "kubernetes_manifest_hcl" "configmap_test" {
  manifest = {
    "apiVersion" = "v1"
    "data" = {
      "TEST" = "test"
    }
    "kind" = "ConfigMap"
    "metadata" = {
      "name" = "test"
    }
  }
}`

	if strings.TrimSpace(expected) != strings.TrimSpace(output) {
		t.Errorf("\nExpected:\n%v\nGot:\n%v\n", strings.TrimSpace(expected), strings.TrimSpace(output))
	}
}

func TestToHCLMultiple(t *testing.T) {
	yaml := `---
apiVersion: v1
kind: ConfigMap
metadata:
  name: one
data:
  TEST: one
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: two
data:
  TEST: two`

	r := strings.NewReader(yaml)
	output, err := ToHCL(r)

	if err != nil {
		t.Fatal("Converting to HCL failed:", err)
	}

	expected := `
resource "kubernetes_manifest_hcl" "configmap_one" {
  manifest = {
    "apiVersion" = "v1"
    "data" = {
      "TEST" = "one"
    }
    "kind" = "ConfigMap"
    "metadata" = {
      "name" = "one"
    }
  }
}

resource "kubernetes_manifest_hcl" "configmap_two" {
  manifest = {
    "apiVersion" = "v1"
    "data" = {
      "TEST" = "two"
    }
    "kind" = "ConfigMap"
    "metadata" = {
      "name" = "two"
    }
  }
}`

	if strings.TrimSpace(expected) != strings.TrimSpace(output) {
		t.Errorf("\nExpected:\n%v\n\nGot:\n%v\n", strings.TrimSpace(expected), strings.TrimSpace(output))
	}
}