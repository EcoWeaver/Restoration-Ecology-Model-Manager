<script lang="ts">
	import Input from "./Input.svelte";
	import { stripURL } from "./util";

	let { value = $bindable() }: { value: string } = $props();

	function set(v: string) {
		v = stripURL(v, "https://orcid.org/").toUpperCase();
		if (v.length > 19) {
			return;
		}

		let idx = 0;
		for (const t of v) {
			if (idx % 5 === 4) {
				if (t !== "-") {
					return;
				}
			} else if (t < "0" || t > "9") {
				if (idx !== 18 || t !== "X") {
					return;
				}
			}

			idx += 1;
		}

		value = v;
	}
</script>

<Input placeholder="orcid" bind:value={() => value, set} />
