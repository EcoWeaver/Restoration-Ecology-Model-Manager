<script lang="ts">
	import Input from "./Input.svelte";

	let { type, value = $bindable() }: { type: string; value: string } = $props();

	function set(v: string) {
		if (/^\d{0,4}$/.test(v)) {
			value = v;
		}
	}
</script>

<Input placeholder={type} bind:value={() => value, set} type="text" />
