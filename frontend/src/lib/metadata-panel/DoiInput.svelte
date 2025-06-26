<script lang="ts">
	import { stripURL } from "./util";
	import Input from "./Input.svelte";

	let { value = $bindable() }: { value: string } = $props();

	function set(v: string) {
		value = stripURL(v, "https://doi.org/");
	}
</script>

<Input placeholder="doi" bind:value={() => value, set} type="text" />
