<script lang="ts">
	import Input from "./Input.svelte";

	let { value = $bindable() }: { value: string } = $props();

	function set(v: string) {
		if (v.length > 10) {
			return;
		}

		let digits = 4;
		for (let t of v) {
			if (digits === 0) {
				if (t !== "-") return;
			} else if (t < "0" || t > "9") {
				return;
			}

			if (digits === 0) {
				digits = 2;
			} else {
				digits -= 1;
			}
		}

		value = v;
	}
</script>

<Input placeholder="date" bind:value={() => value, set} type="text" />
