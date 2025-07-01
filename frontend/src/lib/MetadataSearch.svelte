<script lang="ts">
	import { untrack } from "svelte";
	import { goto } from "$app/navigation";
	import Chip, { type ChipElement } from "./search/Chip.svelte";
	import { suggest } from "./search/suggestions";
	import { Search as SearchIcon } from "@lucide/svelte";

	let { params = [] }: { params?: ChipElement[] } = $props();

	let input = $state("");
	const chips = $state(params);
	let chip = $state<undefined | ChipElement>(undefined);

	let selected = $state<number>(-1);
	let suggestions = $state<string[]>([]);

	let focused = $state(false);
	function onfocusout(event: FocusEvent & { currentTarget: EventTarget & HTMLDivElement }) {
		const { relatedTarget, currentTarget } = event;
		if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) {
			return;
		}

		focused = false;
	}

	let inputElement: HTMLInputElement | undefined = $state();
	function select(event: MouseEvent, selected: number) {
		inputElement?.focus();
		event.preventDefault();

		chip = { kind: suggestions[selected] };
		input = "";
		selected = -1;
	}

	async function submit() {
		const query = new URLSearchParams();
		for (const { kind, value } of chips) {
			query.append(kind, value!);
		}

		await goto(`/search?${query}`, { invalidateAll: true });
		focused = false;
	}

	async function button() {
		if (chip) {
			if (input.length > 0) {
				chip.value = input;
				chips.push(chip);

				chip = undefined;
				input = "";

				await submit();
			}
		} else if (input.length === 0) {
			await submit();
		}
	}

	async function key(key: KeyboardEvent) {
		if (key.code === "Enter") {
			if (chip) {
				if (input.length > 0) {
					chip.value = input;
					chips.push(chip);

					chip = undefined;
					input = "";
				}
			} else {
				if (selected >= 0 && suggestions[selected]) {
					chip = { kind: suggestions[selected] };
					input = "";
					selected = -1;
				} else if (selected === -1 && input.length === 0) {
					await submit();
				}
			}

			key.preventDefault();
		} else if (key.code === "Backspace" && input === "") {
			if (chip) {
				input = chip.kind;
				chip = undefined;
			} else if (chips.length > 0) {
				chip = chips.pop()!;

				input = chip.value ?? "";
				chip.value = undefined;
			}

			key.preventDefault();
		} else if (key.code === "ArrowDown") {
			selected = Math.min(selected + 1, suggestions.length - 1);
			key.preventDefault();
		} else if (key.code === "ArrowUp") {
			selected = Math.max(input === "" ? -1 : 0, selected - 1);
			key.preventDefault();
		}
	}

	let element = $state<Element>();
	$effect(() => element?.scrollIntoView({ block: "nearest" }));

	let containerElement: HTMLDivElement;
	let inputElementOffset = $state<number | undefined>();
	function computeInputElementOffset() {
		const inputElementRect = inputElement?.getBoundingClientRect();
		if (!inputElementRect) {
			inputElementOffset = undefined;
			return;
		}

		const containerRect = containerElement.getBoundingClientRect();
		const max = containerRect.x + containerRect.width - inputElementRect.width;

		inputElementOffset = Math.min(inputElementRect.x, max);
	}

	$effect(() => {
		inputElement?.scrollIntoView();
		computeInputElementOffset();

		if (chip) {
			suggestions = [];
			return;
		}

		suggest(
			input,
			untrack(() => selected),
		).then(([newSuggestions, newSelected]) => {
			suggestions = newSuggestions;
			selected = newSelected;
		});
	});
</script>

<div
	bind:this={containerElement}
	class="scrollbar-neutral-400 h-10 min-w-0 grow overflow-x-auto overflow-y-clip rounded-l-md border border-neutral-600 focus-within:outline-2 focus-within:outline-indigo-400"
	onscroll={() => computeInputElementOffset()}
	onresize={() => computeInputElementOffset()}
	{onfocusout}
>
	<span class="relative z-10 flex h-full items-center gap-1.5 px-1">
		{#each chips as chip, idx (idx)}
			<Chip {chip} />
		{/each}
		{#if chip}
			<Chip {chip} />
		{/if}

		<input
			class="block field-sizing-content min-w-24 grow pl-1 focus-visible:outline-none"
			placeholder={chips.length === 0 && !chip ? "search by metadata" : undefined}
			bind:this={inputElement}
			bind:value={input}
			onkeydown={key}
			onfocus={() => (focused = true)}
			type="text"
		/>
	</span>

	{#if focused && suggestions.length !== 0}
		<div
			style={inputElementOffset ? `left: ${inputElementOffset}px` : undefined}
			class="absolute mt-2 -ml-4 max-h-96 min-w-64 overflow-y-scroll rounded-md border border-neutral-600 bg-white p-2 py-1 shadow-xl dark:bg-neutral-800"
			tabindex="-1"
		>
			{#each suggestions as suggestion, idx (suggestion)}
				{#if idx === selected}
					<button
						tabindex="-1"
						onclick={(event) => select(event, idx)}
						class="my-1 block w-full cursor-pointer rounded-md bg-neutral-200 p-2 text-left ring-2 ring-indigo-400 dark:bg-neutral-700"
						bind:this={element}
					>
						{suggestion}
					</button>
				{:else}
					<button
						tabindex="-1"
						onclick={(event) => select(event, idx)}
						class="my-1 block w-full cursor-pointer rounded-md p-2 text-left hover:bg-neutral-200 dark:hover:bg-neutral-700"
					>
						{suggestion}
					</button>
				{/if}
			{/each}
		</div>
	{/if}
</div>

<button
	class="flex h-10 cursor-pointer place-items-center rounded-r-md border border-l-0 border-neutral-600 px-3 outline-indigo-400 focus-visible:outline-2"
	title="search"
	onclick={button}
>
	<SearchIcon />
</button>
