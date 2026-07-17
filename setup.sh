#!/usr/bin/env bash

set -Eeuo pipefail

readonly SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
readonly COMPOSE_FILE="${SCRIPT_DIR}/docker-compose.yml"
readonly DOMAIN="kaucjo-mapa.pl"
readonly HOST_IP="127.0.0.1"
readonly HOSTS_FILE="/etc/hosts"
readonly HOSTS_MARKER="# kaucjomapa-dev"
readonly HOSTS_ENTRY="${HOST_IP} ${DOMAIN} ${HOSTS_MARKER}"
readonly CERTIFICATES_DIR="${SCRIPT_DIR}/.docker/nginx-proxy-manager/certs"
readonly CERTIFICATE_FILE="${CERTIFICATES_DIR}/kaucjo-mapa.crt"
readonly CERTIFICATE_KEY_FILE="${CERTIFICATES_DIR}/kaucjo-mapa.key"
readonly CERTIFICATE_CA_FILE="${CERTIFICATES_DIR}/kaucjo-mapa-ca.crt"

usage() {
	printf 'Usage: %s <start|stop>\n' "$(basename -- "$0")" >&2
}

require_command() {
	if ! command -v "$1" >/dev/null 2>&1; then
		printf 'Error: required command not found: %s\n' "$1" >&2
		exit 1
	fi
}

run_as_root() {
	if [[ ${EUID} -eq 0 ]]; then
		"$@"
	else
		sudo "$@"
	fi
}

require_hosts_access() {
	if [[ -w ${HOSTS_FILE} ]]; then
		return
	fi

	require_command sudo
	sudo -v
}

domain_is_mapped_to_localhost() {
	awk -v domain="${DOMAIN}" -v ip="${HOST_IP}" '
		$0 !~ /^[[:space:]]*#/ && $1 == ip {
			for (field = 2; field <= NF; field++) {
				if ($field ~ /^#/) break
				if ($field == domain) found = 1
			}
		}
		END { exit !found }
	' "${HOSTS_FILE}"
}

domain_has_another_mapping() {
	awk -v domain="${DOMAIN}" -v ip="${HOST_IP}" '
		$0 !~ /^[[:space:]]*#/ && $1 != ip {
			for (field = 2; field <= NF; field++) {
				if ($field ~ /^#/) break
				if ($field == domain) found = 1
			}
		}
		END { exit !found }
	' "${HOSTS_FILE}"
}

add_host_entry() {
	if domain_is_mapped_to_localhost; then
		printf '%s is already mapped to %s in %s.\n' "${DOMAIN}" "${HOST_IP}" "${HOSTS_FILE}"
		return
	fi

	if domain_has_another_mapping; then
		printf 'Error: %s already has a non-%s mapping in %s.\n' "${DOMAIN}" "${HOST_IP}" "${HOSTS_FILE}" >&2
		printf 'Remove the conflicting mapping manually and run this command again.\n' >&2
		return 1
	fi

	printf '%s\n' "${HOSTS_ENTRY}" | run_as_root tee -a "${HOSTS_FILE}" >/dev/null
	printf 'Added %s to %s.\n' "${HOSTS_ENTRY}" "${HOSTS_FILE}"
}

remove_host_entry() {
	local temporary_file
	temporary_file="$(mktemp)"

	awk -v marker="${HOSTS_MARKER}" 'index($0, marker) == 0' "${HOSTS_FILE}" >"${temporary_file}"

	if cmp -s -- "${HOSTS_FILE}" "${temporary_file}"; then
		rm -f -- "${temporary_file}"
		printf 'No managed %s entry found in %s.\n' "${DOMAIN}" "${HOSTS_FILE}"
		return
	fi

	run_as_root sh -c 'cat -- "$1" >"$2"' sh "${temporary_file}" "${HOSTS_FILE}"
	rm -f -- "${temporary_file}"
	printf 'Removed the managed %s entry from %s.\n' "${DOMAIN}" "${HOSTS_FILE}"
}

