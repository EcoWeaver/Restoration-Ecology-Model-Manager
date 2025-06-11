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
	class="flex h-10 overflow-y-scroll border border-neutral-600 focus-within:outline-2 focus-within:outline-blue-500"
>
	{@render children?.()}

	<input
		class="block grow px-2 py-1 focus-within:outline-none"
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
