let timer: undefined | number = undefined;
const nav = $state({ delayed: false });

export function startNavigate() {
	timer = setTimeout(() => (nav.delayed = true), 500);
}

export function clearNavigate() {
	clearTimeout(timer);
	nav.delayed = false;
}

export { nav };
