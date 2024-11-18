from alpine:3.20

run apk add sqlite

copy entrypoint.sh /
copy lints.sql /
run chmod +x /entrypoint.sh

entrypoint ["/entrypoint.sh"]