certificate_is_valid() {
	local certificate_public_key
	local private_key_public_key

	[[ -s ${CERTIFICATE_FILE} && -s ${CERTIFICATE_KEY_FILE} && -s ${CERTIFICATE_CA_FILE} ]] || return 1
	openssl x509 -in "${CERTIFICATE_FILE}" -noout -checkend 86400 >/dev/null 2>&1 || return 1
	openssl x509 -in "${CERTIFICATE_FILE}" -noout -checkhost "${DOMAIN}" >/dev/null 2>&1 || return 1

	certificate_public_key="$(openssl x509 -in "${CERTIFICATE_FILE}" -pubkey -noout 2>/dev/null | openssl pkey -pubin -outform PEM 2>/dev/null)" || return 1
	private_key_public_key="$(openssl pkey -in "${CERTIFICATE_KEY_FILE}" -pubout -outform PEM 2>/dev/null)" || return 1

	[[ ${certificate_public_key} == "${private_key_public_key}" ]]
}

generate_local_certificate() {
	mkdir -p -- "${CERTIFICATES_DIR}"
	rm -f -- "${CERTIFICATE_FILE}" "${CERTIFICATE_KEY_FILE}"

	if command -v mkcert >/dev/null 2>&1; then
		mkcert -cert-file "${CERTIFICATE_FILE}" -key-file "${CERTIFICATE_KEY_FILE}" "${DOMAIN}"
		cp -- "$(mkcert -CAROOT)/rootCA.pem" "${CERTIFICATE_CA_FILE}"
		printf 'Generated a locally trusted TLS certificate for %s with mkcert.\n' "${DOMAIN}"
	else
		openssl req -x509 -newkey rsa:2048 -sha256 -nodes -days 825 \
			-keyout "${CERTIFICATE_KEY_FILE}" \
			-out "${CERTIFICATE_FILE}" \
			-subj "/CN=${DOMAIN}" \
			-addext "subjectAltName=DNS:${DOMAIN}" \
			-addext "keyUsage=digitalSignature,keyEncipherment" \
			-addext "extendedKeyUsage=serverAuth"
		cp -- "${CERTIFICATE_FILE}" "${CERTIFICATE_CA_FILE}"
		printf 'Generated a self-signed TLS certificate for %s with OpenSSL.\n' "${DOMAIN}"
		printf 'Warning: install mkcert before regenerating the certificate to avoid browser trust warnings.\n' >&2
	fi

	chmod 600 "${CERTIFICATE_KEY_FILE}"
	chmod 644 "${CERTIFICATE_FILE}" "${CERTIFICATE_CA_FILE}"
}

ensure_local_certificate() {
	require_command openssl

	if certificate_is_valid; then
		printf 'TLS certificate for %s is already valid.\n' "${DOMAIN}"
		return
	fi

	generate_local_certificate
	if ! certificate_is_valid; then
		printf 'Error: generated TLS certificate validation failed.\n' >&2
		return 1
	fi
}

start() {
	require_hosts_access
	if domain_has_another_mapping; then
		printf 'Error: %s already has a non-%s mapping in %s.\n' "${DOMAIN}" "${HOST_IP}" "${HOSTS_FILE}" >&2
		printf 'Remove the conflicting mapping manually and run this command again.\n' >&2
		return 1
	fi

	ensure_local_certificate
	PUID="$(id -u)" PGID="$(id -g)" \
		docker compose -f "${COMPOSE_FILE}" up -d --wait --wait-timeout 120 nginx-proxy-manager
	add_host_entry
	printf 'Nginx Proxy Manager: http://localhost:81\n'
	printf 'KaucjoMapa: https://%s (HTTP/2 and HTTP/3)\n' "${DOMAIN}"
	printf 'KaucjoMapa API: https://%s/api/\n' "${DOMAIN}"
}

stop() {
	require_hosts_access
	docker compose -f "${COMPOSE_FILE}" down
	remove_host_entry
}

main() {
	require_command docker
	if ! docker compose version >/dev/null 2>&1; then
		printf 'Error: Docker Compose plugin is not available.\n' >&2
		exit 1
	fi

	case "${1:-}" in
		start)
			start
			;;
		stop)
			stop
			;;
		*)
			usage
			exit 2
			;;
	esac
}

main "$@"
