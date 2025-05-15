#!/bin/bash
set -e

CONTAINER_NAMES=("influxdb2")

log() {
  echo "[INFO] $1" >&2
}

warn() {
  echo "[WARN] $1" >&2
}

collect_host_entries() {
  local ENTRIES=()
  for CONTAINER in "${CONTAINER_NAMES[@]}"; do
    IP=$(docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' "$CONTAINER" 2>/dev/null || true)
    if [[ -n "$IP" ]]; then
      log "容器 $CONTAINER 的 IP: $IP"
      ENTRIES+=("$IP $CONTAINER")
    else
      warn "容器 $CONTAINER 不存在或未运行"
    fi
  done
  printf "%s\n" "${ENTRIES[@]}"
}

generate_ps1_script() {
  local ACTION=$1
  shift
  local ENTRIES=("$@")

  TMP_PS1=$(mktemp --suffix=.ps1)
  log "生成 PowerShell 脚本: $TMP_PS1"

  # 写入 UTF-8 BOM 头（无换行符）
  printf "\xEF\xBB\xBF" > "$TMP_PS1"

  # 追加脚本内容
  {
    echo '$ErrorActionPreference = "Stop"'
    echo '$hostsPath = "$env:SystemRoot\System32\drivers\etc\hosts"'
    echo 'Write-Host "目标 hosts 文件: $hostsPath"'
    echo '$original = @(Get-Content -Path $hostsPath -ErrorAction SilentlyContinue)'

    for CONTAINER in "${CONTAINER_NAMES[@]}"; do
      echo "# 移除容器 $CONTAINER 的所有旧条目"
      echo "\$original = @(\$original | Where-Object { \$_ -notmatch '^\\S+\\s+${CONTAINER}\\b' })"
    done

    # 添加新条目（仅 up 操作）
    if [[ "$ACTION" == "add" ]]; then
      for entry in "${ENTRIES[@]}"; do
        echo "if (-not [string]::IsNullOrEmpty('${entry}')) { \$original += '${entry}' }"
      done
    fi

    # 过滤空行并控制格式
    echo '# 移除所有空行和注释（可选）'
    echo "\$original = \$original | Where-Object { \$_.Trim() -ne '' -and -not \$_.StartsWith('#') }"
    echo '# 添加新条目后重新生成内容'
    echo '$content = $original -join "`r`n"'
    echo '$content = $content.TrimEnd("`r`n")'
    echo '[System.IO.File]::WriteAllText($hostsPath, $content)'

    echo 'ipconfig /flushdns | Out-Null'
    echo 'Write-Host "✅ hosts 文件已更新完成！"'
    echo 'Pause'
    echo '# 自删除临时文件'
    echo 'Remove-Item -Path $MyInvocation.MyCommand.Path -Force -ErrorAction SilentlyContinue'
  } >> "$TMP_PS1"

  WIN_PATH=$(wslpath -w "$TMP_PS1")
  log "以管理员身份运行 PowerShell（查看 hosts 更新日志）"
  powershell.exe -Command "Start-Process powershell -Verb RunAs -ArgumentList '-NoExit','-ExecutionPolicy','Bypass','-File','${WIN_PATH}'"
}

ACTION="$1"
shift

case "$ACTION" in
  up)
    log "启动容器..."
    docker-compose up -d "$@"

    max_retries=5
    retry_interval=2
    HOST_ENTRIES=()
    for ((i=1; i<=max_retries; i++)); do
      mapfile -t HOST_ENTRIES < <(collect_host_entries 2>/dev/null)
      if [[ "${#HOST_ENTRIES[@]}" -gt 0 ]]; then
        break
      fi
      log "等待容器IP分配，重试 $i/$max_retries..."
      sleep $retry_interval
    done

    if [[ "${#HOST_ENTRIES[@]}" -gt 0 ]]; then
      generate_ps1_script add "${HOST_ENTRIES[@]}"
    else
      warn "未找到容器 IP，跳过 hosts 映射添加"
    fi
    ;;

  down)
    log "准备移除 hosts 映射..."
    mapfile -t HOST_ENTRIES < <(collect_host_entries 2>/dev/null)

    if [[ "${#HOST_ENTRIES[@]}" -gt 0 ]]; then
      generate_ps1_script remove "${HOST_ENTRIES[@]}"
    else
      warn "无法获取容器 IP，可能容器已停止"
    fi

    log "停止容器..."
    docker-compose down -v "$@"
    ;;

  *)
    echo "用法: $0 {up|down}"
    exit 1
    ;;
esac