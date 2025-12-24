<script lang="ts">
	import { untrack } from "svelte";
	import { goto } from "$app/navigation";
	import Chip, { type ChipElement } from "./search/Chip.svelte";
	import { metadataSuggest } from "./search/suggestions";
	import { Search as SearchIcon } from "@lucide/svelte";
	import { resolve } from "$app/paths";
	import { useElements } from "./api";

	let { params = [] }: { params?: ChipElement[] } = $props();

	const { getMetadataElements } = useElements();

	let value = $state("");
	let chip = $state<ChipElement>();
	// $derived is not deeply reactive, so you have to put a $state
	// inside of the $derived.
	// see https://github.com/sveltejs/svelte/issues/16189
	let chips = $derived.by(() => {
		const chips = $state(params);
		return chips;
	});

	let suggestions = $state<string[]>([]);
	let selected = $state<number>(-1);

	let focused = $state(false);
	function onfocusout(event: FocusEvent & { currentTarget: EventTarget & HTMLDivElement }) {
		const { relatedTarget, currentTarget } = event;
		if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) {
			return;
		}

		focused = false;
	}

	let inputElement = $state<HTMLInputElement>();
	function select(event: MouseEvent, selected: number) {
		inputElement?.focus();
		event.preventDefault();

		chip = { kind: suggestions[selected] };
		value = "";
		selected = -1;
	}

	async function onkeydown(event: KeyboardEvent) {
		if (event.code === "Enter") {
			if (chip) {
				if (value.length > 0) {
					chip.value = value;
					chips.push(chip);

					chip = undefined;
					value = "";
				}
			} else {
				if (selected >= 0 && suggestions[selected]) {
					chip = { kind: suggestions[selected] };
					value = "";
					selected = -1;
				} else if (selected === -1 && value.length === 0) {
					await submit();
				}
			}

			event.preventDefault();
		} else if (event.code === "Backspace" && value === "") {
			if (chip) {
				value = chip.kind;
				chip = undefined;
			} else if (chips.length > 0) {
				chip = chips.pop()!;

				value = chip.value ?? "";
				chip.value = undefined;
			}

			event.preventDefault();
		} else if (event.code === "ArrowDown") {
			selected = Math.min(selected + 1, suggestions.length - 1);
			event.preventDefault();
		} else if (event.code === "ArrowUp") {
			selected = Math.max(value === "" ? -1 : 0, selected - 1);
			event.preventDefault();
		}
	}

	async function onclick() {
		if (chip) {
			if (value.length > 0) {
				chip.value = value;
				chips.push(chip);

				chip = undefined;
				value = "";

				await submit();
			}
		} else if (value.length === 0) {
			await submit();
		}
	}

	async function submit() {
		const query = new URLSearchParams();
		for (const { kind, value } of chips) {
			query.append(kind, value!);
		}

		// eslint-disable-next-line svelte/no-navigation-without-resolve
		await goto(resolve("/search") + `?${query}`);
		focused = false;
	}

	let containerElement: HTMLDivElement;
	let inputElementOffset = $state<number>();
	function computeInputElementOffset() {
		if (!inputElement) {
			inputElementOffset = undefined;
			return;
		}

		const inputElementRect = inputElement.getBoundingClientRect();
		const containerRect = containerElement.getBoundingClientRect();
		const max = containerRect.x + containerRect.width - inputElementRect.width;

		inputElementOffset = Math.min(inputElementRect.x, max);
	}

	$effect(() => {
		// recompute offset on chip changes, as deleting
		// and editing chips can change their size
		$state.snapshot(chip);
		$state.snapshot(chips);
		computeInputElementOffset();
	});

	const suggest = metadataSuggest(getMetadataElements());
	$effect(() => {
		if (chip) {
			suggestions = [];
		} else {
			[suggestions, selected] = suggest(
				value,
				untrack(() => selected),
			);
		}
	});

	const chipsEditing = $state<boolean[]>([]);

	let focusedElement = $state<Element>();
	$effect(() => focusedElement?.scrollIntoView({ block: "nearest" }));
</script>

<div
	bind:this={containerElement}
	class="h-10 min-w-0 grow overflow-x-auto overflow-y-clip rounded-l-md border border-neutral-600 focus-within:outline-2 focus-within:outline-indigo-400"
	onscroll={() => computeInputElementOffset()}
	onresize={() => computeInputElementOffset()}
	{onfocusout}
>
	<div class:pl-3={chips.length === 0 && !chip} class="relative z-10 flex h-full gap-0.75 px-0.75">
		{#each chips, idx (idx)}
			<Chip
				onremove={() => {
					chips.splice(idx, 1);
					inputElement?.focus();
				}}
				onenter={() => inputElement?.focus()}
				bind:editing={chipsEditing[idx]}
				bind:chip={chips[idx]}
			/>
		{/each}
		{#if chip}
			<Chip {chip} />
		{/if}

		<input
			class="field-sizing-content min-w-24 grow focus-visible:outline-none"
			placeholder={chips.length === 0 && !chip ? "search uploaded models" : undefined}
			bind:this={inputElement}
			bind:value
			{onkeydown}
			onfocus={() => (focused = true)}
			type="text"
		/>
	</div>

	{#if focused && suggestions.length > 0 && !chipsEditing.includes(true)}
		<div
			style={inputElementOffset ? `left: ${inputElementOffset}px` : undefined}
			class="absolute mt-2 -ml-4 max-h-96 min-w-64 overflow-y-scroll rounded-md border border-neutral-600 bg-white p-1.75 py-1 shadow-xl dark:bg-neutral-800"
			tabindex="-1"
		>
			{#each suggestions as suggestion, idx (suggestion)}
				{#if idx === selected}
					<button
						tabindex="-1"
						onclick={(event) => select(event, idx)}
						class="my-1 block w-full cursor-pointer rounded bg-neutral-200 p-2 text-left ring-2 ring-indigo-400 dark:bg-neutral-700"
						bind:this={focusedElement}
					>
						{suggestion}
					</button>
				{:else}
					<button
						tabindex="-1"
						onclick={(event) => select(event, idx)}
						class="my-1 block w-full cursor-pointer rounded p-2 text-left hover:bg-neutral-200 dark:hover:bg-neutral-700"
					>
						{suggestion}
					</button>
				{/if}
			{/each}
		</div>
	{/if}
</div>

<button
	class="flex h-10 cursor-pointer items-center rounded-r-md border border-l-0 border-neutral-600 px-3 outline-indigo-400 focus-visible:outline-2"
	title="search"
	{onclick}
>
	<SearchIcon />
</button>
