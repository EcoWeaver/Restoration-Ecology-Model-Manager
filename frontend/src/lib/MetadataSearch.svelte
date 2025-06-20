<script lang="ts">
	import { goto } from "$app/navigation";
	import { metadataElements } from "./api";
	import Chip, { type ChipElement } from "./search/Chip.svelte";
	import Input from "./search/Input.svelte";
	import SearchButton from "./search/Button.svelte";

	let { params = [] }: { params?: ChipElement[] } = $props();

	const elements = metadataElements.then((e) =>
		Object.values(e)
			.flat()
			.map((e) => e.name),
	);

	let input = $state("");
	let selected = $state<number>(-1);
	let suggestions = $state<string[]>([]);

	let chips = $state<ChipElement[]>(params);
	let chip = $state<ChipElement | undefined>(undefined);

	let focused = $state(false);
	function onfocusout(event: FocusEvent & { currentTarget: EventTarget & HTMLDivElement }) {
		let { relatedTarget, currentTarget } = event;
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
		let query = new URLSearchParams();
		for (let { kind, value } of chips) {
			query.set(kind, value!);
		}

		await goto(`/search?${query}`, { invalidateAll: true });
		focused = false;
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
			selected = Math.max(0, selected - 1);
			key.preventDefault();
		}
	}

	let element = $state<Element>();
	$effect(() => element?.scrollIntoView({ block: "nearest" }));

	$effect(() => {
		if (chip) {
			suggestions = [];
			return;
		}

		input;
		elements.then((elements) => {
			if (input === "") {
				selected = -1;
				suggestions = elements;
			} else {
				suggestions = elements.filter((elem) => elem.startsWith(input));
				if (selected === -1 && suggestions.length > 0) {
					selected = 0;
				}
			}
		});
	});
</script>

<div class="grow" {onfocusout}>
	<Input
		placeholder={chips.length === 0 && !chip ? "search metadata" : undefined}
		bind:value={input}
		bind:self={inputElement}
		onkeydown={key}
		onfocus={() => (focused = true)}
	>
		{#each chips as chip}
			<Chip {chip} />
		{/each}
		{#if chip}
			<Chip {chip} />
		{/if}
	</Input>

	{#if focused && suggestions.length !== 0}
		<div
			tabindex="-1"
			class="absolute mt-px max-h-96 min-w-64 overflow-y-scroll bg-neutral-100 p-2 shadow-xl dark:bg-neutral-700"
		>
			{#each suggestions as suggestion, idx}
				{#if idx === selected}
					<button
						tabindex="-1"
						onclick={(event) => select(event, idx)}
						class="my-1 block w-full cursor-pointer rounded bg-neutral-200 p-2 text-left ring-2 ring-indigo-400 dark:bg-neutral-600"
						bind:this={element}
					>
						{suggestion}
					</button>
				{:else}
					<button
						tabindex="-1"
						onclick={(event) => select(event, idx)}
						class="my-1 block w-full cursor-pointer rounded p-2 text-left hover:bg-neutral-200 dark:hover:bg-neutral-600"
					>
						{suggestion}
					</button>
				{/if}
			{/each}
		</div>
	{/if}
</div>

<SearchButton
	onclick={() => {
		if (selected === -1 && input.length === 0) {
			submit();
		}
	}}
/>
