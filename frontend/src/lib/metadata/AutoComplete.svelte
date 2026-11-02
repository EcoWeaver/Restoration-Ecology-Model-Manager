<script lang="ts">
	import { req, type Metadata } from "#lib/api.js";
	import Input from "./Input.svelte";

	let {
		placeholder,
		meta = $bindable(),
	}: {
		placeholder: string;
		meta: Metadata;
	} = $props();

	let suggestions = $state<string[]>([]);
	let selected = $state(0);

	let focused = $state(false);
	function onfocusout(event: FocusEvent & { currentTarget: EventTarget & HTMLDivElement }) {
		const { relatedTarget, currentTarget } = event;
		if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) {
			return;
		}

		focused = false;
	}

	let inputElement: ReturnType<typeof Input>;
	function select(event: MouseEvent, idx: number) {
		inputElement.focus();
		event.preventDefault();

		if (suggestions[idx]) {
			meta.value = suggestions[idx];
			selected = 0;
			focused = false;
		}
	}

	function onkeydown(event: KeyboardEvent) {
		if (event.key === "Enter") {
			if (suggestions[selected]) {
				meta.value = suggestions[selected];
				selected = 0;
				focused = false;
				event.preventDefault();
			}
		} else if (focused) {
			if (event.key === "ArrowUp") {
				selected = Math.min(suggestions.length - 1, selected + 1);
				event.preventDefault();
			} else if (event.key === "ArrowDown") {
				selected = Math.max(0, selected - 1);
				event.preventDefault();
			}
		} else {
			focused = true;
		}
	}

	let controller: AbortController | undefined = undefined;
	async function suggest() {
		controller?.abort();
		if (meta.value === "") {
			suggestions = [];
			selected = 0;
		} else {
			const c = new AbortController();
			controller = c;
			try {
				const newSuggestions = await req<string[]>(
					"GET",
					`metadata/suggest?kind=${meta.kind}&q=${encodeURIComponent(meta.value)}`,
					undefined,
					c.signal,
				);
				const unique = new Set(newSuggestions);
				suggestions = newSuggestions.filter((suggestion) => unique.delete(suggestion));
				selected = Math.min(selected, Math.max(0, suggestions.length - 1));
			} catch (e) {
				if (!c.signal.aborted) {
					/* eslint-disable-next-line no-console */
					console.error(e);
				}
			}
		}
	}

	$effect(() => {
		$state.snapshot(meta.value);
		suggest();
	});

	let focusedElement = $state<HTMLButtonElement>();
	$effect(() => focusedElement?.scrollIntoView({ block: "nearest" }));
</script>

<div class="relative" {onfocusout}>
	{#if focused && suggestions.length > 0}
		<div
			class="absolute bottom-9.5 -ml-0.75 flex max-h-64 min-w-64 flex-col-reverse overflow-y-scroll rounded-md border border-neutral-600 bg-neutral-200 p-1.75 py-1 shadow-xl dark:bg-neutral-700"
			tabindex="-1"
		>
			{#each suggestions as suggestion, idx (suggestion)}
				{#if idx === selected}
					<button
						tabindex="-1"
						onclick={(event) => select(event, idx)}
						class="my-1 block w-full cursor-pointer rounded bg-neutral-300 p-2 text-left ring-2 ring-indigo-400 dark:bg-neutral-600"
						bind:this={focusedElement}
					>
						{suggestion}
					</button>
				{:else}
					<button
						tabindex="-1"
						onclick={(event) => select(event, idx)}
						class="my-1 block w-full cursor-pointer rounded p-2 text-left hover:bg-neutral-300 dark:hover:bg-neutral-600"
					>
						{suggestion}
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
		bind:value={meta.value}
	/>
</div>
