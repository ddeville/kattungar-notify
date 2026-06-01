{{/*
Expand the name of the chart.
*/}}
{{- define "kattungar-notify.name" -}}
{{- .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "kattungar-notify.fullname" -}}
{{- $name := .Chart.Name -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "kattungar-notify.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels.
*/}}
{{- define "kattungar-notify.labels" -}}
helm.sh/chart: {{ include "kattungar-notify.chart" . }}
{{ include "kattungar-notify.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Selector labels.
*/}}
{{- define "kattungar-notify.selectorLabels" -}}
app.kubernetes.io/name: {{ include "kattungar-notify.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
Secret name used by the deployment.
*/}}
{{- define "kattungar-notify.secretName" -}}
{{- default (include "kattungar-notify.fullname" .) .Values.secret.name -}}
{{- end -}}

{{/*
Persistent volume claim name used by the deployment.
*/}}
{{- define "kattungar-notify.pvcName" -}}
{{- default (include "kattungar-notify.fullname" .) .Values.persistence.existingClaim -}}
{{- end -}}

{{/*
Container image reference.
*/}}
{{- define "kattungar-notify.image" -}}
{{- $tag := .Values.image.tag | default .Chart.Version -}}
{{- if .Values.image.digest -}}
{{- printf "%s:%s@%s" .Values.image.repository $tag .Values.image.digest -}}
{{- else -}}
{{- printf "%s:%s" .Values.image.repository $tag -}}
{{- end -}}
{{- end -}}
