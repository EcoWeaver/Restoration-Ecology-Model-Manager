<script lang="ts">
	import { type NodeProps, Position, useSvelteFlow } from "@xyflow/svelte";
	import type { Node } from "$lib/model";
	import { Handle } from "@xyflow/svelte";

	let { id, data }: NodeProps<Node> = $props();
	const { updateNodeData } = useSvelteFlow();

	let body = $state(data.body);
	$effect(() => updateNodeData(id, { body }));

	function isLight(s: string) {
		const rgb = parseInt(s.slice(1), 16);
		const r = (rgb >> 16) & 0xff;
		const g = (rgb >> 8) & 0xff;
		const b = rgb & 0xff;

		const luma = r * 0.2126 + g * 0.7152 + b * 0.0722;
		return luma > 127;
	}
</script>

<div
	class="rounded-lg px-4 py-2"
	class:text-white={!isLight(data.color)}
	style:background-color={data.color}
>
	{#if data.editable}
		<span class="cursor-text" contenteditable bind:innerText={body}></span>
	{:else}
		<span>{body}</span>
	{/if}
</div>

<Handle type="target" position={Position.Left} />
<Handle type="source" position={Position.Right} />
