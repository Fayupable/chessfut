#!/bin/bash
set -euo pipefail

# ---------------------------------------------------------------------------
# Configuration
# ---------------------------------------------------------------------------
ENVIRONMENT="${BACKUP_ENV:-PROD}"
CONTAINER_NAME="${POSTGRES_CONTAINER:-chessfut_postgres}"
BACKUP_BASE_DIR="${BACKUP_BASE_DIR:-/home/enis/backups/chessfut-prod}"
RETENTION_DAYS="${RETENTION_DAYS:-30}"
MIN_DISK_MB="${MIN_DISK_MB:-500}"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_DIR="${BACKUP_BASE_DIR}/${TIMESTAMP}"
AUDIT_LOG="${BACKUP_BASE_DIR}/backup_audit.log"
ALERT_EMAIL="${ALERT_EMAIL:-enisyaman4@gmail.com}"
APP_NAME="${APP_NAME:-Chessfut}"
MAIL_FROM="${MAIL_FROM:-noreply@chessfut.com}"
SCRIPT_START=$(date +%s)

DATABASES=(
    "${POSTGRES_DB:-chessfut}"
)

R2_REMOTE="${R2_REMOTE:-}"
RCLONE_CONFIG="${RCLONE_CONFIG:-${HOME}/.config/rclone/rclone.conf}"

# ---------------------------------------------------------------------------
# Logging
# ---------------------------------------------------------------------------
log() {
    local level="${2:-INFO}"
    local message="[$(date '+%Y-%m-%d %H:%M:%S')] [${ENVIRONMENT}] [${level}] $1"
    echo "$message"
    echo "$message" >> "$AUDIT_LOG"
}

# ---------------------------------------------------------------------------
# Mail
# ---------------------------------------------------------------------------
send_mail() {
    local subject="$1"
    local status="$2"
    local summary_rows="$3"
    local extra_notes="$4"

    local header_color="#27ae60"
    local status_badge_color="#27ae60"
    [[ "$status" == "WARNING" ]] && header_color="#e67e22" && status_badge_color="#e67e22"
    [[ "$status" == "FAILED"  ]] && header_color="#c0392b" && status_badge_color="#c0392b"

    local total_elapsed=$(( $(date +%s) - SCRIPT_START ))
    local total_min=$(( total_elapsed / 60 ))
    local total_sec=$(( total_elapsed % 60 ))

    local notes_html=""
    if [ -n "$extra_notes" ]; then
        notes_html="<div style='margin-top:16px;padding:12px;background:#fff8e1;border-left:4px solid #f39c12;font-size:12px;color:#555'>${extra_notes}</div>"
    fi

    local html_body
    html_body=$(cat <<-HTML
<html>
<body style="font-family:Arial,sans-serif;max-width:640px;margin:0 auto;padding:20px;background:#f4f4f4">
  <div style="background:${header_color};padding:18px 24px;border-radius:8px 8px 0 0">
    <h2 style="color:white;margin:0;font-size:17px">${APP_NAME} — Database Backup Report</h2>
    <p style="color:rgba(255,255,255,0.85);margin:4px 0 0;font-size:12px">${ENVIRONMENT} · $(date '+%Y-%m-%d %H:%M:%S')</p>
  </div>
  <div style="background:#fff;padding:24px;border:1px solid #ddd;border-top:none;border-radius:0 0 8px 8px">
    <div style="margin-bottom:20px">
      <span style="display:inline-block;background:${status_badge_color};color:white;padding:4px 12px;border-radius:20px;font-size:13px;font-weight:bold">${status}</span>
      <span style="margin-left:12px;color:#555;font-size:13px">Total duration: <strong>${total_min}m ${total_sec}s</strong></span>
    </div>
    <table style="width:100%;border-collapse:collapse;font-size:13px">
      <thead>
        <tr style="background:#f0f0f0">
          <th style="text-align:left;padding:8px 10px;border-bottom:2px solid #ddd">Database</th>
          <th style="text-align:right;padding:8px 10px;border-bottom:2px solid #ddd">Size</th>
          <th style="text-align:right;padding:8px 10px;border-bottom:2px solid #ddd">Duration</th>
          <th style="text-align:center;padding:8px 10px;border-bottom:2px solid #ddd">Status</th>
        </tr>
      </thead>
      <tbody>
        ${summary_rows}
      </tbody>
    </table>
    ${notes_html}
    <hr style="border:none;border-top:1px solid #eee;margin:20px 0">
    <p style="color:#aaa;font-size:11px;margin:0">
      Audit log: ${AUDIT_LOG}<br>
      Backup dir: ${BACKUP_DIR}
    </p>
  </div>
</body>
</html>
HTML
)

    {
        echo "To: $ALERT_EMAIL"
        echo "From: ${APP_NAME} Backup <${MAIL_FROM}>"
        echo "Subject: $subject"
        echo "MIME-Version: 1.0"
        echo "Content-Type: text/html; charset=UTF-8"
        echo ""
        echo "$html_body"
    } | msmtp "$ALERT_EMAIL"
}

# ---------------------------------------------------------------------------
# Error handler
# ---------------------------------------------------------------------------
SUMMARY_ROWS=""
error_handler() {
    log "Backup failed at line $1" "ERROR"
    send_mail "[${ENVIRONMENT}] Backup FAILED — $(date +%Y-%m-%d)" \
        "FAILED" \
        "${SUMMARY_ROWS}" \
        "Script exited unexpectedly at line $1. Check audit log for details."
    exit 1
}
trap 'error_handler $LINENO' ERR

