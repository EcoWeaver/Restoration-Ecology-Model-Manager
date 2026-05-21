<script lang="ts">
	import Input from "./Input.svelte";

	let { value = $bindable() }: { value: string } = $props();

	function set(v: string) {
		v = v.toLowerCase();
		if (v.length > 36) {
			return;
		}

		let idx = 0;
		for (const t of v) {
			if (idx === 8 || idx === 13 || idx === 18 || idx === 23) {
				if (t !== "-") {
					return;
				}
			} else if ((t < "0" || t > "9") && (t < "a" || t > "f")) {
				return;
			}

			idx += 1;
		}

		value = v;
	}
</script>

<Input placeholder="uuid" bind:value={() => value, set} />
