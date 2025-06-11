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

	let dontUnfocus = $state(false);
	let focused = $state(false);

	let chips = $state<ChipElement[]>(params);
	let chip = $state<ChipElement | undefined>(undefined);

	let inputElement: HTMLInputElement | undefined = $state();
	function select(event: MouseEvent, selected: number) {
		dontUnfocus = false;

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
			key.preventDefault();

			if (chip) {
				input = chip.kind;
				chip = undefined;
			} else if (chips.length > 0) {
				chip = chips.pop()!;

				input = chip.value ?? "";
				chip.value = undefined;
			}
		} else if (key.code === "ArrowDown") {
			if (!chip && selected < suggestions.length - 1) selected += 1;
			key.preventDefault();
		} else if (key.code === "ArrowUp") {
			if (!chip && selected >= 0) selected -= 1;
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

<div class="grow py-4">
	<Input
		placeholder={chips.length === 0 && !chip ? "search metadata" : undefined}
		bind:value={input}
		bind:self={inputElement}
		onkeydown={key}
		onfocus={() => (focused = true)}
		onblur={() => {
			if (!dontUnfocus) {
				focused = false;
			}
		}}
	>
		{#each chips as chip}
			<Chip {chip} />
		{/each}
		{#if chip}
			<Chip {chip} />
		{/if}
	</Input>

	{#if focused && suggestions.length != 0}
		<div class="absolute mt-px max-h-96 min-w-64 overflow-y-scroll bg-white p-2 shadow-xl">
			{#each suggestions as suggestion, idx}
				{#if idx === selected}
					<button
						onmousedown={() => (dontUnfocus = true)}
						onclick={(event) => select(event, idx)}
						onmouseleave={() => {
							if (dontUnfocus) {
								focused = false;
								dontUnfocus = false;
							}
						}}
						class="my-1 block w-full cursor-pointer rounded bg-neutral-200 p-2 text-left ring-2 ring-blue-500 hover:bg-neutral-200"
						bind:this={element}
					>
						{suggestion}
					</button>
				{:else}
					<button
						onmousedown={() => (dontUnfocus = true)}
						onclick={(event) => select(event, idx)}
						onmouseleave={() => {
							if (dontUnfocus) {
								focused = false;
								dontUnfocus = false;
							}
						}}
						class="my-1 block w-full cursor-pointer rounded p-2 text-left hover:bg-neutral-200"
					>
						{suggestion}
					</button>
				{/if}
			{/each}
		</div>
	{/if}
</div>

<SearchButton onclick={submit} />
