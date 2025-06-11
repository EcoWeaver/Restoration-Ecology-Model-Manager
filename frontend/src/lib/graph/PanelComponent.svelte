<script lang="ts">
	import { type Snippet } from "svelte";

	let { children, onmetadata }: { children: Snippet; onmetadata?: () => any } = $props();
	let show = $state(true);
</script>

{#if show}
	<div
		class="m-3 mb-2 max-h-[55vh] overflow-y-auto rounded-xl bg-neutral-200/60 px-8 py-6 backdrop-blur-sm"
	>
		{@render children()}
	</div>
{/if}

<div class="flex h-8 w-full flex-row-reverse gap-2 px-3">
	<button
		onclick={() => (show = !show)}
		class="cursor right-0 block cursor-pointer rounded-md bg-neutral-200 p-3 py-1 hover:bg-neutral-300"
	>
		{show ? "hide" : "show metadata"}
	</button>

	{#if onmetadata && show}
		<button
			onclick={onmetadata}
			class="cursor right-0 block cursor-pointer rounded-md bg-neutral-200 p-3 py-1 hover:bg-neutral-300"
		>
			+ add metadata
		</button>
	{/if}
</div>
