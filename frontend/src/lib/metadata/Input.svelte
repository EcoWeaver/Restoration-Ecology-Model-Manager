<script lang="ts">
	import type { FocusEventHandler, KeyboardEventHandler } from "svelte/elements";

	let {
		type = "text",
		placeholder,
		disabled = false,
		onkeydown,
		onfocus,
		value = $bindable(),
	}: {
		type?: string;
		placeholder: string;
		disabled?: boolean;
		onkeydown?: KeyboardEventHandler<HTMLInputElement>;
		onfocus?: FocusEventHandler<HTMLInputElement>;
		value: string;
	} = $props();

	let inputElement: HTMLInputElement;
	export function focus() {
		inputElement.focus();
	}
</script>

<input
	bind:this={inputElement}
	{disabled}
	{placeholder}
	bind:value
	{onkeydown}
	{onfocus}
	class="field-sizing-content h-8 w-full rounded border border-neutral-600 pr-2 pl-3 outline-indigo-400 focus-visible:outline-2"
	{type}
/>
