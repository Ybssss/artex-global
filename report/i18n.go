package report

import "github.com/Autumn-27/artex/config"

// reportCatalog translates this package's user-facing document strings.
//
// The report is a deliverable the operator hands to someone else, so it must be
// written in the language they picked in the UI rather than the Korean the
// templates were originally authored in. Keys are the Korean source text and
// format verbs keep their positional order, so the same fmt.Fprintf arguments
// still line up after substitution.
//
// Korean (ko) is the key itself and is deliberately absent: a lookup miss, or an
// active language of ko, returns the source text unchanged. That keeps a Korean
// deployment byte-identical to before this file existed.
//
// config.ActiveLanguage reads the runtime setting the user picked (settings
// table) layered over the bootstrap default, so a language change applies to the
// next generated report without a restart.
var reportCatalog = map[string]map[string]string{
	"# 침투 테스트 보고서: %s\n\n": {
		"en": "# Penetration Test Report: %s\n\n",
		"zh": "# 渗透测试报告：%s\n\n",
		"es": "# Informe de prueba de penetración: %s\n\n",
		"ja": "# ペネトレーションテストレポート: %s\n\n",
	},
	"이름 없는 작업": {
		"en": "Untitled task",
		"zh": "未命名任务",
		"es": "Tarea sin título",
		"ja": "無題のタスク",
	},
	"- **작업 목표**: %s\n": {
		"en": "- **Task goal**: %s\n",
		"zh": "- **任务目标**：%s\n",
		"es": "- **Objetivo de la tarea**: %s\n",
		"ja": "- **タスク目標**: %s\n",
	},
	"(지정되지 않음)": {
		"en": "(not specified)",
		"zh": "（未指定）",
		"es": "(no especificado)",
		"ja": "（未指定）",
	},
	"- **생성 시각**: %s\n\n": {
		"en": "- **Generated at**: %s\n\n",
		"zh": "- **生成时间**：%s\n\n",
		"es": "- **Generado el**: %s\n\n",
		"ja": "- **生成日時**: %s\n\n",
	},
	"## 요약\n\n": {
		"en": "## Summary\n\n",
		"zh": "## 摘要\n\n",
		"es": "## Resumen\n\n",
		"ja": "## 概要\n\n",
	},
	"- 확인된 취약점: **%d**건\n": {
		"en": "- Confirmed vulnerabilities: **%d**\n",
		"zh": "- 已确认漏洞：**%d** 个\n",
		"es": "- Vulnerabilidades confirmadas: **%d**\n",
		"ja": "- 確認済み脆弱性: **%d** 件\n",
	},
	"- 자산: ": {
		"en": "- Assets: ",
		"zh": "- 资产：",
		"es": "- Activos: ",
		"ja": "- 資産: ",
	},
	"## 취약점\n\n": {
		"en": "## Vulnerabilities\n\n",
		"zh": "## 漏洞\n\n",
		"es": "## Vulnerabilidades\n\n",
		"ja": "## 脆弱性\n\n",
	},
	"_이번에 확인된 취약점이 없습니다._\n\n": {
		"en": "_No vulnerabilities were confirmed this time._\n\n",
		"zh": "_本次未确认任何漏洞。_\n\n",
		"es": "_No se confirmaron vulnerabilidades esta vez._\n\n",
		"ja": "_今回は確認された脆弱性はありません。_\n\n",
	},
	"### %d. [%s] %s\n\n": {
		"en": "### %d. [%s] %s\n\n",
		"zh": "### %d. [%s] %s\n\n",
		"es": "### %d. [%s] %s\n\n",
		"ja": "### %d. [%s] %s\n\n",
	},
	"미분류": {
		"en": "Unclassified",
		"zh": "未分类",
		"es": "Sin clasificar",
		"ja": "未分類",
	},
	"**PoC / 증거:**\n\n```\n%s\n```\n\n": {
		"en": "**PoC / Evidence:**\n\n```\n%s\n```\n\n",
		"zh": "**PoC / 证据：**\n\n```\n%s\n```\n\n",
		"es": "**PoC / Evidencia:**\n\n```\n%s\n```\n\n",
		"ja": "**PoC / 証拠:**\n\n```\n%s\n```\n\n",
	},

	// --- findings export (report/findings.go) ---

	"# 취약점 요약 보고서\n\n": {
		"en": "# Vulnerability Summary Report\n\n",
		"zh": "# 漏洞摘要报告\n\n",
		"es": "# Informe resumen de vulnerabilidades\n\n",
		"ja": "# 脆弱性サマリレポート\n\n",
	},
	"- **생성 시각**: %s\n": {
		"en": "- **Generated at**: %s\n",
		"zh": "- **生成时间**：%s\n",
		"es": "- **Generado el**: %s\n",
		"ja": "- **生成日時**: %s\n",
	},
	"- **취약점 총계**: %d건\n\n": {
		"en": "- **Total vulnerabilities**: %d\n\n",
		"zh": "- **漏洞总数**：%d 个\n\n",
		"es": "- **Total de vulnerabilidades**: %d\n\n",
		"ja": "- **脆弱性合計**: %d 件\n\n",
	},
	"| 심각도 | 개수 |\n| --- | --- |\n": {
		"en": "| Severity | Count |\n| --- | --- |\n",
		"zh": "| 严重程度 | 数量 |\n| --- | --- |\n",
		"es": "| Severidad | Cantidad |\n| --- | --- |\n",
		"ja": "| 重大度 | 件数 |\n| --- | --- |\n",
	},
	"심각": {
		"en": "Critical",
		"zh": "严重",
		"es": "Crítico",
		"ja": "深刻",
	},
	"높음": {
		"en": "High",
		"zh": "高",
		"es": "Alto",
		"ja": "高",
	},
	"중간": {
		"en": "Medium",
		"zh": "中",
		"es": "Medio",
		"ja": "中",
	},
	"낮음": {
		"en": "Low",
		"zh": "低",
		"es": "Bajo",
		"ja": "低",
	},
	"_일치하는 취약점이 없습니다._\n": {
		"en": "_No matching vulnerabilities._\n",
		"zh": "_没有匹配的漏洞。_\n",
		"es": "_No hay vulnerabilidades coincidentes._\n",
		"ja": "_一致する脆弱性はありません。_\n",
	},
	"## 취약점 상세\n\n": {
		"en": "## Vulnerability Details\n\n",
		"zh": "## 漏洞详情\n\n",
		"es": "## Detalles de vulnerabilidades\n\n",
		"ja": "## 脆弱性の詳細\n\n",
	},
	"- **유형**: %s\n": {
		"en": "- **Type**: %s\n",
		"zh": "- **类型**：%s\n",
		"es": "- **Tipo**: %s\n",
		"ja": "- **種類**: %s\n",
	},
	"- **상태**: %s\n": {
		"en": "- **Status**: %s\n",
		"zh": "- **状态**：%s\n",
		"es": "- **Estado**: %s\n",
		"ja": "- **状態**: %s\n",
	},
	"- **소속 작업**: %s\n": {
		"en": "- **Owning task**: %s\n",
		"zh": "- **所属任务**：%s\n",
		"es": "- **Tarea propietaria**: %s\n",
		"ja": "- **所属タスク**: %s\n",
	},
	"- **발견 시각**: %s\n\n": {
		"en": "- **Found at**: %s\n\n",
		"zh": "- **发现时间**：%s\n\n",
		"es": "- **Encontrado el**: %s\n\n",
		"ja": "- **発見日時**: %s\n\n",
	},
	"- **발견 시각**: %s\n": {
		"en": "- **Found at**: %s\n",
		"zh": "- **发现时间**：%s\n",
		"es": "- **Encontrado el**: %s\n",
		"ja": "- **発見日時**: %s\n",
	},
	"**증거:**\n\n```\n%s\n```\n\n": {
		"en": "**Evidence:**\n\n```\n%s\n```\n\n",
		"zh": "**证据：**\n\n```\n%s\n```\n\n",
		"es": "**Evidencia:**\n\n```\n%s\n```\n\n",
		"ja": "**証拠:**\n\n```\n%s\n```\n\n",
	},
	"**상세 보고서:**\n\n": {
		"en": "**Detailed report:**\n\n",
		"zh": "**详细报告：**\n\n",
		"es": "**Informe detallado:**\n\n",
		"ja": "**詳細レポート:**\n\n",
	},
	"- **심각도**: %s\n": {
		"en": "- **Severity**: %s\n",
		"zh": "- **严重程度**：%s\n",
		"es": "- **Severidad**: %s\n",
		"ja": "- **重大度**: %s\n",
	},
	"## 개요\n\n%s\n\n": {
		"en": "## Overview\n\n%s\n\n",
		"zh": "## 概述\n\n%s\n\n",
		"es": "## Resumen\n\n%s\n\n",
		"ja": "## 概要\n\n%s\n\n",
	},
	"## 증거\n\n```\n%s\n```\n\n": {
		"en": "## Evidence\n\n```\n%s\n```\n\n",
		"zh": "## 证据\n\n```\n%s\n```\n\n",
		"es": "## Evidencia\n\n```\n%s\n```\n\n",
		"ja": "## 証拠\n\n```\n%s\n```\n\n",
	},
	"## 상세 보고서\n\n": {
		"en": "## Detailed Report\n\n",
		"zh": "## 详细报告\n\n",
		"es": "## Informe detallado\n\n",
		"ja": "## 詳細レポート\n\n",
	},

	// --- CSV export header (report/findings.go) ---

	"ID": {
		"en": "ID",
		"zh": "ID",
		"es": "ID",
		"ja": "ID",
	},
	"이름": {
		"en": "Name",
		"zh": "名称",
		"es": "Nombre",
		"ja": "名前",
	},
	"심각도": {
		"en": "Severity",
		"zh": "严重程度",
		"es": "Severidad",
		"ja": "重大度",
	},
	"유형": {
		"en": "Type",
		"zh": "类型",
		"es": "Tipo",
		"ja": "種類",
	},
	"상태": {
		"en": "Status",
		"zh": "状态",
		"es": "Estado",
		"ja": "状態",
	},
	"소속 작업": {
		"en": "Owning task",
		"zh": "所属任务",
		"es": "Tarea propietaria",
		"ja": "所属タスク",
	},
	"발견 시각": {
		"en": "Found at",
		"zh": "发现时间",
		"es": "Encontrado el",
		"ja": "発見日時",
	},
	"개요": {
		"en": "Overview",
		"zh": "概述",
		"es": "Resumen",
		"ja": "概要",
	},
	"트래픽 증거 개수": {
		"en": "Traffic evidence count",
		"zh": "流量证据数量",
		"es": "Cantidad de evidencia de tráfico",
		"ja": "トラフィック証拠数",
	},
	"트래픽 증거 ID": {
		"en": "Traffic evidence IDs",
		"zh": "流量证据 ID",
		"es": "IDs de evidencia de tráfico",
		"ja": "トラフィック証拠 ID",
	},

	// --- traffic evidence appendix (report/findings.go) ---

	"\n## 관련 트래픽 증거\n\n": {
		"en": "\n## Related Traffic Evidence\n\n",
		"zh": "\n## 相关流量证据\n\n",
		"es": "\n## Evidencia de tráfico relacionada\n\n",
		"ja": "\n## 関連トラフィック証拠\n\n",
	},
	"증거 버전: %d, 바인딩 개수: %d.\n\n": {
		"en": "Evidence version: %d, binding count: %d.\n\n",
		"zh": "证据版本：%d，绑定数量：%d。\n\n",
		"es": "Versión de evidencia: %d, cantidad de vínculos: %d.\n\n",
		"ja": "証拠バージョン: %d、バインディング数: %d。\n\n",
	},
	"증거가 변경되어 상세 보고서를 갱신해야 합니다.\n\n": {
		"en": "Evidence changed; the detailed report must be refreshed.\n\n",
		"zh": "证据已变更，需要更新详细报告。\n\n",
		"es": "La evidencia cambió; es necesario actualizar el informe detallado.\n\n",
		"ja": "証拠が変更されたため、詳細レポートを更新する必要があります。\n\n",
	},
	"%d. **증거 #%d · %s**: `%s %s`, 상태 코드 %d\n": {
		"en": "%d. **Evidence #%d · %s**: `%s %s`, status code %d\n",
		"zh": "%d. **证据 #%d · %s**：`%s %s`，状态码 %d\n",
		"es": "%d. **Evidencia #%d · %s**: `%s %s`, código de estado %d\n",
		"ja": "%d. **証拠 #%d · %s**: `%s %s`、ステータスコード %d\n",
	},
	"   [요청 메시지](evidence/%d/%d/request.http) · [응답 메시지](evidence/%d/%d/response.http)\n": {
		"en": "   [Request message](evidence/%d/%d/request.http) · [Response message](evidence/%d/%d/response.http)\n",
		"zh": "   [请求报文](evidence/%d/%d/request.http) · [响应报文](evidence/%d/%d/response.http)\n",
		"es": "   [Mensaje de solicitud](evidence/%d/%d/request.http) · [Mensaje de respuesta](evidence/%d/%d/response.http)\n",
		"ja": "   [リクエストメッセージ](evidence/%d/%d/request.http) · [レスポンスメッセージ](evidence/%d/%d/response.http)\n",
	},
}

// tr renders one report string in the active output language.
//
// It mirrors config.T but reads ActiveLanguage rather than the bootstrap
// Language: a report is generated on demand, after the user has had a chance to
// change the setting, so the runtime value is the right one.
func tr(ko string) string {
	if m, ok := reportCatalog[ko]; ok {
		if v, ok := m[config.ActiveLanguage()]; ok && v != "" {
			return v
		}
	}
	return ko
}
