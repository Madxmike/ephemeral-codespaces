{{- define "instance.name" -}}
    {{- printf "%s-%s" .Release.Name .Values.instance -}}
{{- end -}}


{{- define "mount.path" -}}
    {{- printf "%s/code-server" .Values.container.mount.basePath | clean -}}
{{- end -}}


{{- define "args" }}
    {{- $args := "code-server --extensions-dir ~/extensions" }}
    {{- if .Values.extensions }}
        {{- range .Values.extensions }}
            {{- $args = (printf "%s code-server --install-extension %s" $args .)}}
        {{- end }}
    {{- end}}
    {{- $args = (printf "%s; /init" $args)}}
    {{- $args -}}
{{- end }}