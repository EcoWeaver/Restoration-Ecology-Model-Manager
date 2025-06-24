<script lang="ts">
	import { type Snippet } from "svelte";
	import type { FocusEventHandler, KeyboardEventHandler } from "svelte/elements";

	let {
		children,
		placeholder,
		value = $bindable(),
		self = $bindable(),
		onkeydown,
		onfocus,
		onblur,
	}: {
		children?: Snippet;
		placeholder?: string;
		value: string;
		self?: HTMLInputElement;
		onkeydown?: KeyboardEventHandler<HTMLInputElement>;
		onfocus?: FocusEventHandler<HTMLInputElement>;
		onblur?: FocusEventHandler<HTMLInputElement>;
	} = $props();
</script>

<span class="relative z-10 flex h-10 overflow-x-auto border border-neutral-600">
	{@render children?.()}

	<input
		class="block field-sizing-content min-w-16 grow px-2 py-1 focus-visible:outline-none"
		{placeholder}
		bind:this={self}
		bind:value
		defaultValue={value}
		{onkeydown}
		{onfocus}
		{onblur}
		type="text"
	/>
</span>
