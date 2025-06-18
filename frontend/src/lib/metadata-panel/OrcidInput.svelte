<script lang="ts">
	import Input from "./Input.svelte";

	let { value = $bindable() }: { value: string } = $props();

	function set(v: string) {
		if (v.length > 19) return;
		v = v.toUpperCase();

		let digits = 0;
		let clusters = 0;
		for (let t of v) {
			if (digits === 4) {
				if (t !== "-") return;
			} else if (t < "0" || t > "9") {
				if (clusters === 3 && digits === 3 && t === "X") break;
				return;
			}

			if (digits === 4) {
				digits = 0;
				clusters += 1;
			} else {
				digits += 1;
			}
		}

		value = v;
	}
</script>

<Input placeholder="orcid" bind:value={() => value, set} type="text" />
