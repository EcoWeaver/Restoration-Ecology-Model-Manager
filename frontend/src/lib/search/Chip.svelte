<script lang="ts" module>
	export type ChipElement = {
		kind: string;
		value?: string;
	};
</script>

<script lang="ts">
	import { X } from "@lucide/svelte";

	let {
		chip = $bindable(),
		editing = $bindable(),
		onenter,
		onremove,
	}: {
		chip: ChipElement;
		editing?: boolean;
		onenter?: () => void;
		onremove?: () => void;
	} = $props();

	let prev: string | undefined;

	function onfocus() {
		prev = chip.value;
	}

	function onblur() {
		if (prev && chip.value?.length === 0) {
			// if the contenteditable gets unfocussed, but the
			// current value for the chip is empty, restore the
			// previous value so that the chip cannot be empty
			chip.value = prev;
		}
	}

	function oninput(event: Event & { currentTarget: EventTarget & HTMLDivElement }) {
		if (chip.value?.includes("\n")) {
			chip.value = chip.value.replaceAll("\n", "");
			if (chip.value.length > 0) {
				event.currentTarget.blur();
				onenter?.();
			}
		}
	}
</script>

<span class="flex items-center gap-0.5">
	<div class="h-8 content-center rounded-l bg-neutral-200 px-2 text-nowrap dark:bg-neutral-600">
		{chip.kind}
	</div>
	<div class="h-8 content-center bg-neutral-200 px-2 dark:bg-neutral-600">≈</div>

	{#if chip.value}
		<div
			class="flex h-8 items-center rounded-r bg-neutral-200 pr-2 text-nowrap outline-indigo-400 focus-within:outline-2 dark:bg-neutral-600"
		>
			<div
				class="px-2 outline-none"
				contenteditable="plaintext-only"
				role="textbox"
				bind:innerText={chip.value}
				bind:focused={editing}
				{oninput}
				{onfocus}
				{onblur}
			></div>

			<button
				onclick={onremove}
				class="cursor-pointer rounded-full outline-offset-2 outline-indigo-400 focus-visible:outline-2"
			>
				<X size={16} />
			</button>
		</div>
	{/if}
</span>
