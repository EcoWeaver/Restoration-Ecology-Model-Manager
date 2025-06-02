<script lang="ts">
	let { type, value = $bindable() }: { type: string; value: string } = $props();

	function set(v: string) {
		if (/^\d{0,4}$/.test(v)) {
			value = v;
		}
	}
</script>

<input placeholder={type} class="h-6 w-full outline" bind:value={() => value, set} type="text" />
