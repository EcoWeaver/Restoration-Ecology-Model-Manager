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

	function onkeydown(event: KeyboardEvent) {
		if (event.key === "Enter") {
			if (suggestions[selected]) {
				value = suggestions[selected];
			}
		} else if (event.key === "ArrowUp") {
			event.preventDefault();
			selected = Math.min(10, suggestions.length - 1, selected + 1);
		} else if (event.key === "ArrowDown") {
			event.preventDefault();
			selected = Math.max(0, selected - 1);
		}
	}

	let dontUnfocus = $state(false);
	let focused = $state(false);
	let inputElement: ReturnType<typeof Input>;
	function select(event: MouseEvent, idx: number) {
		dontUnfocus = false;

		inputElement.focus();
		event.preventDefault();

		if (suggestions[idx]) {
			value = suggestions[idx];
			selected = 0;
		}
	}

	$effect(() => {
		suggest(value).then((newSuggestions) => {
			suggestions = newSuggestions;
			if (suggestions.length === 0) {
				selected = 0;
			}
		});
	});
</script>

<div class="relative">
	{#if suggestions.length > 0 && focused}
		<div
			class="absolute bottom-9 flex min-w-64 flex-col-reverse rounded bg-neutral-200 p-4 dark:bg-neutral-700"
		>
			{#each suggestions.slice(0, 10) as choice, idx}
				<button
					tabindex="-1"
					class:ring-2={idx === selected}
					class:bg-neutral-300={idx === selected}
					class:dark:bg-neutral-600={idx === selected}
					onmousedown={() => (dontUnfocus = true)}
					onmouseleave={() => {
						if (dontUnfocus) {
							focused = false;
							dontUnfocus = false;
						}
					}}
					onclick={(event) => select(event, idx)}
					class="my-1 block w-full cursor-pointer rounded p-2 text-left ring-indigo-400 hover:bg-neutral-300 dark:hover:bg-neutral-600"
				>
					{choice}
				</button>
			{/each}
		</div>
	{/if}

	<Input
		bind:this={inputElement}
		onfocus={() => (focused = true)}
		onblur={() => {
			if (!dontUnfocus) {
				focused = false;
			}
		}}
		{onkeydown}
		{placeholder}
		bind:value
	/>
</div>
