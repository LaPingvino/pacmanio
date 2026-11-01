#!/usr/bin/env bash
# caddy-pacmanio.sh — serve the private Pacmanio web build at https://kiefte.eu/pacmanio/
# behind the same HTTP basic auth as the goloco test page (the build holds stock
# GameMaker assets and must not be public).
#
# Copies the user/hash line from the goloco-test block (run caddy-goloco.sh first),
# inserts (or replaces) a marked pacmanio block in the kiefte.eu site, keeps a
# Caddyfile.bak.<epoch> backup, validates, and reloads Caddy; on a validation
# failure the backup is restored and Caddy is left untouched.
#
#   sudo bash ~/Projects/pacmanio-recovered/scripts/caddy-pacmanio.sh [WEB_DIR]
set -euo pipefail

WEB=${1:-/home/joop/Projects/pacmanio-web}
CADDYFILE=/etc/caddy/Caddyfile
SITE_RE='^kiefte[.]eu, www[.]kiefte[.]eu'

[[ $EUID -eq 0 ]] || { echo "run with sudo" >&2; exit 1; }
[[ -f $WEB/index.html ]] || { echo "no index.html in $WEB" >&2; exit 1; }
grep -qE "$SITE_RE" "$CADDYFILE" || { echo "kiefte.eu site block not found in $CADDYFILE" >&2; exit 1; }

# The goloco block's credentials line ("user hash"), taken as it is.
cred=$(awk '/# goloco-test BEGIN/{b=1} b&&/basic_auth/{a=1;next} a&&NF==2{print $1" "$2; exit}' "$CADDYFILE")
[[ -n $cred ]] || { echo "no goloco-test basic_auth credentials found; run caddy-goloco.sh first" >&2; exit 1; }

chmod -R o+rX "$WEB"
backup="$CADDYFILE.bak.$(date +%s)"
cp -p "$CADDYFILE" "$backup"

block=$(cat <<EOB
        # pacmanio-test BEGIN (managed by pacmanio-recovered/scripts/caddy-pacmanio.sh)
        redir /pacmanio /pacmanio/
        handle_path /pacmanio/* {
                basic_auth {
                        $cred
                }
                root * $WEB
                file_server
        }
        # pacmanio-test END
EOB
)

tmp=$(mktemp)
awk -v block="$block" -v site="$SITE_RE" '
        /# pacmanio-test BEGIN/ { skip = 1 }
        skip { if (/# pacmanio-test END/) skip = 0; next }
        { print }
        $0 ~ site && !done { print block; done = 1 }
' "$CADDYFILE" > "$tmp"
cat "$tmp" > "$CADDYFILE"
rm -f "$tmp"

if ! caddy validate --config "$CADDYFILE" --adapter caddyfile >/dev/null 2>&1; then
        cp -p "$backup" "$CADDYFILE"
        echo "Caddyfile did not validate — restored $backup. Details:" >&2
        caddy validate --config "$CADDYFILE" --adapter caddyfile >&2 || true
        exit 1
fi
systemctl reload caddy
echo "Pacmanio test page: https://kiefte.eu/pacmanio/ (same login as goloco)"
