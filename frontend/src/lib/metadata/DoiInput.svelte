<script lang="ts">
	import { stripURL } from "./util";
	import Input from "./Input.svelte";

	let { value = $bindable() }: { value: string } = $props();

	function set(v: string) {
		v = stripURL(v, "https://doi.org/").toLowerCase();

		let idx = 0;
		let prefix = true;
		for (const t of v) {
			if (prefix) {
				if (idx === 0) {
					if (t !== "1") {
						return;
					}
				}
				if (idx === 1) {
					if (t !== "0") {
						return;
					}
				}
				if (idx === 2) {
					if (t !== ".") {
						return;
					}
				} else if (t < "0" || t > "9") {
					if (t === "/") {
						if (idx < 7) {
							return;
						} else {
							prefix = false;
						}
					} else {
						return;
					}
				} else if (idx === 12) {
					return;
				}

				idx += 1;
			} else {
				if (
					(t < "0" || t > "9") &&
					(t < "a" || t > "z") &&
					["-", ".", "_", ";", "(", ")", "\\", ":"].indexOf(t) === -1
				) {
					return;
				}
			}
		}

		value = v;
	}
</script>

<Input placeholder="doi" bind:value={() => value, set} />
