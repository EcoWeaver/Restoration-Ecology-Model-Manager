<script lang="ts">
	import Input from "./Input.svelte";

	let {
		placeholder,
		suggest,
		value = $bindable(),
	}: {
		placeholder: string;
		suggest: (value: string) => Promise<string[]>;
		value: string;
	} = $props();

	let suggestions = $state<string[]>([]);
	let selected = $state(0);

	let focused = $state(false);
	let inputElement: ReturnType<typeof Input>;

	function onkeydown(event: KeyboardEvent) {
		if (event.key === "Enter") {
			if (suggestions[selected]) {
				event.preventDefault();
				value = suggestions[selected];
				selected = 0;
				focused = false;
			}
		} else {
			if (event.key === "ArrowUp") {
				event.preventDefault();
				if (focused) {
					selected = Math.min(suggestions.length - 1, selected + 1);
				}
			} else if (event.key === "ArrowDown") {
				event.preventDefault();
				if (focused) {
					selected = Math.max(0, selected - 1);
				}
			}
			focused = true;
		}
	}

	function select(event: MouseEvent, idx: number) {
		inputElement.focus();
		event.preventDefault();

		if (suggestions[idx]) {
			value = suggestions[idx];
			selected = 0;
			focused = false;
		}
	}

	function onfocusout(event: FocusEvent & { currentTarget: EventTarget & HTMLDivElement }) {
		const { relatedTarget, currentTarget } = event;
		if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) {
			return;
		}

		focused = false;
	}

	$effect(() => {
		suggest(value).then((newSuggestions) => {
			suggestions = newSuggestions;
			selected = Math.min(selected, Math.max(0, newSuggestions.length - 1));
		});
	});

	let focusedElement = $state<HTMLButtonElement>();
	$effect(() => focusedElement?.scrollIntoView({ block: "nearest" }));
</script>

<div class="relative" {onfocusout}>
	{#if suggestions.length > 0 && focused}
		<div
			class="absolute bottom-9 flex max-h-64 min-w-64 flex-col-reverse overflow-y-auto rounded bg-neutral-200 p-4 shadow-xl dark:bg-neutral-700"
		>
			{#each suggestions as choice, idx (choice)}
				{#if idx === selected}
					<button
						tabindex="-1"
						onclick={(event) => select(event, idx)}
						class="my-1 block w-full cursor-pointer rounded bg-neutral-300 p-2 text-left ring-2 ring-indigo-400 dark:bg-neutral-600"
						bind:this={focusedElement}
					>
						{choice}
					</button>
				{:else}
					<button
						tabindex="-1"
						onclick={(event) => select(event, idx)}
						class="my-1 block w-full cursor-pointer rounded p-2 text-left hover:bg-neutral-300 dark:hover:bg-neutral-600"
					>
						{choice}
					</button>
				{/if}
			{/each}
		</div>
	{/if}

	<Input
		bind:this={inputElement}
		onfocus={() => (focused = true)}
		{onkeydown}
		{placeholder}
		bind:value
	/>
</div>
