splits() {
	PWD="$(pwd)"
	for cmd in "$@"; do
		kitty @ launch --cwd="$PWD" sh -c "$cmd || sleep infinity" || return 1
	done
}

splits \
	"cd backend && make up" \
	"cd backend && make run" \
	"cd frontend && npm run dev"
