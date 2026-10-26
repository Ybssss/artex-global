package server

// 由 err.Error() 直通的用户可见错误消息翻译（补充 errCatalog）。
//
// errmsg_catalog.go 的 errCatalog 只覆盖 writeErr 能看到的模板；而下面这些字符串是
// 经 errors.New / fmt.Errorf 构造、再以 err.Error() 原样塞进 JSON 响应或写进数据库的
// （归档任务状态、通知投递失败原因、侧问错误、拦截校验等）。writeErr 的中央 trMsg
// 看不到它们，所以此前无论选什么语言都显示韩文。
//
// 单独成文件而不是并进 errmsg_catalog.go：那个 map 很长，脚本拼接容易落错位置；
// 用 init() 合并两份 map 更安全，也便于 review。
//
// 键是韩文原文（含格式化占位符，顺序与调用点一致），ko 不存，查不到就原样返回。
var errCatalogExtra = map[string]map[string]string{
	"%s 필드의 내용에 마스킹 표시(%q)가 들어 있습니다. 이 필드는 전체를 비워 두어 기존 값을 유지하거나 전체를 새 값으로 제출해야 하며, 구조체 내부에 마스킹 표시를 끼워 넣을 수 없습니다": {
		"en": "the %s field contains a masking marker (%q). Clear the whole field to keep the existing value, or submit the whole field as a new value; a masking marker cannot be embedded inside a structure",
		"zh": "%s 字段的内容中含有掩码标记（%q）。该字段须整体留空以保留原值，或整体提交为新值；不能在结构体内部插入掩码标记",
		"es": "el campo %s contiene una marca de enmascarado (%q). Deje el campo entero vacío para conservar el valor actual, o envíe el campo entero como valor nuevo; no se puede insertar una marca de enmascarado dentro de una estructura", "ja": "%s フィールドの内容にマスク表示（%q）が含まれています。このフィールドは全体を空にして既存値を維持するか、全体を新しい値として送信する必要があり、構造体の内部にマスク表示を挿入することはできません",
	},
	"action 값은 allow, deny, ask 중 하나여야 합니다": {
		"en": "action must be one of allow, deny, ask",
		"zh": "action 必须是 allow、deny、ask 之一",
		"es": "action debe ser allow, deny o ask", "ja": "action は allow、deny、ask のいずれかである必要があります",
	},
	"API 주소가 올바르지 않습니다: %w": {
		"en": "the API address is not valid: %w",
		"zh": "API 地址无效：%w",
		"es": "la dirección de la API no es válida: %w", "ja": "API アドレスが正しくありません: %w",
	},
	"API 주소를 조합하지 못했습니다 (API 주소: %s)": {
		"en": "could not assemble the API address (API address: %s)",
		"zh": "无法拼接 API 地址（API 地址：%s）",
		"es": "no se pudo componer la dirección de la API (dirección de la API: %s)", "ja": "API アドレスを組み立てられませんでした（API アドレス: %s）",
	},
	"archive_ids 값은 비워 둘 수 없습니다": {
		"en": "archive_ids cannot be empty",
		"zh": "archive_ids 不能为空",
		"es": "archive_ids no puede estar vacío", "ja": "archive_ids は空にできません",
	},
	"Bot Token이 없습니다": {
		"en": "Bot Token is missing",
		"zh": "Bot Token 缺失",
		"es": "Falta el Bot Token", "ja": "Bot Token がありません",
	},
	"Chat ID가 없습니다": {
		"en": "Chat ID is missing",
		"zh": "Chat ID 缺失",
		"es": "Falta el Chat ID", "ja": "Chat ID がありません",
	},
	"decision_source 값은 model, rule, unknown 중 하나여야 합니다": {
		"en": "decision_source must be one of model, rule, unknown",
		"zh": "decision_source 必须是 model、rule、unknown 之一",
		"es": "decision_source debe ser model, rule o unknown", "ja": "decision_source は model、rule、unknown のいずれかである必要があります",
	},
	"DingTalk 응답을 해석하지 못했습니다: %w (%s)": {
		"en": "could not parse the DingTalk response: %w (%s)",
		"zh": "无法解析钉钉响应：%w（%s）",
		"es": "no se pudo analizar la respuesta de DingTalk: %w (%s)", "ja": "DingTalk の応答を解析できませんでした: %w (%s)",
	},
	"DingTalk에서 오류가 발생했습니다 (%d): %s": {
		"en": "DingTalk returned an error (%d): %s",
		"zh": "钉钉返回错误（%d）：%s",
		"es": "DingTalk devolvió un error (%d): %s", "ja": "DingTalk でエラーが発生しました（%d）: %s",
	},
	"Feishu 응답을 해석하지 못했습니다: %w (%s)": {
		"en": "could not parse the Feishu response: %w (%s)",
		"zh": "无法解析飞书响应：%w（%s）",
		"es": "no se pudo analizar la respuesta de Feishu: %w (%s)", "ja": "Feishu の応答を解析できませんでした: %w (%s)",
	},
	"Feishu에서 오류가 발생했습니다 (%d): %s": {
		"en": "Feishu returned an error (%d): %s",
		"zh": "飞书返回错误（%d）：%s",
		"es": "Feishu devolvió un error (%d): %s", "ja": "Feishu でエラーが発生しました（%d）: %s",
	},
	"http 또는 https 만 지원합니다. 받은 값: %q": {
		"en": "only http or https is supported. Got: %q",
		"zh": "仅支持 http 或 https。收到的值：%q",
		"es": "solo se admite http o https. Recibido: %q", "ja": "http または https のみ対応しています。受け取った値: %q",
	},
	"LLM 폴백 체인: 사용 가능한 설정이 없습니다": {
		"en": "LLM fallback chain: no usable configuration",
		"zh": "LLM 回退链：没有可用的配置",
		"es": "cadena de respaldo de LLM: no hay ninguna configuración utilizable", "ja": "LLM フォールバックチェーン: 使用可能な設定がありません",
	},
	"match_target 값은 tool_name 또는 tool_input 중 하나여야 합니다": {
		"en": "match_target must be either tool_name or tool_input",
		"zh": "match_target 必须是 tool_name 或 tool_input",
		"es": "match_target debe ser tool_name o tool_input", "ja": "match_target は tool_name または tool_input のいずれかである必要があります",
	},
	"match_type 값은 string 또는 regex 중 하나여야 합니다": {
		"en": "match_type must be either string or regex",
		"zh": "match_type 必须是 string 或 regex",
		"es": "match_type debe ser string o regex", "ja": "match_type は string または regex のいずれかである必要があります",
	},
	"name 값은 비워 둘 수 없습니다": {
		"en": "name cannot be empty",
		"zh": "name 不能为空",
		"es": "name no puede estar vacío", "ja": "name は空にできません",
	},
	"pattern 값은 비워 둘 수 없습니다": {
		"en": "pattern cannot be empty",
		"zh": "pattern 不能为空",
		"es": "pattern no puede estar vacío", "ja": "pattern は空にできません",
	},
	"pattern 값이 올바른 정규식이 아닙니다: %w": {
		"en": "pattern is not a valid regular expression: %w",
		"zh": "pattern 不是有效的正则表达式：%w",
		"es": "pattern no es una expresión regular válida: %w", "ja": "pattern が正しい正規表現ではありません: %w",
	},
	"SMTP DATA 실패: %w": {
		"en": "SMTP DATA failed: %w",
		"zh": "SMTP DATA 失败：%w",
		"es": "SMTP DATA falló: %w", "ja": "SMTP DATA に失敗しました: %w",
	},
	"SMTP 서버 연결 실패: %w": {
		"en": "failed to connect to the SMTP server: %w",
		"zh": "连接 SMTP 服务器失败：%w",
		"es": "no se pudo conectar al servidor SMTP: %w", "ja": "SMTP サーバーへの接続に失敗しました: %w",
	},
	"SMTP 서버 주소가 없습니다": {
		"en": "the SMTP server address is missing",
		"zh": "SMTP 服务器地址缺失",
		"es": "Falta la dirección del servidor SMTP", "ja": "SMTP サーバーアドレスがありません",
	},
	"SMTP 인증 실패: %w": {
		"en": "SMTP authentication failed: %w",
		"zh": "SMTP 认证失败：%w",
		"es": "la autenticación SMTP falló: %w", "ja": "SMTP 認証に失敗しました: %w",
	},
	"SMTP 포트가 올바르지 않습니다 (1-65535 범위여야 합니다)": {
		"en": "the SMTP port is not valid (it must be in the 1-65535 range)",
		"zh": "SMTP 端口无效（必须在 1-65535 范围内）",
		"es": "el puerto SMTP no es válido (debe estar en el rango 1-65535)", "ja": "SMTP ポートが正しくありません（1〜65535 の範囲である必要があります）",
	},
	"SMTP 핸드셰이크 실패: %w": {
		"en": "SMTP handshake failed: %w",
		"zh": "SMTP 握手失败：%w",
		"es": "el handshake SMTP falló: %w", "ja": "SMTP ハンドシェイクに失敗しました: %w",
	},
	"STARTTLS 실패: %w": {
		"en": "STARTTLS failed: %w",
		"zh": "STARTTLS 失败：%w",
		"es": "STARTTLS falló: %w", "ja": "STARTTLS に失敗しました: %w",
	},
	"status 값은 pending, allowed, denied, timeout 중 하나여야 합니다": {
		"en": "status must be one of pending, allowed, denied, timeout",
		"zh": "status 必须是 pending、allowed、denied、timeout 之一",
		"es": "status debe ser pending, allowed, denied o timeout", "ja": "status は pending、allowed、denied、timeout のいずれかである必要があります",
	},
	"Telegram 요청 제한에 걸렸습니다: %s": {
		"en": "the Telegram rate limit was hit: %s",
		"zh": "已触发 Telegram 请求限制：%s",
		"es": "se alcanzó el límite de solicitudes de Telegram: %s", "ja": "Telegram のリクエスト制限に達しました: %s",
	},
	"Telegram 응답 오류 %d: %s": {
		"en": "Telegram response error %d: %s",
		"zh": "Telegram 响应错误 %d：%s",
		"es": "error de respuesta de Telegram %d: %s", "ja": "Telegram の応答エラー %d: %s",
	},
	"Telegram 응답을 해석하지 못했습니다: %w (%s)": {
		"en": "could not parse the Telegram response: %w (%s)",
		"zh": "无法解析 Telegram 响应：%w（%s）",
		"es": "no se pudo analizar la respuesta de Telegram: %w (%s)", "ja": "Telegram の応答を解析できませんでした: %w (%s)",
	},
	"Webhook 주소가 없습니다": {
		"en": "the webhook address is missing",
		"zh": "Webhook 地址缺失",
		"es": "Falta la dirección del webhook", "ja": "Webhook アドレスがありません",
	},
	"Webhook 주소가 올바르지 않습니다: %w": {
		"en": "the webhook address is not valid: %w",
		"zh": "Webhook 地址无效：%w",
		"es": "la dirección del webhook no es válida: %w", "ja": "Webhook アドレスが正しくありません: %w",
	},
	"Webhook 주소를 해석하지 못했습니다: %s": {
		"en": "could not resolve the webhook address: %s",
		"zh": "无法解析 Webhook 地址：%s",
		"es": "no se pudo resolver la dirección del webhook: %s", "ja": "Webhook アドレスを解決できませんでした: %s",
	},
	"WeCom 요청이 제한되었습니다 (%d): %s": {
		"en": "the WeCom request was rate limited (%d): %s",
		"zh": "企业微信请求被限流（%d）：%s",
		"es": "la solicitud a WeCom fue limitada (%d): %s", "ja": "WeCom のリクエストが制限されました（%d）: %s",
	},
	"WeCom 응답을 해석하지 못했습니다: %w (%s)": {
		"en": "could not parse the WeCom response: %w (%s)",
		"zh": "无法解析企业微信响应：%w（%s）",
		"es": "no se pudo analizar la respuesta de WeCom: %w (%s)", "ja": "WeCom の応答を解析できませんでした: %w (%s)",
	},
	"WeCom에서 오류가 발생했습니다 (%d): %s": {
		"en": "WeCom returned an error (%d): %s",
		"zh": "企业微信返回错误（%d）：%s",
		"es": "WeCom devolvió un error (%d): %s", "ja": "WeCom でエラーが発生しました（%d）: %s",
	},
	"곁질문 기록 커서가 유효하지 않습니다.": {
		"en": "the side-question history cursor is not valid.",
		"zh": "侧问历史游标无效。",
		"es": "el cursor del historial de preguntas secundarias no es válido.", "ja": "脇質問の履歴カーソルが正しくありません。",
	},
	"곁질문 압축으로 컨텍스트를 더 줄이지 못해 재시도를 중단했습니다.": {
		"en": "compaction could not shrink the side-question context any further, so the retry was stopped.",
		"zh": "侧问压缩已无法进一步缩小上下文，重试已停止。",
		"es": "la compactación no pudo reducir más el contexto de la pregunta secundaria, por lo que se detuvo el reintento.", "ja": "脇質問の圧縮でこれ以上コンテキストを減らせないため、再試行を中止しました。",
	},
	"곁질문 요약이 예산 이내로 줄어들지 않았습니다. 다시 시도해 주세요.": {
		"en": "the side-question summary could not be reduced within budget. Try again",
		"zh": "侧问摘要未能压缩到预算以内。请重试",
		"es": "el resumen de la pregunta secundaria no se pudo reducir dentro del presupuesto. Inténtelo de nuevo", "ja": "脇質問の要約を予算内に収められませんでした。もう一度お試しください",
	},
	"곁질문 요약이 완전히 생성되지 않았습니다. 다시 시도해 주세요.": {
		"en": "the side-question summary was not fully generated. Try again",
		"zh": "侧问摘要未完整生成。请重试",
		"es": "el resumen de la pregunta secundaria no se generó por completo. Inténtelo de nuevo", "ja": "脇質問の要約が完全に生成されませんでした。もう一度お試しください",
	},
	"곁질문 컨텍스트 정리가 이번 처리 상한에 도달했습니다. 질문 범위를 좁힌 뒤 다시 시도해 주세요.": {
		"en": "side-question context compaction hit its per-run cap. Narrow the question and try again",
		"zh": "侧问上下文整理已达到本次处理上限。请缩小问题范围后重试",
		"es": "la compactación del contexto de la pregunta secundaria alcanzó su límite por ejecución. Acote la pregunta e inténtelo de nuevo", "ja": "脇質問のコンテキスト整理が今回の処理上限に達しました。質問の範囲を狭めてから再試行してください",
	},
	"곁질문 컨텍스트를 압축해도 모델 예산을 초과합니다. 질문 범위를 좁히거나 모델 컨텍스트 설정을 조정해 주세요.": {
		"en": "even after compacting, the side-question context exceeds the model budget. Narrow the question or adjust the model context setting",
		"zh": "即使压缩后，侧问上下文仍超出模型预算。请缩小问题范围或调整模型上下文设置",
		"es": "incluso tras compactar, el contexto de la pregunta secundaria supera el presupuesto del modelo. Acote la pregunta o ajuste la configuración de contexto del modelo", "ja": "圧縮しても脇質問のコンテキストがモデル予算を超えます。質問の範囲を狭めるか、モデルのコンテキスト設定を調整してください",
	},
	"다른 호스트로의 리다이렉트는 허용하지 않습니다 (%s → %s)": {
		"en": "redirects to a different host are not allowed (%s → %s)",
		"zh": "不允许重定向到其他主机（%s → %s）",
		"es": "no se permiten redirecciones a otro host (%s → %s)", "ja": "別のホストへのリダイレクトは許可されていません（%s → %s）",
	},
	"대상 URL이 없습니다": {
		"en": "the target URL is missing",
		"zh": "目标 URL 缺失",
		"es": "Falta la URL de destino", "ja": "対象 URL がありません",
	},
	"대상 URL이 올바르지 않습니다: %w": {
		"en": "the target URL is not valid: %w",
		"zh": "目标 URL 无效：%w",
		"es": "la URL de destino no es válida: %w", "ja": "対象 URL が正しくありません: %w",
	},
	"대상 주소를 해석할 수 없습니다: %q": {
		"en": "cannot resolve the target address: %q",
		"zh": "无法解析目标地址：%q",
		"es": "no se puede resolver la dirección de destino: %q", "ja": "対象アドレスを解決できません: %q",
	},
	"등록되지 않은 채널 유형입니다: %q": {
		"en": "unregistered channel type: %q",
		"zh": "未注册的渠道类型：%q",
		"es": "tipo de canal no registrado: %q", "ja": "登録されていないチャネル種別です: %q",
	},
	"로컬 또는 링크 로컬 주소 %s 로는 전송하지 않습니다 (로컬 서비스로 보내야 하면 %s=1 로 설정하세요)": {
		"en": "not sending to the local or link-local address %s (if you must reach a local service, set %s=1)",
		"zh": "不会向本地或链路本地地址 %s 发送（如果确实要发送到本地服务，请设置 %s=1）",
		"es": "no se envía a la dirección local o link-local %s (si necesita llegar a un servicio local, configure %s=1)", "ja": "ローカルまたはリンクローカルアドレス %s には送信しません（ローカルサービスに送る必要がある場合は %s=1 を設定してください）",
	},
	"리다이렉트가 너무 많습니다": {
		"en": "too many redirects",
		"zh": "重定向次数过多",
		"es": "demasiadas redirecciones", "ja": "リダイレクトが多すぎます",
	},
	"메일 본문 쓰기 실패: %w": {
		"en": "failed to write the mail body: %w",
		"zh": "写入邮件正文失败：%w",
		"es": "no se pudo escribir el cuerpo del correo: %w", "ja": "メール本文の書き込みに失敗しました: %w",
	},
	"메일 제출 실패: %w": {
		"en": "failed to submit the mail: %w",
		"zh": "提交邮件失败：%w",
		"es": "no se pudo enviar el correo: %w", "ja": "メールの送信に失敗しました: %w",
	},
	"모델 응답이 중단되었습니다. 다시 질문해 주세요.": {
		"en": "the model response was interrupted. Please ask again",
		"zh": "模型响应中断。请重新提问",
		"es": "la respuesta del modelo se interrumpió. Vuelva a preguntar", "ja": "モデルの応答が中断されました。もう一度質問してください",
	},
	"모델이 답변을 반환하지 않았습니다.": {
		"en": "the model returned no answer",
		"zh": "模型未返回回答",
		"es": "el modelo no devolvió ninguna respuesta", "ja": "モデルが回答を返しませんでした",
	},
	"발신자 주소가 없습니다": {
		"en": "the sender address is missing",
		"zh": "发件人地址缺失",
		"es": "Falta la dirección del remitente", "ja": "送信元アドレスがありません",
	},
	"보관 id %d 가 유효하지 않습니다": {
		"en": "archive id %d is not valid",
		"zh": "归档 id %d 无效",
		"es": "el id de archivo %d no es válido", "ja": "アーカイブ id %d が正しくありません",
	},
	"보관 스냅샷을 만드는 동안 작업 탐색 기록이 변경되었습니다": {
		"en": "the task exploration record changed while the archive snapshot was being taken",
		"zh": "生成归档快照期间，任务的探索记录发生了变化",
		"es": "el registro de exploración de la tarea cambió mientras se tomaba la instantánea del archivo", "ja": "アーカイブスナップショットの作成中にタスクの探索記録が変更されました",
	},
	"보관 패키지 경로가 관리 대상 디렉터리 안에 있지 않습니다": {
		"en": "the archive package path is not inside the managed directory",
		"zh": "归档包路径不在受管目录内",
		"es": "la ruta del paquete de archivo no está dentro del directorio gestionado", "ja": "アーカイブパッケージのパスが管理対象ディレクトリの内側にありません",
	},
	"보관 패키지 경로가 비어 있습니다": {
		"en": "the archive package path is empty",
		"zh": "归档包路径为空",
		"es": "la ruta del paquete de archivo está vacía", "ja": "アーカイブパッケージのパスが空です",
	},
	"상대 서버가 요청을 거부했습니다 (HTTP %d): %s": {
		"en": "the remote server rejected the request (HTTP %d): %s",
		"zh": "对端服务器拒绝了请求（HTTP %d）：%s",
		"es": "el servidor remoto rechazó la solicitud (HTTP %d): %s", "ja": "相手サーバーがリクエストを拒否しました（HTTP %d）: %s",
	},
	"상대 서버가 요청을 제한하거나 응답이 지연되었습니다 (HTTP %d): %s": {
		"en": "the remote server rate limited the request or the response was delayed (HTTP %d): %s",
		"zh": "对端服务器限制了请求或响应延迟（HTTP %d）：%s",
		"es": "el servidor remoto limitó la solicitud o la respuesta se retrasó (HTTP %d): %s", "ja": "相手サーバーがリクエストを制限したか、応答が遅延しました（HTTP %d）: %s",
	},
	"상대 서버에 오류가 발생했습니다 (HTTP %d): %s": {
		"en": "the remote server returned an error (HTTP %d): %s",
		"zh": "对端服务器发生错误（HTTP %d）：%s",
		"es": "el servidor remoto devolvió un error (HTTP %d): %s", "ja": "相手サーバーでエラーが発生しました（HTTP %d）: %s",
	},
	"수신자 주소가 최소 하나 필요합니다": {
		"en": "at least one recipient address is required",
		"zh": "至少需要一个收件人地址",
		"es": "se requiere al menos una dirección de destinatario", "ja": "宛先アドレスが少なくとも 1 つ必要です",
	},
	"올바르지 않은 최소 심각도입니다: %q. low / medium / high / critical 중에서 선택하거나, 비워 두면 제한하지 않습니다": {
		"en": "invalid minimum severity: %q. Choose from low / medium / high / critical, or leave it empty to impose no limit",
		"zh": "最小严重程度无效：%q。请从 low / medium / high / critical 中选择，留空则不限制",
		"es": "severidad mínima no válida: %q. Elija entre low / medium / high / critical, o déjelo vacío para no limitar", "ja": "最小重大度が正しくありません: %q。low / medium / high / critical から選択するか、空欄にすると制限しません",
	},
	"요청 본문 템플릿 문법이 올바르지 않습니다: %w": {
		"en": "the request body template syntax is not valid: %w",
		"zh": "请求正文模板语法无效：%w",
		"es": "la sintaxis de la plantilla del cuerpo de la solicitud no es válida: %w", "ja": "リクエスト本文テンプレートの構文が正しくありません: %w",
	},
	"요청 본문 템플릿을 렌더링하지 못했습니다: %w": {
		"en": "failed to render the request body template: %w",
		"zh": "渲染请求正文模板失败：%w",
		"es": "no se pudo renderizar la plantilla del cuerpo de la solicitud: %w", "ja": "リクエスト本文テンプレートのレンダリングに失敗しました: %w",
	},
	"요청 본문 템플릿의 렌더링 결과가 올바른 JSON이 아닙니다": {
		"en": "the rendered request body template is not valid JSON",
		"zh": "请求正文模板的渲染结果不是有效的 JSON",
		"es": "el resultado renderizado de la plantilla del cuerpo no es JSON válido", "ja": "リクエスト本文テンプレートのレンダリング結果が正しい JSON ではありません",
	},
	"요청 본문을 만들지 못했습니다: %w": {
		"en": "failed to build the request body: %w",
		"zh": "构造请求正文失败：%w",
		"es": "no se pudo construir el cuerpo de la solicitud: %w", "ja": "リクエスト本文を作成できませんでした: %w",
	},
	"요청 주소가 올바르지 않습니다: %s": {
		"en": "the request address is not valid: %s",
		"zh": "请求地址无效：%s",
		"es": "la dirección de la solicitud no es válida: %s", "ja": "リクエストアドレスが正しくありません: %s",
	},
	"요청에 실패했습니다: %s": {
		"en": "the request failed: %s",
		"zh": "请求失败：%s",
		"es": "la solicitud falló: %s", "ja": "リクエストに失敗しました: %s",
	},
	"응답을 읽지 못했습니다: %w": {
		"en": "failed to read the response: %w",
		"zh": "读取响应失败：%w",
		"es": "no se pudo leer la respuesta: %w", "ja": "応答を読み取れませんでした: %w",
	},
	"이 작업은 다른 보관 또는 삭제 작업을 진행하고 있습니다": {
		"en": "this task already has another archive or delete operation in progress",
		"zh": "该任务正在进行其他归档或删除操作",
		"es": "esta tarea ya tiene otra operación de archivo o borrado en curso", "ja": "このタスクでは別のアーカイブまたは削除操作が進行中です",
	},
	"이미 처리되었거나 존재하지 않는 승인입니다. 기록을 새로 고쳐 주세요": {
		"en": "this approval was already handled or does not exist. Refresh the record",
		"zh": "该审批已处理或不存在。请刷新记录",
		"es": "esta aprobación ya se procesó o no existe. Actualice el registro", "ja": "すでに処理済みか存在しない承認です。記録を更新してください",
	},
	"자격 증명 전송을 거부했습니다: 연결이 암호화되지 않았습니다. TLS 를 켜거나, 465 포트(암시적 TLS)를 사용하거나, 「TLS 사용」을 체크하세요 (%w)": {
		"en": "refused to send credentials: the connection is not encrypted. Enable TLS, use port 465 (implicit TLS), or check \"Use TLS\" (%w)",
		"zh": "拒绝发送凭据：连接未加密。请启用 TLS、使用 465 端口（隐式 TLS），或勾选「使用 TLS」（%w）",
		"es": "se rechazó el envío de credenciales: la conexión no está cifrada. Active TLS, use el puerto 465 (TLS implícito) o marque «Usar TLS» (%w)", "ja": "資格情報の送信を拒否しました: 接続が暗号化されていません。TLS を有効にするか、465 ポート（暗黙的 TLS）を使用するか、「TLS を使用」にチェックしてください（%w）",
	},
	"작업에 아직 실행 중인 에이전트가 있습니다. 먼저 일시정지한 뒤 보관을 다시 시도하세요": {
		"en": "the task still has running agents. Pause it first, then retry archiving",
		"zh": "任务中仍有正在运行的 agent。请先暂停，然后重试归档",
		"es": "la tarea aún tiene agentes en ejecución. Páusela primero y vuelva a intentar el archivo", "ja": "タスクにはまだ実行中のエージェントがあります。まず一時停止してから、アーカイブを再試行してください",
	},
	"주소를 해석할 수 없습니다 (%s)": {
		"en": "cannot resolve the address (%s)",
		"zh": "无法解析地址（%s）",
		"es": "no se puede resolver la dirección (%s)", "ja": "アドレスを解決できません（%s）",
	},
	"지원하지 않는 메서드입니다: %s (GET/POST/PUT/PATCH 중에서 선택하세요)": {
		"en": "unsupported method: %s (choose from GET/POST/PUT/PATCH)",
		"zh": "不支持的方法：%s（请从 GET/POST/PUT/PATCH 中选择）",
		"es": "método no admitido: %s (elija entre GET/POST/PUT/PATCH)", "ja": "サポートされていないメソッドです: %s（GET/POST/PUT/PATCH から選択してください）",
	},
	"한 번에 최대 100개 보관까지 처리할 수 있습니다": {
		"en": "at most 100 archives can be processed at once",
		"zh": "一次最多可处理 100 个归档",
		"es": "se pueden procesar como máximo 100 archivos a la vez", "ja": "一度に最大 100 件のアーカイブまで処理できます",
	},
	"호스트 이름이 없습니다": {
		"en": "the host name is missing",
		"zh": "主机名缺失",
		"es": "Falta el nombre del host", "ja": "ホスト名がありません",
	},
}

func init() {
	for k, v := range errCatalogExtra {
		errCatalog[k] = v
	}
}
