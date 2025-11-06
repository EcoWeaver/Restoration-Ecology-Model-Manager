<script lang="ts">
	import { type Snippet } from "svelte";
	import type { MouseEventHandler } from "svelte/elements";
	import { PlusIcon, RedoIcon, SquareChartGantt, UndoIcon } from "@lucide/svelte";
	import { slide } from "svelte/transition";

	let {
		children,
		onaddmetadata,
		history,
	}: {
		children: Snippet;
		onaddmetadata?: MouseEventHandler<HTMLButtonElement>;
		history?: {
			undo?: MouseEventHandler<HTMLButtonElement>;
			redo?: MouseEventHandler<HTMLButtonElement>;
		};
	} = $props();

	let show = $state(true);
</script>

{#if show}
	<div
		transition:slide={{ duration: 100 }}
		tabindex="-1"
		class="m-3 mb-2 max-h-[55vh] overflow-y-auto rounded-xl bg-neutral-200/72 py-6 pr-10 pl-6 backdrop-blur dark:bg-neutral-700/72"
	>
		{@render children()}
	</div>
{/if}

<div class="flex h-10 w-full flex-row justify-end gap-2 px-3">
	{#if history}
		<button
			title={history.undo ? "undo" : undefined}
			onclick={history.undo}
			class:cursor-pointer={history.undo}
			class:hover:bg-neutral-300={history.undo}
			class:dark:hover:bg-neutral-600={history.undo}
			class:text-neutral-400={!history.undo}
			class="cursor right-0 block rounded-full bg-neutral-200 p-2 outline-indigo-400 focus-visible:outline-2 dark:bg-neutral-700"
		>
			<UndoIcon />
		</button>
		<button
			title={history.redo ? "redo" : undefined}
			onclick={history.redo}
			class:cursor-pointer={history.redo}
			class:hover:bg-neutral-300={history.redo}
			class:dark:hover:bg-neutral-600={history.redo}
			class:text-neutral-400={!history.redo}
			class="cursor right-0 block rounded-full bg-neutral-200 p-2 outline-indigo-400 focus-visible:outline-2 dark:bg-neutral-700"
		>
			<RedoIcon />
		</button>
	{/if}

	{#if onaddmetadata && show}
		<button
			title="add metadata"
			transition:slide={{ duration: 100, axis: "x" }}
			onclick={onaddmetadata}
			class="cursor right-0 flex cursor-pointer items-center gap-1 rounded-full bg-neutral-200 px-4 py-1 text-nowrap outline-indigo-400 hover:bg-neutral-300 focus-visible:outline-2 dark:bg-neutral-700 dark:hover:bg-neutral-600"
		>
			<PlusIcon class="inline size-5.5" />
			<span>add metadata</span>
		</button>
	{/if}

	<button
		title="toggle metadata"
		onclick={() => (show = !show)}
		class:bg-neutral-300={show}
		class:bg-neutral-200={!show}
		class:dark:bg-neutral-600={show}
		class:dark:bg-neutral-700={!show}
		class="cursor right-0 block cursor-pointer rounded-full p-2 outline-indigo-400 hover:bg-neutral-300 focus-visible:outline-2 dark:hover:bg-neutral-600"
	>
		<SquareChartGantt />
	</button>
</div>
