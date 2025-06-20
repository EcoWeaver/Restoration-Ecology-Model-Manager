<script lang="ts">
	import { type Snippet } from "svelte";

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
		onsubmit?: () => any;
		onkeydown?: (key: KeyboardEvent) => any;
		onfocus?: () => any;
		onblur?: () => any;
	} = $props();
</script>

<span
	class="relative z-10 flex h-10 overflow-x-auto border border-neutral-600 focus-within:outline-2 focus-within:outline-indigo-400"
>
	{@render children?.()}

	<input
		class="block field-sizing-content min-w-16 grow px-2 py-1 focus-within:outline-none"
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
