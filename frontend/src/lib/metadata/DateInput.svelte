<script lang="ts">
	import Input from "./Input.svelte";

	let { value = $bindable() }: { value: string } = $props();

	function set(v: string) {
		if (v.length > 10) {
			return;
		}

		let idx = 0;
		for (const t of v) {
			if (idx === 4 || idx === 7) {
				if (t !== "-") {
					return;
				}
			} else if (t < "0" || t > "9") {
				return;
			}

			idx += 1;
		}

		value = v;
	}
</script>

<Input placeholder="date" bind:value={() => value, set} />
