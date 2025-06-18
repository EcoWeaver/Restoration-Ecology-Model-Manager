<script lang="ts">
	import { type NodeProps, Position, useSvelteFlow } from "@xyflow/svelte";
	import type { Node } from "$lib/model";
	import { Handle } from "@xyflow/svelte";
	import { Move, Palette, Trash2 } from "@lucide/svelte";

	const { id, data }: NodeProps<Omit<Node, "dragHandle"> & { dragHandle?: string }> = $props();
	const { updateNodeData, deleteElements } = useSvelteFlow();

	let body = $state(data.body);
	let color = $state(data.color);
	$effect(() => updateNodeData(id, { body, color }));

	function isLight(s: string) {
		const rgb = parseInt(s.slice(1), 16);
		const r = (rgb >> 16) & 0xff;
		const g = (rgb >> 8) & 0xff;
		const b = rgb & 0xff;

		const luma = r * 0.2126 + g * 0.7152 + b * 0.0722;
		return luma > 127;
	}
	const light = $derived(isLight(color));

	function changeColor() {
		const e = document.createElement("input");
		e.type = "color";
		e.value = color;
		e.click();
		e.onchange = () => (color = e.value);
	}

	function deleteSelf() {
		deleteElements({ nodes: [{ id }] });
	}
</script>

<div class="relative rounded-lg px-4 py-2" style:background-color={color}>
	<Handle
		type="source"
		position={Position.Top}
		class="!absolute !inset-0 !size-full !transform-none !rounded-none !border-0 opacity-0"
	/>
	{#if data.editable}
		<div
			class="relative max-w-prose cursor-text wrap-break-word"
			contenteditable
			class:text-black={light}
			class:text-white={!light}
			bind:innerText={body}
		></div>
		<div
			class="absolute top-full left-0 mt-1 h-10 w-max cursor-default rounded-full bg-neutral-200 p-2 dark:bg-neutral-700"
		>
			<div class="drag-handle inline-block cursor-grab"><Move /></div>
			<button class="cursor-pointer" onclick={changeColor} title="Change color"><Palette /></button>
			<button class="cursor-pointer" onclick={deleteSelf} title="Delete"><Trash2 /></button>
		</div>
	{:else}
		<div class="max-w-prose wrap-break-word">{body}</div>
	{/if}
</div>
