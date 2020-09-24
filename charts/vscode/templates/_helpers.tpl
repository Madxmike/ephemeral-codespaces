{{- define "instance.name" -}}
    {{- printf "%s-%s" .Release.Name .Values.instance -}}
{{- end -}}


{{- define "mount.path" -}}
    {{- printf "%s/code-server" .Values.container.mount.basePath | clean -}}
{{- end -}}


{{- define "args.extensions" }}
    {{- if .Values.extensions }}
        {{- range .Values.extensions }}
            {{- printf " code-server --install-extension %s" . }};
        {{- end }}
    {{- end}}
{{- end }}

{{- define "args.runtimes" }}
    {{- if .Values.runtimes }}
        {{- $files := .Files }}
        {{- range .Values.runtimes }}
            {{- range $files.Lines (printf "%s-%s.sh" .name .version) }}
                {{- printf " %s;" .}}
            {{- end}}
        {{- end }}
    {{- end}}
{{- end }}
