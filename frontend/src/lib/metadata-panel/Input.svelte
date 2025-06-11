<script lang="ts">
	let {
		type = "text",
		placeholder,
		disabled = false,
		value = $bindable(),
	}: { type?: string; placeholder: string; disabled?: boolean; value: string } = $props();
</script>

<input
	{disabled}
	{placeholder}
	class="h-8 w-full rounded border border-neutral-600 pr-2 pl-3"
	{type}
/>
