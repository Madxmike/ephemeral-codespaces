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

    {{- $args = (printf "%s;" $args)}}

    {{- if .Values.runtimes }}
        {{- $files := .Files }}
        {{- range .Values.runtimes }}
            {{- $runtimeScript := $files.Get (printf "/runtimes/%s-%s.sh" .name .version) }}
            {{- $args = (printf " %s %s" $args $runtimeScript) }}
        {{- end }}
    {{- end }}
    {{- $args = (printf "%s; /init" $args)}}
    {{- $args -}}
{{- end }}