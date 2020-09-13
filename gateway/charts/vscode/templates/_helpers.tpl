{{- define "instance.name" -}}
    {{- printf "%s-%s" .Release.Name .Values.instance -}}
{{- end -}}


{{- define "mount.path" -}}
    {{- printf "%s/mount" .Values.container.mount.basePath | clean -}}
{{- end -}}
