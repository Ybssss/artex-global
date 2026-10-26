package selfupdate

import "github.com/Autumn-27/artex/config"

// Self-update progress and failure messages.
//
// These reach the operator twice: as progress lines in the update panel and as
// the Detail field of the state the server reports when an update fails. Both are
// user-facing, so they follow the language picked in the UI rather than the
// Korean the templates were authored in.
//
// Keys are the Korean source templates; format verbs keep their positional order
// so the same fmt.Errorf/Err arguments line up after substitution. Korean is the
// key itself and is absent from the maps, so a miss or an active language of ko
// returns the source unchanged.
//
// config.ActiveLanguage reads the runtime setting (settings table) layered over
// the bootstrap default, so a language change applies without a restart.
var updateCatalog = map[string]map[string]string{
	// --- stage.go: download / verify / extract ---
	"이 버전은 %s/%s 용 릴리스 패키지를 제공하지 않습니다(%s 누락)": {
		"en": "this version has no release package for %s/%s (%s is missing)",
		"zh": "此版本未提供适用于 %s/%s 的发行包（缺少 %s）",
		"es": "esta versión no tiene un paquete de publicación para %s/%s (falta %s)",
		"ja": "このバージョンには %s/%s 用のリリースパッケージがありません（%s が欠落）",
	},
	"체크섬 목록을 가져오는 중…": {
		"en": "Fetching checksum list…",
		"zh": "正在获取校验和列表…",
		"es": "Obteniendo la lista de sumas de verificación…",
		"ja": "チェックサム一覧を取得中…",
	},
	"%s 에 %s 항목이 없어 검증되지 않은 바이너리 설치를 거부합니다": {
		"en": "%s has no entry for %s, so installing an unverified binary is refused",
		"zh": "%s 中没有 %s 条目，因此拒绝安装未经校验的二进制文件",
		"es": "%s no tiene una entrada para %s, por lo que se rechaza instalar un binario sin verificar",
		"ja": "%s に %s の項目がないため、検証されていないバイナリのインストールを拒否します",
	},
	"%s 내려받는 중(%s)…": {
		"en": "Downloading %s (%s)…",
		"zh": "正在下载 %s（%s）…",
		"es": "Descargando %s (%s)…",
		"ja": "%s をダウンロード中（%s）…",
	},
	"SHA256 검증 중…": {
		"en": "Verifying SHA256…",
		"zh": "正在校验 SHA256…",
		"es": "Verificando SHA256…",
		"ja": "SHA256 を検証中…",
	},
	"SHA256 이 일치하지 않습니다: 기대값 %s, 실제값 %s(다운로드가 손상됐거나 변조됨)": {
		"en": "SHA256 mismatch: expected %s, got %s (the download is corrupt or has been tampered with)",
		"zh": "SHA256 不匹配：期望 %s，实际 %s（下载已损坏或被篡改）",
		"es": "SHA256 no coincide: se esperaba %s, se obtuvo %s (la descarga está dañada o ha sido manipulada)",
		"ja": "SHA256 が一致しません: 期待値 %s、実際値 %s（ダウンロードが破損しているか改ざんされています）",
	},
	"압축을 풀고 스모크 테스트 중…": {
		"en": "Extracting and running a smoke test…",
		"zh": "正在解压并执行冒烟测试…",
		"es": "Extrayendo y ejecutando una prueba de humo…",
		"ja": "展開してスモークテスト中…",
	},
	"새 버전이 현재 시스템에서 실행되지 않습니다: %w": {
		"en": "the new version does not run on this system: %w",
		"zh": "新版本无法在当前系统上运行：%w",
		"es": "la nueva versión no se ejecuta en este sistema: %w",
		"ja": "新しいバージョンがこのシステムで実行できません: %w",
	},
	"새 바이너리 체크섬 계산 실패: %w": {
		"en": "failed to compute the new binary's checksum: %w",
		"zh": "计算新二进制文件校验和失败：%w",
		"es": "no se pudo calcular la suma de verificación del nuevo binario: %w",
		"ja": "新しいバイナリのチェックサム計算に失敗しました: %w",
	},
	"체크섬 기록 실패: %w": {
		"en": "failed to record the checksum: %w",
		"zh": "记录校验和失败：%w",
		"es": "no se pudo registrar la suma de verificación: %w",
		"ja": "チェックサムの記録に失敗しました: %w",
	},
	"새 버전 준비 실패: %w": {
		"en": "failed to prepare the new version: %w",
		"zh": "准备新版本失败：%w",
		"es": "no se pudo preparar la nueva versión: %w",
		"ja": "新しいバージョンの準備に失敗しました: %w",
	},
	"경고: 업데이트 표시 파일 기록에 실패해 이번 업데이트는 자동 롤백 보호를 받지 못합니다": {
		"en": "warning: failed to write the update marker file, so this update is not covered by automatic rollback protection",
		"zh": "警告：写入更新标记文件失败，本次更新不受自动回滚保护",
		"es": "advertencia: no se pudo escribir el archivo marcador de actualización, por lo que esta actualización no cuenta con protección de reversión automática",
		"ja": "警告: 更新マーカーファイルの書き込みに失敗したため、この更新は自動ロールバック保護の対象外です",
	},
	"새 버전이 준비되었습니다. 재시작 중…": {
		"en": "The new version is ready. Restarting…",
		"zh": "新版本已准备就绪，正在重启…",
		"es": "La nueva versión está lista. Reiniciando…",
		"ja": "新しいバージョンの準備ができました。再起動中…",
	},
	"이 릴리스에 %s 항목이 없어 무결성을 검증할 수 없으므로 업데이트를 거부합니다": {
		"en": "this release has no %s entry, so integrity cannot be verified and the update is refused",
		"zh": "此发行版缺少 %s 条目，无法校验完整性，因此拒绝更新",
		"es": "esta publicación no tiene una entrada %s, por lo que no se puede verificar la integridad y se rechaza la actualización",
		"ja": "このリリースには %s の項目がないため整合性を検証できず、更新を拒否します",
	},
	"%s 내려받기 실패: %w": {
		"en": "failed to download %s: %w",
		"zh": "下载 %s 失败：%w",
		"es": "no se pudo descargar %s: %w",
		"ja": "%s のダウンロードに失敗しました: %w",
	},
	"%s 읽기 실패: %w": {
		"en": "failed to read %s: %w",
		"zh": "读取 %s 失败：%w",
		"es": "no se pudo leer %s: %w",
		"ja": "%s の読み取りに失敗しました: %w",
	},
	"%s 내용이 비어 있거나 형식을 인식할 수 없습니다": {
		"en": "%s is empty or its format is unrecognised",
		"zh": "%s 内容为空或格式无法识别",
		"es": "%s está vacío o su formato no se reconoce",
		"ja": "%s の内容が空か、形式を認識できません",
	},
	"임시 파일 생성 실패: %w": {
		"en": "failed to create a temporary file: %w",
		"zh": "创建临时文件失败：%w",
		"es": "no se pudo crear un archivo temporal: %w",
		"ja": "一時ファイルの作成に失敗しました: %w",
	},
	"다운로드가 중단되었습니다: %w": {
		"en": "the download was interrupted: %w",
		"zh": "下载已中断：%w",
		"es": "la descarga se interrumpió: %w",
		"ja": "ダウンロードが中断されました: %w",
	},
	"디스크 기록 실패: %w": {
		"en": "failed to write to disk: %w",
		"zh": "写入磁盘失败：%w",
		"es": "no se pudo escribir en el disco: %w",
		"ja": "ディスクへの書き込みに失敗しました: %w",
	},
	"다운로드가 완전하지 않습니다: 기대 %d바이트, 실제 %d바이트": {
		"en": "the download is incomplete: expected %d bytes, got %d bytes",
		"zh": "下载不完整：期望 %d 字节，实际 %d 字节",
		"es": "la descarga está incompleta: se esperaban %d bytes, se obtuvieron %d bytes",
		"ja": "ダウンロードが不完全です: 期待 %d バイト、実際 %d バイト",
	},
	"릴리스 패키지 열기 실패: %w": {
		"en": "failed to open the release package: %w",
		"zh": "打开发行包失败：%w",
		"es": "no se pudo abrir el paquete de publicación: %w",
		"ja": "リリースパッケージを開くのに失敗しました: %w",
	},
	"새 바이너리 쓰기 실패: %w": {
		"en": "failed to write the new binary: %w",
		"zh": "写入新二进制文件失败：%w",
		"es": "no se pudo escribir el nuevo binario: %w",
		"ja": "新しいバイナリの書き込みに失敗しました: %w",
	},
	"%s 압축 해제 실패: %w": {
		"en": "failed to extract %s: %w",
		"zh": "解压 %s 失败：%w",
		"es": "no se pudo extraer %s: %w",
		"ja": "%s の展開に失敗しました: %w",
	},
	"릴리스 패키지 안의 실행 파일이 %s 보다 커서 압축 해제를 거부합니다": {
		"en": "the executable inside the release package is larger than %s, so extraction is refused",
		"zh": "发行包内的可执行文件大于 %s，因此拒绝解压",
		"es": "el ejecutable dentro del paquete de publicación es mayor que %s, por lo que se rechaza la extracción",
		"ja": "リリースパッケージ内の実行ファイルが %s より大きいため、展開を拒否します",
	},
	"릴리스 패키지 안의 %s 파일이 비어 있습니다": {
		"en": "the %s file inside the release package is empty",
		"zh": "发行包内的 %s 文件为空",
		"es": "el archivo %s dentro del paquete de publicación está vacío",
		"ja": "リリースパッケージ内の %s ファイルが空です",
	},
	"릴리스 패키지에서 %s 파일을 찾지 못했습니다": {
		"en": "could not find the %s file in the release package",
		"zh": "在发行包中找不到 %s 文件",
		"es": "no se encontró el archivo %s en el paquete de publicación",
		"ja": "リリースパッケージ内に %s ファイルが見つかりませんでした",
	},
	"프로그램 디렉터리 %s 에 쓸 수 없어 자동 업데이트를 할 수 없습니다(권한을 확인하거나 수동 업데이트를 사용하세요): %w": {
		"en": "cannot write to the program directory %s, so automatic update is unavailable (check permissions or use a manual update): %w",
		"zh": "无法写入程序目录 %s，因此无法自动更新（请检查权限或使用手动更新）：%w",
		"es": "no se puede escribir en el directorio del programa %s, por lo que la actualización automática no está disponible (compruebe los permisos o use una actualización manual): %w",
		"ja": "プログラムディレクトリ %s に書き込めないため自動更新できません（権限を確認するか手動更新を使用してください）: %w",
	},
	"내려받는 중 %s / %s": {
		"en": "Downloading %s / %s",
		"zh": "正在下载 %s / %s",
		"es": "Descargando %s / %s",
		"ja": "ダウンロード中 %s / %s",
	},

	// --- bootstrap.go: swap / rollback ---
	"새 버전 검증에 실패해 폐기했습니다: ": {
		"en": "verification of the new version failed, so it was discarded: ",
		"zh": "新版本校验失败，已丢弃：",
		"es": "la verificación de la nueva versión falló, por lo que se descartó: ",
		"ja": "新しいバージョンの検証に失敗したため破棄しました: ",
	},
	"버전 교체에 실패했습니다: ": {
		"en": "failed to swap versions: ",
		"zh": "版本替换失败：",
		"es": "no se pudo intercambiar las versiones: ",
		"ja": "バージョンの入れ替えに失敗しました: ",
	},
	"새 버전 시작에 실패하고 롤백도 실패했습니다: ": {
		"en": "the new version failed to start and the rollback also failed: ",
		"zh": "新版本启动失败，回滚也失败：",
		"es": "la nueva versión no arrancó y la reversión también falló: ",
		"ja": "新しいバージョンの起動に失敗し、ロールバックも失敗しました: ",
	},
	"새 버전 시작에 실패해 이전 버전으로 롤백했습니다(%s)": {
		"en": "the new version failed to start, so it rolled back to the previous version (%s)",
		"zh": "新版本启动失败，已回滚到上一版本（%s）",
		"es": "la nueva versión no arrancó, por lo que se revirtió a la versión anterior (%s)",
		"ja": "新しいバージョンの起動に失敗したため、以前のバージョンにロールバックしました（%s）",
	},
	"체크섬 읽기 실패: %w": {
		"en": "failed to read the checksum: %w",
		"zh": "读取校验和失败：%w",
		"es": "no se pudo leer la suma de verificación: %w",
		"ja": "チェックサムの読み取りに失敗しました: %w",
	},
	"체크섬 계산 실패: %w": {
		"en": "failed to compute the checksum: %w",
		"zh": "计算校验和失败：%w",
		"es": "no se pudo calcular la suma de verificación: %w",
		"ja": "チェックサムの計算に失敗しました: %w",
	},
	"SHA256 이 일치하지 않습니다(다운로드가 손상됐거나 변조됨)": {
		"en": "SHA256 mismatch (the download is corrupt or has been tampered with)",
		"zh": "SHA256 不匹配（下载已损坏或被篡改）",
		"es": "SHA256 no coincide (la descarga está dañada o ha sido manipulada)",
		"ja": "SHA256 が一致しません（ダウンロードが破損しているか改ざんされています）",
	},
	"실행 권한 부여 실패: %w": {
		"en": "failed to grant execute permission: %w",
		"zh": "授予执行权限失败：%w",
		"es": "no se pudo conceder permiso de ejecución: %w",
		"ja": "実行権限の付与に失敗しました: %w",
	},
	"스모크 테스트 시간 초과(새 바이너리가 응답하지 않음)": {
		"en": "smoke test timed out (the new binary did not respond)",
		"zh": "冒烟测试超时（新二进制文件无响应）",
		"es": "la prueba de humo agotó el tiempo de espera (el nuevo binario no respondió)",
		"ja": "スモークテストがタイムアウトしました（新しいバイナリが応答しません）",
	},
	"스모크 테스트 실패: %v: %s": {
		"en": "smoke test failed: %v: %s",
		"zh": "冒烟测试失败：%v：%s",
		"es": "la prueba de humo falló: %v: %s",
		"ja": "スモークテストに失敗しました: %v: %s",
	},
	"이전 백업 %s 정리 실패: %w": {
		"en": "failed to clean up the previous backup %s: %w",
		"zh": "清理上一备份 %s 失败：%w",
		"es": "no se pudo limpiar la copia de seguridad anterior %s: %w",
		"ja": "以前のバックアップ %s の整理に失敗しました: %w",
	},
	"현재 버전 백업 실패: %w": {
		"en": "failed to back up the current version: %w",
		"zh": "备份当前版本失败：%w",
		"es": "no se pudo respaldar la versión actual: %w",
		"ja": "現在のバージョンのバックアップに失敗しました: %w",
	},
	"새 버전 설치 실패(%v), 현재 버전 복구도 실패: %w": {
		"en": "installing the new version failed (%v), and restoring the current version also failed: %w",
		"zh": "安装新版本失败（%v），恢复当前版本也失败：%w",
		"es": "la instalación de la nueva versión falló (%v) y la restauración de la versión actual también falló: %w",
		"ja": "新しいバージョンのインストールに失敗し（%v）、現在のバージョンの復元にも失敗しました: %w",
	},
	"새 버전 설치 실패: %w": {
		"en": "failed to install the new version: %w",
		"zh": "安装新版本失败：%w",
		"es": "no se pudo instalar la nueva versión: %w",
		"ja": "新しいバージョンのインストールに失敗しました: %w",
	},
	"롤백할 백업이 없습니다(%s): %w": {
		"en": "there is no backup to roll back to (%s): %w",
		"zh": "没有可回滚的备份（%s）：%w",
		"es": "no hay copia de seguridad a la que revertir (%s): %w",
		"ja": "ロールバックするバックアップがありません（%s）: %w",
	},
	"실패한 버전 이동 실패: %w": {
		"en": "failed to move the failed version aside: %w",
		"zh": "移动失败的版本失败：%w",
		"es": "no se pudo apartar la versión fallida: %w",
		"ja": "失敗したバージョンの移動に失敗しました: %w",
	},
	"이전 버전 복구 실패: %w": {
		"en": "failed to restore the previous version: %w",
		"zh": "恢复上一版本失败：%w",
		"es": "no se pudo restaurar la versión anterior: %w",
		"ja": "以前のバージョンの復元に失敗しました: %w",
	},
	"롤백할 이전 버전이 없습니다(": {
		"en": "there is no previous version to roll back to (",
		"zh": "没有可回滚的上一版本（",
		"es": "no hay una versión anterior a la que revertir (",
		"ja": "ロールバックする以前のバージョンがありません（",
	},
	" 없음)": {
		"en": " not found)",
		"zh": " 不存在）",
		"es": " no encontrado)",
		"ja": " なし）",
	},
	"이전 버전을 실행할 수 없어 롤백을 거부합니다: %w": {
		"en": "the previous version cannot be executed, so the rollback is refused: %w",
		"zh": "上一版本无法执行，因此拒绝回滚：%w",
		"es": "la versión anterior no se puede ejecutar, por lo que se rechaza la reversión: %w",
		"ja": "以前のバージョンを実行できないため、ロールバックを拒否します: %w",
	},
	"현재 버전 이동 실패: %w": {
		"en": "failed to move the current version: %w",
		"zh": "移动当前版本失败：%w",
		"es": "no se pudo mover la versión actual: %w",
		"ja": "現在のバージョンの移動に失敗しました: %w",
	},
	"이전 버전 설치 실패: %w": {
		"en": "failed to install the previous version: %w",
		"zh": "安装上一版本失败：%w",
		"es": "no se pudo instalar la versión anterior: %w",
		"ja": "以前のバージョンのインストールに失敗しました: %w",
	},
	"알 수 없는 버전": {
		"en": "Unknown version",
		"zh": "未知版本",
		"es": "Versión desconocida",
		"ja": "不明なバージョン",
	},

	// --- github.go ---
	"리다이렉트 횟수가 너무 많습니다": {
		"en": "too many redirects",
		"zh": "重定向次数过多",
		"es": "demasiadas redirecciones",
		"ja": "リダイレクトが多すぎます",
	},
	"HTTPS 가 아닌 주소를 거부합니다: %s": {
		"en": "refusing a non-HTTPS address: %s",
		"zh": "拒绝非 HTTPS 地址：%s",
		"es": "se rechaza una dirección que no es HTTPS: %s",
		"ja": "HTTPS 以外のアドレスを拒否します: %s",
	},
	"GitHub 도메인이 아닌 주소를 거부합니다: %s": {
		"en": "refusing an address that is not a GitHub domain: %s",
		"zh": "拒绝非 GitHub 域名的地址：%s",
		"es": "se rechaza una dirección que no es un dominio de GitHub: %s",
		"ja": "GitHub ドメイン以外のアドレスを拒否します: %s",
	},

	// --- selfupdate.go ---
	"실행 파일 위치 확인 실패: %w": {
		"en": "failed to determine the executable location: %w",
		"zh": "确定可执行文件位置失败：%w",
		"es": "no se pudo determinar la ubicación del ejecutable: %w",
		"ja": "実行ファイルの位置確認に失敗しました: %w",
	},
	"GitHub 에 접속하지 못했습니다(시스템 설정에서 전역 프록시를 설정할 수 있습니다): %w": {
		"en": "could not reach GitHub (a global proxy can be configured in the system settings): %w",
		"zh": "无法连接 GitHub（可在系统设置中配置全局代理）：%w",
		"es": "no se pudo conectar a GitHub (puede configurar un proxy global en los ajustes del sistema): %w",
		"ja": "GitHub に接続できませんでした（システム設定でグローバルプロキシを設定できます）: %w",
	},
	"GitHub API 요청 제한에 걸렸습니다(시간당 60회). 잠시 후 다시 시도하세요": {
		"en": "the GitHub API rate limit was hit (60 requests per hour). Try again shortly",
		"zh": "已触发 GitHub API 请求限制（每小时 60 次）。请稍后重试",
		"es": "se alcanzó el límite de solicitudes de la API de GitHub (60 por hora). Inténtelo de nuevo en un momento",
		"ja": "GitHub API のリクエスト制限（1 時間あたり 60 回）に達しました。しばらくしてから再試行してください",
	},
	"저장소 %s 에 아직 정식 릴리스가 없습니다": {
		"en": "repository %s has no official release yet",
		"zh": "仓库 %s 尚无正式发行版",
		"es": "el repositorio %s aún no tiene una publicación oficial",
		"ja": "リポジトリ %s にはまだ正式リリースがありません",
	},
	"GitHub 응답 코드 %d": {
		"en": "GitHub responded with status %d",
		"zh": "GitHub 响应状态码 %d",
		"es": "GitHub respondió con el código %d",
		"ja": "GitHub のレスポンスコード %d",
	},
	"Release 해석 실패: %w": {
		"en": "failed to parse the release: %w",
		"zh": "解析发行版失败：%w",
		"es": "no se pudo analizar la publicación: %w",
		"ja": "リリースの解析に失敗しました: %w",
	},
	"Release 에 tag 가 없습니다": {
		"en": "the release has no tag",
		"zh": "发行版缺少 tag",
		"es": "la publicación no tiene tag",
		"ja": "リリースに tag がありません",
	},
}

// trMsg renders one update message in the active output language.
func trMsg(ko string) string {
	if m, ok := updateCatalog[ko]; ok {
		if v, ok := m[config.ActiveLanguage()]; ok && v != "" {
			return v
		}
	}
	return ko
}
