<script lang="ts">
	import { type Snippet } from "svelte";
	import { GanttChartSquare, PlusIcon } from "@lucide/svelte";
	import { slide } from "svelte/transition";

	let { children, onmetadata }: { children: Snippet; onmetadata?: () => any } = $props();
	let show = $state(true);
</script>

{#if show}
	<div
		transition:slide={{ duration: 100 }}
		class="m-3 mb-2 max-h-[55vh] overflow-y-auto rounded-xl bg-neutral-200/60 py-6 pr-10 pl-6 backdrop-blur-sm dark:bg-neutral-700/60"
	>
		{@render children()}
	</div>
{/if}

<div class="flex h-10 w-full flex-row justify-end gap-2 px-3">
	{#if onmetadata && show}
		<button
			transition:slide={{ duration: 100, axis: "x" }}
			onclick={onmetadata}
			class="cursor right-0 flex cursor-pointer items-center gap-1 rounded-full bg-neutral-200 px-4 py-1 text-nowrap focus-within:outline-2 focus-within:outline-indigo-400 hover:bg-neutral-300 dark:bg-neutral-700 dark:hover:bg-neutral-600"
		>
			<PlusIcon class="inline size-5.5" /> add metadata
		</button>
	{/if}

	<button
		title="toggle metadata"
		onclick={() => (show = !show)}
		class:bg-neutral-300={show}
		class:bg-neutral-200={!show}
		class:dark:bg-neutral-600={show}
		class:dark:bg-neutral-700={!show}
		class="cursor right-0 block cursor-pointer rounded-full p-2 focus-within:outline-2 focus-within:outline-indigo-400 hover:bg-neutral-300 dark:hover:bg-neutral-600"
	>
		<GanttChartSquare />
	</button>
</div>
