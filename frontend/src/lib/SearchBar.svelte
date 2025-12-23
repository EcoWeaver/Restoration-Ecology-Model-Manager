<script lang="ts">
	import { untrack } from "svelte";
	import { goto } from "$app/navigation";
	import Chip, { type ChipElement } from "./search/Chip.svelte";
	import { metadataSuggest } from "./search/suggestions";
	import { Search as SearchIcon } from "@lucide/svelte";
	import { resolve } from "$app/paths";
	import { useMetadataElements } from "./api";

	let { params = [] }: { params?: ChipElement[] } = $props();

	const { getMetadataElements } = useMetadataElements();
	const suggest = metadataSuggest(getMetadataElements());

	let input = $state("");
	let chip = $state<ChipElement>();
	// $derived is not deeply reactive, so you have to put a $state
	// inside of the $derived.
	// see https://github.com/sveltejs/svelte/issues/16189
	let chips = $derived.by(() => {
		const chips = $state(params);
		return chips;
	});

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

	let inputElement = $state<HTMLInputElement>();
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

		// eslint-disable-next-line svelte/no-navigation-without-resolve
		await goto(resolve("/search") + `?${query}`);
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
	let inputElementOffset = $state<number>();
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
		// recompute offset on chip changes, as deleting
		// and editing chips can change their size
		$state.snapshot(chips);
		computeInputElementOffset();
	});

	const chipsEditing = $state<boolean[]>([]);
	$effect(() => {
		inputElement?.scrollIntoView();
		computeInputElementOffset();

		if (chip) {
			suggestions = [];
			return;
		}

		[suggestions, selected] = suggest(
			input,
			untrack(() => selected),
		);
	});
</script>

<div
	bind:this={containerElement}
	class="h-10 min-w-0 grow overflow-x-auto overflow-y-clip rounded-l-md border border-neutral-600 focus-within:outline-2 focus-within:outline-indigo-400"
	onscroll={() => computeInputElementOffset()}
	onresize={() => computeInputElementOffset()}
	{onfocusout}
>
	<div
		class:pl-3={chips.length === 0 && !chip}
		class="relative z-10 flex h-full items-center gap-0.75 px-0.75"
	>
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
			class="block field-sizing-content min-w-24 grow focus-visible:outline-none"
			placeholder={chips.length === 0 && !chip ? "search uploaded models" : undefined}
			bind:this={inputElement}
			bind:value={input}
			onkeydown={key}
			onfocus={() => (focused = true)}
			type="text"
		/>
	</div>

	{#if focused && suggestions.length !== 0 && !chipsEditing.includes(true)}
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