# ---------------------------------------------------------------------------
# Load env
# ---------------------------------------------------------------------------
ENV_FILE="$(dirname "$0")/.env"
if [ -f "$ENV_FILE" ]; then
    source "$ENV_FILE"
else
    log ".env file not found at $ENV_FILE — relying on environment variables" "WARN"
fi

# ---------------------------------------------------------------------------
# Pre-flight checks
# ---------------------------------------------------------------------------
mkdir -p "$BACKUP_DIR"
log "========== Backup started =========="

# Disk space check
AVAILABLE_MB=$(df -m "$BACKUP_BASE_DIR" | awk 'NR==2 {print $4}')
if [ "$AVAILABLE_MB" -lt "$MIN_DISK_MB" ]; then
    log "Insufficient disk space: ${AVAILABLE_MB}MB available, ${MIN_DISK_MB}MB required" "ERROR"
    exit 1
fi
log "Disk space OK: ${AVAILABLE_MB}MB available"

# Container check
if ! docker inspect "$CONTAINER_NAME" &>/dev/null; then
    log "Container $CONTAINER_NAME not found" "ERROR"
    exit 1
fi

if ! docker exec "$CONTAINER_NAME" pg_isready -U "$POSTGRES_USER" &>/dev/null; then
    log "PostgreSQL not ready in $CONTAINER_NAME" "ERROR"
    exit 1
fi
log "PostgreSQL health check passed"

# ---------------------------------------------------------------------------
# Backup
# ---------------------------------------------------------------------------
perform_backup() {
    local db_name="$1"
    local backup_file="${BACKUP_DIR}/${db_name}_${TIMESTAMP}.sql.gz"
    local checksum_file="${backup_file}.md5"
    local db_start=$(date +%s)

    log "Starting backup: $db_name"

    docker exec "$CONTAINER_NAME" pg_dump \
        -U "$POSTGRES_USER" \
        --no-owner \
        --no-privileges \
        "$db_name" | gzip > "$backup_file"

    if [ ! -s "$backup_file" ]; then
        log "Backup file is empty: $db_name" "ERROR"
        return 1
    fi

    if ! gzip -t "$backup_file" 2>/dev/null; then
        log "Gzip integrity check failed: $db_name" "ERROR"
        return 1
    fi

    md5sum "$backup_file" > "$checksum_file"

    local db_elapsed=$(( $(date +%s) - db_start ))
    local db_min=$(( db_elapsed / 60 ))
    local db_sec=$(( db_elapsed % 60 ))
    local size_human=$(du -h "$backup_file" | cut -f1)

    log "SUCCESS: $db_name — ${size_human} in ${db_min}m${db_sec}s"

    SUMMARY_ROWS+="<tr>
      <td style='padding:8px 10px;border-bottom:1px solid #eee'>${db_name}</td>
      <td style='padding:8px 10px;border-bottom:1px solid #eee;text-align:right'>${size_human}</td>
      <td style='padding:8px 10px;border-bottom:1px solid #eee;text-align:right'>${db_min}m ${db_sec}s</td>
      <td style='padding:8px 10px;border-bottom:1px solid #eee;text-align:center'>
        <span style='color:#27ae60;font-weight:bold'>✓</span>
      </td>
    </tr>"
}

for db in "${DATABASES[@]}"; do
    perform_backup "$db"
done

# ---------------------------------------------------------------------------
# R2 Upload (optional)
# ---------------------------------------------------------------------------
EXTRA_NOTES=""
if [ -n "$R2_REMOTE" ] && command -v rclone &>/dev/null && [ -f "$RCLONE_CONFIG" ]; then
    log "Uploading to R2: ${R2_REMOTE}/${TIMESTAMP}"
    if rclone copy "$BACKUP_DIR" "${R2_REMOTE}/${TIMESTAMP}" \
        --config "$RCLONE_CONFIG" \
        --progress 2>> "$AUDIT_LOG"; then
        log "R2 upload successful"
    else
        log "R2 upload failed — local backup still intact" "WARN"
        EXTRA_NOTES="R2 upload failed. Local backup is available at: ${BACKUP_DIR}"
    fi
else
    log "R2 upload skipped (rclone not configured or R2_REMOTE not set)" "INFO"
    EXTRA_NOTES="R2 upload skipped — configure R2_REMOTE and rclone to enable."
fi

# ---------------------------------------------------------------------------
# Retention cleanup
# ---------------------------------------------------------------------------
log "Cleaning up backups older than ${RETENTION_DAYS} days"
find "$BACKUP_BASE_DIR" -maxdepth 1 -type d -mtime +${RETENTION_DAYS} -exec rm -rf {} + 2>/dev/null || true

# ---------------------------------------------------------------------------
# Final report
# ---------------------------------------------------------------------------
TOTAL_SIZE=$(du -sh "$BACKUP_DIR" | cut -f1)
FINAL_STATUS="SUCCESS"
[ -n "$EXTRA_NOTES" ] && FINAL_STATUS="WARNING"

log "========== Backup completed. Total size: ${TOTAL_SIZE} =========="

send_mail "[${ENVIRONMENT}] Backup ${FINAL_STATUS} — $(date +%Y-%m-%d)" \
    "$FINAL_STATUS" \
    "$SUMMARY_ROWS" \
    "$EXTRA_NOTES"