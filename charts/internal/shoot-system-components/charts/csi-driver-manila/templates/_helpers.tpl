{{- define "csi-driver-node.extensionsGroup" -}}
extensions.gardener.cloud
{{- end -}}

{{- define "csi-driver-node.name" -}}
provider-openstack
{{- end -}}

{{- define "csi-driver-manila-node.provisioner" -}}
manila.csi.openstack.org
{{- end -}}