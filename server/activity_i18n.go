package server

import "github.com/Autumn-27/artex/config"

// 活动流摘要翻译。
//
// 这些字符串写进 exploration_activities / conversation_activities 的 Summary 与
// Detail 字段，是用户在任务时间线上直接读到的文字。它们此前全部硬编码韩文，
// 于是语言设置对时间线无效——把界面切成中文，收尾仍然显示
// "모든 목표를 달성했고 …"。
//
// 与 errmsg_catalog 同一套做法：韩文原文作 key，ko 不存（key 即原文），查不到
// 就原样返回。这样韩文部署的落库内容与改动前逐字节一致，历史记录也不会因为
// 改代码而变成另一种语言。
//
// 为什么用 ActiveLanguage 而不是 config.T：摘要是在任务运行过程中生成的，
// 用户可能已经在界面上改过语言；config.T 读的是启动期默认语言（ARTEX_LANG），
// 对运行期改动无感。
//
// 注意：这些摘要不进入任何 agent 的上下文（见 engine.go 中 emitActivity 的注释——
// 显示专用活动不带 node_id，planner/worker 不会回读），所以翻译它们不会污染
// 模型输入。
var activityCatalog = map[string]map[string]string{
	// --- 任务收尾 / 规划轮次 (engine.go) ---
	"모든 목표를 달성했고 직접 투입한 의도의 실행도 끝나, 작업을 종료합니다": {
		"en": "All goals are met and the directly injected intents have finished, so the task is ending",
		"zh": "所有目标均已达成，且直接投入的意图也已执行完毕，任务结束",
		"es": "Se han alcanzado todos los objetivos y las intenciones inyectadas directamente han terminado, así que la tarea finaliza",
		"ja": "すべての目標を達成し、直接投入した意図の実行も終わったため、タスクを終了します",
	},
	"%d차 계획 수립": {
		"en": "Planning round %d",
		"zh": "第 %d 轮规划",
		"es": "Ronda de planificación %d",
		"ja": "%d 回目の計画立案",
	},
	// --- 超时收尾 (engine_timeout.go) ---
	"작업 시간 초과 마무리·최종 판정(%d차)": {
		"en": "Task timeout wrap-up and final verdict (round %d)",
		"zh": "任务超时收尾与最终判定（第 %d 轮）",
		"es": "Cierre por tiempo de espera agotado y veredicto final (ronda %d)",
		"ja": "タスクタイムアウトの締めくくり・最終判定（%d 回目）",
	},
	// --- 目标分解 / 排队 (goals.go) ---
	"0차 목표 분해": {
		"en": "Round-0 goal breakdown",
		"zh": "第 0 轮目标分解",
		"es": "Desglose de objetivos de la ronda 0",
		"ja": "0 回目の目標分解",
	},
	"대기열 등록: 동시 실행 상한 %d개에 도달해, 빈자리가 나면 자동으로 시작합니다": {
		"en": "Queued: the concurrent-run limit of %d has been reached; it will start automatically when a slot frees up",
		"zh": "已加入队列：并发执行上限 %d 已达，出现空位后将自动开始",
		"es": "En cola: se alcanzó el límite de ejecución simultánea de %d; comenzará automáticamente cuando se libere un espacio",
		"ja": "キューに登録: 同時実行の上限 %d に達したため、空きが出次第自動的に開始します",
	},
	"대기열 등록: 현재 실행 가능한 LLM 설정이 없어, 설정이 복구되면 자동으로 시작합니다": {
		"en": "Queued: no LLM configuration is currently usable; it will start automatically once one is restored",
		"zh": "已加入队列：当前没有可用的 LLM 配置，配置恢复后将自动开始",
		"es": "En cola: no hay ninguna configuración de LLM disponible; comenzará automáticamente cuando se restablezca una",
		"ja": "キューに登録: 現在使用可能な LLM 設定がないため、設定が復旧次第自動的に開始します",
	},
	"대기열 등록: 먼저 대기 중인 작업이 있어, FIFO 순서대로 자동으로 시작합니다": {
		"en": "Queued: earlier tasks are already waiting, so it will start automatically in FIFO order",
		"zh": "已加入队列：已有任务在等待，将按 FIFO 顺序自动开始",
		"es": "En cola: ya hay tareas esperando, así que comenzará automáticamente en orden FIFO",
		"ja": "キューに登録: 先に待機中のタスクがあるため、FIFO 順に自動的に開始します",
	},
	// --- 资产来源标签 (assets.go) ---
	"자산 API 로 등록": {
		"en": "Registered via the assets API",
		"zh": "通过资产 API 注册",
		"es": "Registrado mediante la API de activos",
		"ja": "資産 API で登録",
	},
	"작업 설명 또는 목표로 초기화": {
		"en": "Initialized from the task description or goal",
		"zh": "由任务描述或目标初始化",
		"es": "Inicializado desde la descripción o el objetivo de la tarea",
		"ja": "タスクの説明または目標から初期化",
	},
	// --- 漏洞深化意图 (findings_groups.go) ---
	"사용자가 제출한 취약점 심화 익스플로잇 의도": {
		"en": "User-submitted intent to deepen exploitation of a finding",
		"zh": "用户提交的漏洞深入利用意图",
		"es": "Intención enviada por el usuario para profundizar la explotación de una vulnerabilidad",
		"ja": "ユーザーが提出した脆弱性の深掘り悪用の意図",
	},
	// --- 通知渠道自测 (notify_api.go) ---
	"ARTEX 알림 채널 테스트 메시지입니다. 이 메시지를 받으셨다면 채널 설정이 정상입니다.": {
		"en": "This is an ARTEX notification channel test message. If you received it, the channel is configured correctly.",
		"zh": "这是一条 ARTEX 通知渠道测试消息。收到此消息说明渠道配置正常。",
		"es": "Este es un mensaje de prueba del canal de notificaciones de ARTEX. Si lo recibió, el canal está configurado correctamente.",
		"ja": "これは ARTEX 通知チャネルのテストメッセージです。これを受け取っていれば、チャネル設定は正常です。",
	},
	// --- 归档恢复 (task_archives.go) ---
	"작업 복원을 완료했지만 일부 연관 항목이 온전히 복원되지 못했습니다": {
		"en": "The task was restored, but some related items could not be fully restored",
		"zh": "任务已恢复，但部分关联项未能完整恢复",
		"es": "La tarea se restauró, pero algunos elementos relacionados no pudieron restaurarse por completo",
		"ja": "タスクの復元は完了しましたが、一部の関連項目が完全には復元できませんでした",
	},
}

// actText renders one activity-stream string in the active output language.
func actText(ko string) string {
	if m, ok := activityCatalog[ko]; ok {
		if v, ok := m[config.ActiveLanguage()]; ok && v != "" {
			return v
		}
	}
	return ko
}
