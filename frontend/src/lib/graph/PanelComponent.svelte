<script lang="ts">
	import { type Snippet } from "svelte";

	let { children }: { children: Snippet } = $props();
	let show = $state(true);
</script>

{#if show}
	<div
		class="m-3 mb-2 max-h-[55vh] overflow-y-scroll rounded-xl bg-neutral-200/60 px-8 py-4 backdrop-blur-sm"
	>
		{@render children()}
	</div>
{/if}

<div class="flex w-full flex-row-reverse px-3">
	<button
		onclick={() => (show = !show)}
		class="cursor right-0 block cursor-pointer rounded bg-neutral-200 p-3 py-1 hover:bg-neutral-300"
	>
		{show ? "hide" : "show metadata"}
	</button>
</div>
