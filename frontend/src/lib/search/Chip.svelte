<script lang="ts" module>
	export type ChipElement = {
		kind: string;
		value?: string;
	};
</script>

<script lang="ts">
	let {
		chip = $bindable(),
		editing = $bindable(),
		onenter,
	}: {
		chip: ChipElement;
		editing?: boolean;
		onenter?: () => void;
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

	{#if typeof chip.value === "string"}
		<div
			class="h-8 content-center items-center rounded-r bg-neutral-200 px-2 text-nowrap outline-indigo-400 empty:py-1 focus-visible:outline-2 dark:bg-neutral-600"
			contenteditable
			role="textbox"
			bind:innerText={chip.value}
			bind:focused={editing}
			{oninput}
			{onfocus}
			{onblur}
		></div>
	{/if}
</span>
