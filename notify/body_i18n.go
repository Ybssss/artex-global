package notify

// Message-body translations.
//
// The notify package already localises severity and status LABELS through
// labelLang() (see notify.go). The surrounding prose — headings, the
// "view on platform" footer, the truncation note — was left in Korean, so a
// channel configured for English or Chinese still delivered Korean sentences.
//
// This catalog closes that gap. Keys are the Korean source templates and format
// verbs keep their positional order, so the same fmt.Fprintf arguments line up
// after substitution.
//
// Language resolution deliberately matches the labels: labelLang() reads the
// server-wired Lang hook, falling back to "ko" only when nothing is wired (unit
// tests / standalone). Korean is the key itself and is absent from the maps, so
// a lookup miss or an active language of ko returns the source unchanged —
// keeping a Korean deployment byte-identical to before this file existed.
var bodyCatalog = map[string]map[string]string{
	// --- digest headings (markdown.go) ---
	"취약점 요약 · 총 %d건": {
		"en": "Vulnerability summary · %d total",
		"zh": "漏洞摘要 · 共 %d 条",
		"es": "Resumen de vulnerabilidades · %d en total",
		"ja": "脆弱性サマリ · 合計 %d 件",
	},
	"취약점 알림": {
		"en": "Vulnerability alert",
		"zh": "漏洞告警",
		"es": "Alerta de vulnerabilidad",
		"ja": "脆弱性アラート",
	},
	"\n[플랫폼에서 전체 보기](%s)\n": {
		"en": "\n[View all on the platform](%s)\n",
		"zh": "\n[在平台查看全部](%s)\n",
		"es": "\n[Ver todo en la plataforma](%s)\n",
		"ja": "\n[プラットフォームで全体を見る](%s)\n",
	},
	"**최근 %d분간 신규 취약점 %d건**": {
		"en": "**%d new vulnerabilities in the last %d minutes**",
		"zh": "**最近 %d 分钟新增漏洞 %d 个**",
		"es": "**%d nuevas vulnerabilidades en los últimos %d minutos**",
		"ja": "**直近 %d 分間の新規脆弱性 %d 件**",
	},
	"**신규 취약점 %d건**": {
		"en": "**%d new vulnerabilities**",
		"zh": "**新增漏洞 %d 个**",
		"es": "**%d nuevas vulnerabilidades**",
		"ja": "**新規脆弱性 %d 件**",
	},
	"(이 메시지에는 앞 %d건만 표시하며, 나머지 %d건은 다음 메시지에서 이어서 보냅니다)": {
		"en": "(only the first %d are shown here; the remaining %d continue in the next message)",
		"zh": "（本条仅显示前 %d 条，其余 %d 条将在下一条消息中继续发送）",
		"es": "(aquí solo se muestran los primeros %d; los %d restantes continúan en el siguiente mensaje)",
		"ja": "（このメッセージには先頭 %d 件のみを表示し、残り %d 件は次のメッセージで続けて送信します）",
	},

	// --- per-item fields (markdown.go) ---
	"**상태 변경**: %s → %s\n": {
		"en": "**Status change**: %s → %s\n",
		"zh": "**状态变更**：%s → %s\n",
		"es": "**Cambio de estado**: %s → %s\n",
		"ja": "**ステータス変更**: %s → %s\n",
	},
	"**유형**: %s\n": {
		"en": "**Type**: %s\n",
		"zh": "**类型**：%s\n",
		"es": "**Tipo**: %s\n",
		"ja": "**種類**: %s\n",
	},
	"**자산**: %s\n": {
		"en": "**Assets**: %s\n",
		"zh": "**资产**：%s\n",
		"es": "**Activos**: %s\n",
		"ja": "**資産**: %s\n",
	},
	"**개요**: %s\n": {
		"en": "**Overview**: %s\n",
		"zh": "**概述**：%s\n",
		"es": "**Resumen**: %s\n",
		"ja": "**概要**: %s\n",
	},
	"[상세 보기](%s)\n": {
		"en": "[View details](%s)\n",
		"zh": "[查看详情](%s)\n",
		"es": "[Ver detalles](%s)\n",
		"ja": "[詳細を見る](%s)\n",
	},

	// --- HTML channel (html.go) ---
	`<p style="margin:16px 0 0;"><a href="%s" style="color:#1677ff;">플랫폼에서 전체 보기</a></p>`: {
		"en": `<p style="margin:16px 0 0;"><a href="%s" style="color:#1677ff;">View all on the platform</a></p>`,
		"zh": `<p style="margin:16px 0 0;"><a href="%s" style="color:#1677ff;">在平台查看全部</a></p>`,
		"es": `<p style="margin:16px 0 0;"><a href="%s" style="color:#1677ff;">Ver todo en la plataforma</a></p>`,
		"ja": `<p style="margin:16px 0 0;"><a href="%s" style="color:#1677ff;">プラットフォームで全体を見る</a></p>`,
	},
	`<h2 style="font-size:16px;margin:0 0 4px;">최근 %d분간 신규 취약점 %d건</h2>`: {
		"en": `<h2 style="font-size:16px;margin:0 0 4px;">%d new vulnerabilities in the last %d minutes</h2>`,
		"zh": `<h2 style="font-size:16px;margin:0 0 4px;">最近 %d 分钟新增漏洞 %d 个</h2>`,
		"es": `<h2 style="font-size:16px;margin:0 0 4px;">%d nuevas vulnerabilidades en los últimos %d minutos</h2>`,
		"ja": `<h2 style="font-size:16px;margin:0 0 4px;">直近 %d 分間の新規脆弱性 %d 件</h2>`,
	},
	`<h2 style="font-size:16px;margin:0 0 4px;">신규 취약점 %d건</h2>`: {
		"en": `<h2 style="font-size:16px;margin:0 0 4px;">%d new vulnerabilities</h2>`,
		"zh": `<h2 style="font-size:16px;margin:0 0 4px;">新增漏洞 %d 个</h2>`,
		"es": `<h2 style="font-size:16px;margin:0 0 4px;">%d nuevas vulnerabilidades</h2>`,
		"ja": `<h2 style="font-size:16px;margin:0 0 4px;">新規脆弱性 %d 件</h2>`,
	},
	`<div><b>상태 변경</b>: %s → %s</div>`: {
		"en": `<div><b>Status change</b>: %s → %s</div>`,
		"zh": `<div><b>状态变更</b>：%s → %s</div>`,
		"es": `<div><b>Cambio de estado</b>: %s → %s</div>`,
		"ja": `<div><b>ステータス変更</b>: %s → %s</div>`,
	},
	`<div><b>유형</b>: %s</div>`: {
		"en": `<div><b>Type</b>: %s</div>`,
		"zh": `<div><b>类型</b>：%s</div>`,
		"es": `<div><b>Tipo</b>: %s</div>`,
		"ja": `<div><b>種類</b>: %s</div>`,
	},
	`<div><b>자산</b>: %s</div>`: {
		"en": `<div><b>Assets</b>: %s</div>`,
		"zh": `<div><b>资产</b>：%s</div>`,
		"es": `<div><b>Activos</b>: %s</div>`,
		"ja": `<div><b>資産</b>: %s</div>`,
	},
	`<div><b>개요</b>: %s</div>`: {
		"en": `<div><b>Overview</b>: %s</div>`,
		"zh": `<div><b>概述</b>：%s</div>`,
		"es": `<div><b>Resumen</b>: %s</div>`,
		"ja": `<div><b>概要</b>: %s</div>`,
	},
	`<div style="margin-top:6px;"><a href="%s" style="color:#1677ff;">상세 보기</a></div>`: {
		"en": `<div style="margin-top:6px;"><a href="%s" style="color:#1677ff;">View details</a></div>`,
		"zh": `<div style="margin-top:6px;"><a href="%s" style="color:#1677ff;">查看详情</a></div>`,
		"es": `<div style="margin-top:6px;"><a href="%s" style="color:#1677ff;">Ver detalles</a></div>`,
		"ja": `<div style="margin-top:6px;"><a href="%s" style="color:#1677ff;">詳細を見る</a></div>`,
	},
}

// bodyText renders one message-body string in the active label language.
func bodyText(ko string) string {
	if m, ok := bodyCatalog[ko]; ok {
		if v, ok := m[labelLang()]; ok && v != "" {
			return v
		}
	}
	return ko
}
