kitty_splits() {
	CWD="$(pwd)"
	for cmd in "$@"; do
		kitty @ launch --cwd="$CWD" sh -c "$cmd || sleep infinity" || return 1
	done
}

ptyxis_splits() {
	for cmd in "$@"; do
		ptyxis --tab --title "$cmd" -- sh -c "$cmd || sleep infinity"
	done
}

tmux_splits() {
	tmux new-session -d "$1 || sleep infinity"
	for cmd in "${@:2}"; do
		tmux split-window "$cmd || sleep infinity"
	done
	tmux select-layout even-horizontal
	tmux attach
}

splits() {
	if [ "$TERM" == "xterm-kitty" ]; then
		kitty_splits "$@" && return
	fi

	if [ -v PTYXIS_VERSION ]; then
		ptyxis_splits "$@" && return
	fi

	tmux_splits "$@"
}

splits \
	"cd backend && make up" \
	"cd backend && make run" \
	"cd frontend && npm run dev"
