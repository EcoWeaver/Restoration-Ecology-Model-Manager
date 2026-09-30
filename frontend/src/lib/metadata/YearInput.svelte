<script lang="ts">
	import Input from "./Input.svelte";

	let { type, value = $bindable() }: { type: string; value: string } = $props();

	function set(v: string) {
		if (v.length > 4) {
			return;
		}

		for (const t of v) {
			if (t < "0" || t > "9") {
				return;
			}
		}

		value = v;
	}
</script>

<Input placeholder={type} bind:value={() => value, set} />
